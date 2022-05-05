package main

import (
	"errors"
	"log"
	"path/filepath"
)

type Memory struct {
	indexParts  [memoryBlocksCount]int64
	memoryParts [memoryBlocksCount]*[block_size]byte
}

const memoryBlocksCount = 100
const block_size = 1024 * 1024 // 1MB
const endFile = -1
const emptyBlock = 0

func (m *Memory) findFreeSpace(needBlocks int) (blocksNumber []int64, err error) {
	for x := 1; x < memoryBlocksCount; x++ {
		if needBlocks == 0 {
			return blocksNumber, nil
		}

		if m.indexParts[x] == emptyBlock {
			blocksNumber = append(blocksNumber, int64(x))
			needBlocks--
		}
	}

	return nil, errors.New("not enough free space")
}

func (m *Memory) saveToMemory(blocksIndexNumber []int64, infoBlocks [][block_size]byte) error {
	if len(blocksIndexNumber) != len(infoBlocks) {
		err := errors.New("wrong input data")
		log.Println(err)
		return err
	}

	for i := 0; i < len(blocksIndexNumber)-1; i++ {
		m.indexParts[blocksIndexNumber[i]] = blocksIndexNumber[i+1]
	}
	m.indexParts[blocksIndexNumber[len(blocksIndexNumber)-1]] = endFile

	for i, block := range infoBlocks {
		index := blocksIndexNumber[i]
		m.memoryParts[index] = &block
	}

	return nil
}

func (m *Memory) loadFromMemory(blocksIndexNumber []int64) [][block_size]byte {
	output := make([][block_size]byte, len(blocksIndexNumber))
	for i, i2 := range blocksIndexNumber {
		output[i] = *m.memoryParts[i2]
	}

	return output
}

func (m *Memory) SaveFile(fileName string, fileBody []byte) (FileInfo, error) {
	var amountOfBlocks int
	var bytesRemove int
	var extraBlock int

	isEven := (len(fileBody) % block_size) == 0
	if !isEven {
		extraBlock = 1
		bytesRemove = block_size - len(fileBody)%block_size
	}

	amountOfBlocks += len(fileBody) / block_size
	splitFileBody := make([][block_size]byte, amountOfBlocks+extraBlock)

	for i := 0; i < amountOfBlocks; i++ {
		buff := [block_size]byte{}
		for j := 0; j < block_size; j++ {
			buff[j] = fileBody[(block_size*i)+j]
		}

		splitFileBody[i] = buff
	}

	buff := [block_size]byte{}
	startIndex := block_size * amountOfBlocks
	for j := startIndex; j < len(fileBody); j++ {
		buff[j-startIndex] = fileBody[j]
	}

	if !isEven {
		splitFileBody[len(splitFileBody)-1] = buff
	}

	blocksIndex, err := m.findFreeSpace(amountOfBlocks + extraBlock)
	if err != nil {
		log.Println(err)
		return FileInfo{}, err
	}

	err = m.saveToMemory(blocksIndex, splitFileBody)
	if err != nil {
		log.Println(err)
		return FileInfo{}, err
	}

	name := filepath.Base(fileName)
	format := filepath.Ext(fileName)
	/*out := s.mem.loadFromMemory(blocksIndex)

	fmt.Println(string(out[0][:block_size-bytesRemove]), "s")*/
	return FileInfo{
		fileStart:   int(blocksIndex[0]),
		fileName:    name,
		fileFormat:  format,
		BytesRemove: bytesRemove,
	}, nil
}
