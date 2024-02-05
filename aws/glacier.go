package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
)

const (
	multipartThreshold = 1024 * 1024 * 100
	partSize = 1024 * 1024 * 4
)

type part struct {
	reader io.Reader
	sha256 string
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
	input := glacier.UploadArchiveInput{
		AccountId: &account,
		VaultName: &vault,
		Body:      file,
	}

	_, err := client.UploadArchive(context.TODO(), &input)

	return err
}

func initiateMultipart(client *glacier.Client, account string, vault string) (string, error) {
	size := strconv.Itoa(partSize)
	input := glacier.InitiateMultipartUploadInput{
		AccountId: &account,
		VaultName: &vault,
		PartSize:  &size,
	}

	resp, err := client.InitiateMultipartUpload(context.TODO(), &input)
	if err != nil {
		return "", err
	}

	return *resp.UploadId, nil
}

func getPart(file *os.File, fileSize int64, partNumber int64, dbuffer []byte) (part, error) {
	buffer := &dbuffer
	start := partSize * (partNumber - 1)
	fmt.Printf("Start: %d\n", start)
	end := start + partSize
	fmt.Printf("End: %d\n", end)
	if fileSize - partNumber * partSize < 0 {
		size := fileSize - start
		fmt.Printf("fileSize: %d\n", fileSize)
		fmt.Printf("Size: %d\n", size)
		b := make([]byte, size)
		buffer = &b
		end = start + size
	}

	_, err := file.Read(*buffer)
	reader := bytes.NewReader(*buffer)
	sha256 := sha256.New()
	sha256.Write(*buffer)
	sha256Str := fmt.Sprintf("%x", sha256.Sum(nil))

	return part{reader, sha256Str, fmt.Sprintf("bytes %d-%d/*", start, end - 1)}, err
}

func uploadPart(client *glacier.Client, account string, vault string, uploadId string, part part) error {
	log.Println("Uploading part")
	log.Println("Part SHA256: ", part.sha256)
	input := glacier.UploadMultipartPartInput{
		AccountId: &account,
		Body:      part.reader,
		VaultName: &vault,
		UploadId:  &uploadId,
		// Checksum:  &part.sha256,
		Range:     &part.contentRange,
	}

	_, err := client.UploadMultipartPart(context.TODO(), &input)
	return err
}

func abortMultipart(client *glacier.Client, account string, vault string, uploadId string) error {
	log.Println("Aborting multipart")
	input := glacier.AbortMultipartUploadInput{
		AccountId: &account,
		VaultName: &vault,
		UploadId:  &uploadId,
	}

	_, err := client.AbortMultipartUpload(context.TODO(), &input)
	return err
}

func completeMultipart(client *glacier.Client, account string, vault string, uploadId string, size int64, sha256 string) error {
	log.Println("Completing multipart")
	sizeStr := strconv.Itoa(int(size))
	input := glacier.CompleteMultipartUploadInput{
		AccountId: &account,
		VaultName: &vault,
		UploadId:  &uploadId,
		ArchiveSize: &sizeStr,
		Checksum: &sha256,
	}

	_, err := client.CompleteMultipartUpload(context.TODO(), &input)
	return err
}

func calculateHash(file *os.File) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	hashStr := fmt.Sprintf("%x", hash.Sum(nil))
	log.Println("File Hash: ", hashStr)
	return hashStr, nil
}


func uploadMultipart(client *glacier.Client, account string, vault string, file *os.File, fileSize int64) error {
	log.Println("Uploading multipart")
	partNumber := int64(1)
	dBuffer := make([]byte, partSize)

	uploadId, err := initiateMultipart(client, account, vault)
	if err != nil {
		return err
	}

	log.Printf("Upload ID: %s\n", uploadId)

	for sent := int64(0); sent < fileSize; sent += partSize {
		part, err := getPart(file, fileSize, partNumber, dBuffer)
		if err != nil && err != io.EOF {
			return err
		}

		uploadErr := uploadPart(client, account, vault, uploadId, part)
		if uploadErr != nil {
			if err := abortMultipart(client, account, vault, uploadId); err != nil {
				return err
			}

			return uploadErr
		}

		partNumber++

		if err == io.EOF {
			// TODO: check if it is necessary
			break
		}
	}

	hash, err := calculateHash(file)
	if err != nil {
		log.Println("Error while calculating hash")
		if err := abortMultipart(client, account, vault, uploadId); err != nil {
			return err
		}
		return err
	}

	if err := completeMultipart(client, account, vault, uploadId, fileSize, hash); err != nil {
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
