package main

import (
	"os"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	r := gin.Default()

	r.GET("/music", getFileNames)
	r.POST("/upload", upFile)
	r.GET("/start", startSound)

	r.Run()
}

func startSound(c *gin.Context) {
	path := c.Query("path")
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
	done := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))
	_ = <-done
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