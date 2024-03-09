package tools

import (
	"fmt"
	"sync"
	"os/exec"
	"bytes"
)

type splitter struct {
	processCount int
	maxBufferSize int

	buffers [][]byte
	bufferIndex int
	cmdCtor CmdConstructor
	wg sync.WaitGroup
}

type CmdConstructor func() *exec.Cmd;

func NewSplitter(processCount, maxBufferSize int, cmdCtor CmdConstructor) *splitter {
	buffers := make([][]byte, 0)
	return &splitter{processCount, maxBufferSize, buffers, 0, cmdCtor, sync.WaitGroup{}}
}

// func (s *spliter) Close() error {
// 	if s.currentFile != nil {
// 		return s.currentFile.Close()
// 	}

// 	return nil
// }

// func (s *spliter) newFile() error {
// 	var err error
// 	s.wrote = 0
// 	s.index += 1

// 	fileName := fmt.Sprintf("%s.part%d", s.filePattern, s.index)

// 	if s.currentFile != nil {
// 		err = s.currentFile.Close()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	s.currentFile, err = os.Create(fileName)
// 	return err
// }

func (s *splitter) Write(p []byte) (n int, err error) {
	toWrite := len(p)
	wrote := 0
	totalWrote := 0

	for toWrite > 0 {
		if s.maxBufferSize - wrote == 0 {
			s.bufferIndex += 1
			wrote = 0
		}

		sLeft := s.maxBufferSize - wrote
		end := sLeft

		if toWrite < sLeft {
			end = toWrite
		}

		end += totalWrote

		s.buffers = append(s.buffers, p[totalWrote:end])
		if err != nil {
			return 0, err
		}

		wrote += end - totalWrote
		toWrite -= wrote
		totalWrote += wrote


		s.wg.Add(1)
		idx := len(s.buffers) - 1
		go func() {
			defer s.wg.Done()
			cmd := s.cmdCtor()

			cmd.Stdin = bytes.NewReader(s.buffers[idx])
			cmd.Run()
		}()
	}

	if totalWrote != len(p) {
		return totalWrote, fmt.Errorf("Wrote %d, expected %d", totalWrote, len(p))
	}

	s.wg.Wait()

	return totalWrote, nil
}
