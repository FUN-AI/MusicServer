package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	var musicListName string

	r.GET("/musicFileNameAll", getFileNames)
	r.POST("/upload", upMusicFile)
	r.GET("/start/music", func(c *gin.Context) {
		path := c.Query("path")
		startSound(path)
		c.JSON(http.StatusOK, gin.H{
			"res": "started music",
		})
	})
	r.GET("/start/playList", func(c *gin.Context) {
		listName := c.Query("name")
		if listName != musicListName {
			musicPointer = 0
			musicListName = listName
		}
		musicPointer = startPlayList(listName)
		c.JSON(http.StatusOK, gin.H{
			"res": "started play list",
		})
	})
	r.GET("/next", func(c *gin.Context) {
		startPlayList(musicListName)
		c.JSON(http.StatusOK, gin.H{
			"res": "started play list",
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
	r.POST("/makePlayList", func(c *gin.Context) {
		listName := c.Query("name")
		strMusicList := c.Query("musics")
		var MusicList []MusicFile
		for _, i := range strings.Split(strMusicList, ",") {
			MusicList = append(MusicList, MusicFile{"./music/", i})
		}
		makePlaylist(listName, MusicList)
		c.JSON(http.StatusOK, gin.H{
			"res": "made play list",
		})
	})
	r.GET("/playListNameAll", func(c *gin.Context) {
		out := readPlayList()
		var names []string
		for _, n := range out {
			names = append(names, n.ListName)
		}
		c.JSON(http.StatusOK, names)
	})
	r.GET("/playListDetail", func(c *gin.Context) {
		listName := c.Query("name")
		var aPlayList PlayList
		for _, n := range readPlayList() {
			if n.ListName == listName {
				aPlayList = n
				break
			}
		}
		jsonPL, _ := json.Marshal(aPlayList)
		c.JSON(http.StatusOK, string(jsonPL))
	})

	r.Run()
}

type MusicFile struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type PlayList struct {
	ListName string      `json:"listName"`
	Musics   []MusicFile `json:"musics"`
}

func startPlayList(listName string) {
	var MusicFiles []MusicFile
	for _, n := range readPlayList() {
		if n.ListName == listName {
			MusicFiles = n.Musics
			break
		}
	}
	var stream beep.Streamer
	for i, n := range MusicFiles {
		f, _ := os.Open(n.Path + n.Name)
		if i == 0 {
			s, format, _ := mp3.Decode(f)
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			stream = beep.Seq(s)
		} else {
			s, _, _ := mp3.Decode(f)
			stream = beep.Seq(s)
		}
	}
	speaker.Play(stream)
}

func nextPlayList(listName string) {
	var MusicFiles []MusicFile
	for _, n := range readPlayList() {
		if n.ListName == listName {
			MusicFiles = n.Musics
			break
		}
	}
	var stream beep.Streamer
	for i, n := range MusicFiles {
	}
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

func makePlaylist(listName string, musicList []MusicFile) {
	list := PlayList{listName, musicList}
	lists := readPlayList()
	lists = append(lists, list)
	out, err := json.MarshalIndent(lists, "", "    ")
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile("PlayList.json", out, os.ModePerm)
}

func readPlayList() []PlayList {
	raw, err := ioutil.ReadFile("./PlayList.json")
	if err != nil {
		log.Println(err.Error())
	}
	var lists []PlayList
	json.Unmarshal(raw, &lists)
	return lists
}

func upMusicFile(c *gin.Context) {
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
