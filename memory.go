package main

import (
	"errors"
	"log"
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
