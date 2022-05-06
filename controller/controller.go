package controller

import (
	"errors"
	"log"
	"splitFiles/memory"
)

type Controller struct {
	allFiles map[string]memory.FileInfo
	mem      memory.Mem
}

func (c *Controller) SetMemory(mem memory.Mem) {
	c.mem = mem
}

func (s *Controller) SaveFile(fileName string, fileBody []byte) (memory.FileInfo, error) {
	fileInfo, err := s.mem.SaveFile(fileName, fileBody)
	if err != nil {
		log.Println(err)
		return memory.FileInfo{}, err
	}
	s.allFiles[fileName] = fileInfo

	return fileInfo, nil
}

func (c *Controller) LoadFile(fileName string) ([]byte, error) {
	info, isExist := c.allFiles[fileName]
	if !isExist {
		return []byte{}, errors.New("not found")
	}

	return c.mem.LoadFile(info)
}
