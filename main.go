package main

import (
	"io/ioutil"
	"os"
	"splitFiles/controller"
	memory "splitFiles/memory"
)

func main() {
	bytes1, err := ioutil.ReadFile("1.docx")
	if err != nil {
		return
	}

	mem := memory.Memory{}
	c := controller.Controller{}
	c.SetMemory(&mem)

	info, _ := c.SaveFile("s", bytes1)
	out, _ := c.LoadFile(info)

	file, err := os.Create("12.docx")
	if err != nil {
		return
	}
	file.Write(out)
}
