package main

import (
	"errors"
	"log"
)

type Memory struct {
	indexParts  [side][side]int64
	memoryParts [side][side]*[]byte
}

const side = 100
const block_size = 1024 * 1024 // 1MB
const freeBlockFlag = -1

/*func (m *Memory) SaveBlock(block [block_size]byte) (numberOfBlock int64, err error) {

}*/

func (m *Memory) findFreeSpace(needBlocks int) (blocksNumber []int64, err error) {
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			if needBlocks == 0 {
				return blocksNumber, nil
			}

			if m.indexParts[y][x] == freeBlockFlag {
				blockNumber := int64((y + 1) * (x + 1))
				blocksNumber = append(blocksNumber, blockNumber)
				needBlocks--
			}
		}
	}

	return nil, errors.New("not enough free space")
}

func (m *Memory) saveToMemory(blocksNumber []int64, infoBlocks [][]byte) error {
	if len(blocksNumber) != len(infoBlocks) {
		err := errors.New("wrong input data")
		log.Println(err)
		return err
	}

	for i, i2 := range blocksNumber {

	}
}
