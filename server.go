package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/music", getFileNames)
	r.POST("/upload", upFile)

	r.Run()
}

type MusicFile struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func getFileNames(c *gin.Context) {
	dir := "./music"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err})
	}

	var MF []MusicFile
	for _, file := range files {
		MF = append(MF, MusicFile{Path:dir, Name:file.Name()} )
	}
	out, _ := json.Marshal(MF)
	if len(files) <= 0 {
		c.JSON(http.StatusOK, gin.H{"response":"No music file"})
	}else{
		c.JSON(http.StatusOK, string(out))
	}
}

func upFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	for _, file := range files {
		log.Println(file.Filename)
		err := c.SaveUploadedFile(file, "music/"+file.Filename)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(c.PostForm("key"))
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}