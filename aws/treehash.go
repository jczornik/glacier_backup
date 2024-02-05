package aws

import (
	"crypto/sha256"
)

const chunkSize = 1024 * 1024

func computeHashes(buffer []byte) [][]byte {
	hashCount := len(buffer) / chunkSize
	if len(buffer) % chunkSize != 0 {
		hashCount += 1
	}

	hashes := make([][]byte, hashCount)
	for i := 0; i < hashCount; i++ {
		last := (i + 1) * chunkSize
		if last > len(buffer) {
			last = len(buffer)
		}
		h := sha256.Sum256(buffer[i*chunkSize:last])
		hashes[i] = h[:]
	}

	return hashes
}

func computeTreeHash(hashes [][]byte) []byte {
	hashCount := len(hashes)
	switch hashCount {
	case 0:
		return nil
	case 1:
		return hashes[0]
	}
	leaves := make([][32]byte, hashCount)
	for i := range leaves {
		copy(leaves[i][:], hashes[i])
	}
	var (
		queue = leaves[:0]
		h256  = sha256.New()
		buf   [32]byte
	)
	for len(leaves) > 1 {
		for i := 0; i < len(leaves); i += 2 {
			if i+1 == len(leaves) {
				queue = append(queue, leaves[i])
				break
			}
			h256.Write(leaves[i][:])
			h256.Write(leaves[i+1][:])
			h256.Sum(buf[:0])
			queue = append(queue, buf)
			h256.Reset()
		}
		leaves = queue
		queue = queue[:0]
	}
	return leaves[0][:]
}
