package aws

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
)

type multipartUploader struct {
	client *glacier.Client
	account string
	vault string
	uploadId string
	maxUploads int
	awsConfig aws.Config

	activeUploads atomic.Int32
	results       []uploadRes
	resCh         chan uploadRes
	endCh         chan bool
	wg	      sync.WaitGroup
	inited        bool
}

func newMultipartUploader(client *glacier.Client, account string, vault string, uploadId string, maxUploads int, awsCfg aws.Config) *multipartUploader {
	inited := false
	results := make([]uploadRes, 0)
	resCh := make(chan uploadRes)
	endCh := make(chan bool)

	uploader := &multipartUploader{client: client, account: account, vault: vault, uploadId: uploadId, maxUploads: maxUploads, awsConfig: awsCfg, results: results, resCh: resCh, endCh: endCh, inited: inited}
	uploader.run()
	return uploader
}

func (m *multipartUploader) run() {
	if !m.inited {
		go func() {
			for {
				select {
				case res := <-m.resCh:
					if res.err != nil {
						log.Fatal("Error while uploading part", res.err)
					}
					m.results = append(m.results, res)
					m.activeUploads.Add(-1)

				case <-m.endCh:
					close(m.resCh)
					return
				}
			}
		}()
		m.inited = true
	}
}

func (m *multipartUploader) upload(part part) {
	for m.activeUploads.Load() >= int32(m.maxUploads) {
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("active uploads", m.activeUploads.Load())
	fmt.Println("uploading part", part.partNo)

	m.activeUploads.Add(1)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		client := glacier.NewFromConfig(m.awsConfig)
		checksum, err := uploadPart(client, m.account, m.vault, m.uploadId, part)
		m.resCh <- uploadRes{checksum, err, part.partNo}
	}()
}

func (m *multipartUploader) getResults() []uploadRes {
	fmt.Println("waiting for uploads to finish")
	m.wg.Wait()
	fmt.Println("uploads finished")
	return m.results
}

func (m *multipartUploader) close() {
	m.endCh <- true
}

type uploadRes struct {
	checksum []byte
	err      error
	partNo   int64
}
