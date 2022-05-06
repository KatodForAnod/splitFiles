package memory

import (
	"errors"
	"log"
	"path/filepath"
	"sync"
)

type Memory struct {
	indexParts  [memoryBlocksCount]int64
	memoryParts [memoryBlocksCount]*[block_size]byte

	mt sync.Mutex
}

type FileInfo struct {
	fileStart   int
	fileName    string
	fileFormat  string
	BytesRemove int
}

const memoryBlocksCount = 100
const block_size = 1024 * 1024 // 1MB
const endFile = -1
const emptyBlock = 0
const reservedBlock = -2

func (m *Memory) findFreeSpace(needBlocks int) (blocksNumber []int64, err error) {
	m.mt.Lock()
	defer m.mt.Unlock()

	for x := 1; x < memoryBlocksCount; x++ {
		if needBlocks == 0 {
			return blocksNumber, nil
		}

		if m.indexParts[x] == emptyBlock {
			blocksNumber = append(blocksNumber, int64(x))
			needBlocks--
			m.indexParts[x] = reservedBlock
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

	for i, _ := range infoBlocks {
		index := blocksIndexNumber[i]
		m.memoryParts[index] = &infoBlocks[i]
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

func (m *Memory) findBlocksIndexNumber(start int64) (blocksIndexNumber []int64) {
	currIndex := start

	for {
		blocksIndexNumber = append(blocksIndexNumber, currIndex)
		index := m.indexParts[currIndex]
		if index == endFile {
			return blocksIndexNumber
		}

		currIndex = index
	}
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

	return FileInfo{
		fileStart:   int(blocksIndex[0]),
		fileName:    name,
		fileFormat:  format,
		BytesRemove: bytesRemove,
	}, nil
}

func (m *Memory) LoadFile(info FileInfo) ([]byte, error) {
	if info.fileStart > len(m.indexParts) || int64(info.fileStart) <= 0 {
		return []byte{}, errors.New("wrong file")
	}

	blocksIndex := m.findBlocksIndexNumber(int64(info.fileStart))
	blocks := m.loadFromMemory(blocksIndex)

	fileLen := block_size*len(blocksIndex) - info.BytesRemove
	file := make([]byte, 0, fileLen)

	for i := 0; i < len(blocks)-1; i++ {
		file = append(file, blocks[i][:block_size]...)
	}

	file = append(file, blocks[len(blocks)-1][:block_size-info.BytesRemove]...)
	return file, nil
}
