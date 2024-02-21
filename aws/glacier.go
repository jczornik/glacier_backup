package aws

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
)

const (
	multipartThreshold = 1024 * 1024 * 100
	minPartSize        = 1024 * 1024
	maxPartSize        = 1024 * 1024 * 1024 * 4
	maxParts           = 10000
)

type part struct {
	reader       io.Reader
	sha256       string
	contentRange string
}

func GetVaults(cfg aws.Config, account string) ([]string, error) {
	client := glacier.NewFromConfig(cfg)
	vaults, err := client.ListVaults(context.TODO(), &glacier.ListVaultsInput{AccountId: &account})
	if err != nil {
		return nil, err
	}

	vaultNames := make([]string, len(vaults.VaultList))
	for i, element := range vaults.VaultList {
		vaultNames[i] = *element.VaultName
	}

	return vaultNames, nil
}

func upload(client *glacier.Client, account string, vault string, file *os.File) error {
	fileName := file.Name()

	input := glacier.UploadArchiveInput{
		AccountId:          &account,
		VaultName:          &vault,
		Body:               file,
		ArchiveDescription: &fileName,
	}

	r, err := client.UploadArchive(context.TODO(), &input)
	if err != nil {
		return err
	}

	file.Seek(0, 0)
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	hashes := computeHashes(content)
	treeHash := computeTreeHash(hashes)

	if *r.Checksum != hex.EncodeToString(treeHash) {
		return errors.New("Checksums mismatch")
	}

	return err
}

func powInt(x, y int) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}

func computePartSize(archiveSize int64) (int64, error) {
	var partSize int64
	if archiveSize < minPartSize {
		return partSize, errors.New("Archive size is too small")
	}

	if archiveSize > maxPartSize*maxParts {
		return partSize, errors.New("Archive size is too big")
	}

	partSize = archiveSize / maxParts

	mult := 0
	for minPartSize*powInt(2, mult) < partSize {
		mult++
	}

	partSize = minPartSize * powInt(2, mult)
	numOfParts := archiveSize / partSize
	if archiveSize%partSize != 0 {
		numOfParts += 1
	}

	if numOfParts > maxParts {
		partSize *= 2
	}

	return partSize, nil
}

func initiateMultipart(client *glacier.Client, account string, vault string, fileName string, partSize int64) (string, error) {
	size := strconv.FormatInt(partSize, 10)
	input := glacier.InitiateMultipartUploadInput{
		AccountId:          &account,
		VaultName:          &vault,
		PartSize:           &size,
		ArchiveDescription: &fileName,
	}

	resp, err := client.InitiateMultipartUpload(context.TODO(), &input)
	if err != nil {
		return "", err
	}

	return *resp.UploadId, nil
}

func getPart(file *os.File, fileSize int64, partNumber int64, buffer []byte, partSize int64) (part, error) {
	start := partSize * (partNumber - 1)
	end := start + partSize
	if fileSize-partNumber*partSize < 0 {
		size := fileSize - start
		end = start + size
	}

	n, err := file.Read(buffer)
	reader := bytes.NewReader(buffer[:n])
	hashes := computeHashes(buffer[:n])
	treeHash := computeTreeHash(hashes)
	sha256Str := fmt.Sprintf("%x", treeHash)

	return part{reader, sha256Str, fmt.Sprintf("bytes %d-%d/*", start, end-1)}, err
}

func uploadPart(client *glacier.Client, account string, vault string, uploadId string, part part) ([]byte, error) {
	input := glacier.UploadMultipartPartInput{
		AccountId: &account,
		Body:      part.reader,
		VaultName: &vault,
		UploadId:  &uploadId,
		Range:     &part.contentRange,
	}

	r, err := client.UploadMultipartPart(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	if *r.Checksum != part.sha256 {
		return nil, errors.New("Checksums mismatch")
	}

	return hex.DecodeString(*r.Checksum)
}

func abortMultipart(client *glacier.Client, account string, vault string, uploadId string) error {
	input := glacier.AbortMultipartUploadInput{
		AccountId: &account,
		VaultName: &vault,
		UploadId:  &uploadId,
	}

	_, err := client.AbortMultipartUpload(context.TODO(), &input)
	return err
}

func completeMultipart(client *glacier.Client, account string, vault string, uploadId string, size int64, sha256 string) error {
	sizeStr := strconv.FormatInt(size, 10)
	input := glacier.CompleteMultipartUploadInput{
		AccountId:   &account,
		VaultName:   &vault,
		UploadId:    &uploadId,
		ArchiveSize: &sizeStr,
		Checksum:    &sha256,
	}

	_, err := client.CompleteMultipartUpload(context.TODO(), &input)
	return err
}

func uploadMultipart(client *glacier.Client, account string, vault string, file *os.File, fileSize int64) error {
	partNumber := int64(1)
	partSize, err := computePartSize(fileSize)
	if err != nil {
		return err
	}

	dBuffer := make([]byte, partSize)

	uploadId, err := initiateMultipart(client, account, vault, file.Name(), partSize)
	if err != nil {
		return err
	}

	nPartsToSend := int64(fileSize / partSize)
	if fileSize%partSize != 0 {
		nPartsToSend += 1
	}

	checksums := make([][]byte, nPartsToSend)

	for i := int64(0); i < nPartsToSend; i += 1 {
		part, err := getPart(file, fileSize, partNumber, dBuffer, partSize)
		if err != nil && err != io.EOF {
			return err
		}

		hash, uploadErr := uploadPart(client, account, vault, uploadId, part)
		if uploadErr != nil {
			if err := abortMultipart(client, account, vault, uploadId); err != nil {
				return err
			}

			return uploadErr
		}

		checksums[i] = hash
		partNumber++

		if err == io.EOF {
			// TODO: check if it is necessary
			break
		}
	}

	if err != nil {
		if err := abortMultipart(client, account, vault, uploadId); err != nil {
			return err
		}
		return err
	}

	treeHash := computeTreeHash(checksums)
	if err := completeMultipart(client, account, vault, uploadId, fileSize, fmt.Sprintf("%x", treeHash)); err != nil {
		if err := abortMultipart(client, account, vault, uploadId); err != nil {
			return err
		}
		return err
	}
	return nil
}

func UploadData(cfg aws.Config, account string, vault string, archive string) error {
	stat, err := os.Stat(archive)
	if err != nil {
		return err
	}

	size := stat.Size()

	file, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer file.Close()

	client := glacier.NewFromConfig(cfg)

	if size >= multipartThreshold {
		return uploadMultipart(client, account, vault, file, size)
	} else {
		return upload(client, account, vault, file)
	}
}
