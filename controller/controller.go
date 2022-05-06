package controller

import (
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
	return s.mem.SaveFile(fileName, fileBody)
}

func (c *Controller) LoadFile(info memory.FileInfo) ([]byte, error) {
	return c.mem.LoadFile(info)
}
