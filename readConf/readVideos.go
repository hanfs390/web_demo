package readConf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"web_demo/global"
)

type VideoNode struct {
	Label string
	Url string
	Picture string
	FileName string
	Size uint64
	Desc string
}
type VideoDir struct {
	Dir       string
	Label     []string
	VideoList []VideoNode
}
var VideoData []VideoDir
func readVideoLabel(path string, dirNode *VideoDir, label string) {
	hexVideo := strings.Replace(path, global.VideoDir, global.RouteVideoDir, 1)
	rd, _ := ioutil.ReadDir(path) //read all file in this path
	for _, fi := range rd {
		if !fi.IsDir() {
			if !strings.HasSuffix(fi.Name(), ".mp4") {
				fmt.Println(fi.Name())
				continue
			}
			video := VideoNode{}
			video.Label = label
			video.FileName = strings.Replace(fi.Name(), ".mp4", "", 1)
			if len(video.FileName) > 40 {
				video.FileName = video.FileName[:40]
			}
			video.Url = hexVideo + "/" + fi.Name()
			picSrc := strings.Replace(fi.Name(), ".mp4", ".jpg", 1)
			video.Picture = global.RouteVideoDir + "/pictures/" + picSrc
			dirNode.VideoList = append(dirNode.VideoList, video)
		}
	}
}

func readVideoDir(path string, dirNode *VideoDir) {
	rd, _ := ioutil.ReadDir(path) //read all file and dir in this path
	for _, fi := range rd {
		if fi.IsDir() {
			/* this is label */
			label := fi.Name()
			dirNode.Label = append(dirNode.Label, label)
			readVideoLabel(path + "/" + fi.Name(), dirNode, fi.Name())
		}
	}
}
func readAllVideoFiles(path string) {
	fi, e := os.Stat(path)
	if e != nil {
		return
	}
	if !fi.IsDir() {
		fmt.Println(path, " is not a dir")
		return
	}
	rd, _ := ioutil.ReadDir(path) //read all file and dir in this path
	for _, fi := range rd {
		if fi.IsDir() {
			if fi.Name() == "pictures" {
				continue
			}
			tmpDir := VideoDir{}
			tmpDir.Dir = fi.Name()
			readVideoDir(path + "/" + fi.Name(), &tmpDir)/*
			for i:=0; i < len(tmpDir.Label);i++ {
				fmt.Println("label", tmpDir.Label[i])
			}
			for i:=0; i < len(tmpDir.VideoList);i++ {
				fmt.Println("video", tmpDir.VideoList[i].FileName)
			}*/
			VideoData = append(VideoData, tmpDir)
		}
	}
}
