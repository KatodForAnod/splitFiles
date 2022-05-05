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
const freeBlockFlag = -1

func (m *Memory) findFreeSpace(needBlocks int) (blocksNumber []int64, err error) {
	for x := 0; x < memoryBlocksCount; x++ {
		if needBlocks == 0 {
			return blocksNumber, nil
		}

		if m.indexParts[x] == freeBlockFlag {
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

	for i, block := range infoBlocks {
		index := blocksIndexNumber[i]
		m.memoryParts[index] = &block
	}

	return nil
}
