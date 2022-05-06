package server

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"splitFiles/controller"
)

type Handlers struct {
	controller controller.Controller
}

func (h *Handlers) GetSavePage(c *gin.Context) {
	ts, err := template.ParseFiles("data/loadFile.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(c.Writer, nil)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "Internal Server Error", 500)
		return
	}
}

func (h *Handlers) SaveFile(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := c.Request.FormFile("datafile")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	buff := make([]byte, fileHeader.Size)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		log.Println(err)
		http.Error(c.Writer, "", http.StatusInternalServerError)
		return
	}

	_, err = h.controller.SaveFile(fileHeader.Filename, buff)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) LoadFile(c *gin.Context) {
	fileName := c.Param("filename")
	fileBytes, err := h.controller.LoadFile(fileName)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "", http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Write(fileBytes)
	return
}
