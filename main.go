package main

import (
	"splitFiles/controller"
	memory "splitFiles/memory"
	"splitFiles/server"
)

func main() {
	mem := memory.Memory{}
	c := controller.Controller{}
	c.InitController(&mem)

	server.StartServer(c)
}
