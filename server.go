package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/music", getFileNames)
	r.POST("/upload", upFile)
	r.GET("/start", func(c *gin.Context) {
		path := c.Query("path")
		startSound(path)
		c.JSON(http.StatusOK, gin.H{
			"res": "started music",
		})
	})
	r.GET("/stop", func(c *gin.Context) {
		speaker.Lock()
		c.JSON(http.StatusOK, gin.H{
			"res": "stoped music",
		})
	})
	r.GET("/restart", func(c *gin.Context) {
		speaker.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"res": "restarted music",
		})
	})

	r.Run()
}

func startSound(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	s, format, err := mp3.Decode(f)
	if err != nil {
		log.Println(err)
	}
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Println(err)
	}
	speaker.Play(beep.Seq(s))
}

type MusicFile struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func getFileNames(c *gin.Context) {
	dir := "./music"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	var MF []MusicFile
	for _, file := range files {
		MF = append(MF, MusicFile{Path: dir, Name: file.Name()})
	}
	out, _ := json.Marshal(MF)
	if len(files) <= 0 {
		c.JSON(http.StatusOK, gin.H{"res": "No music file"})
	} else {
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
	c.JSON(http.StatusOK, gin.H{
		"res": fmt.Sprintf("%d files uploaded!", len(files)),
	})
}
