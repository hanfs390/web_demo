package readConf

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"web_demo/global"
)

type PictureNode struct {
	FileName string
	Show string /* the name for show */
	Url string
	Number int
	Label string
}
type PictureDir struct {
	Dir       string
	Label     []string
	PictureList []PictureNode
}
var PictureData []PictureDir
func md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func getFileNumberFromDir(path string) int {
	var count int
	rd, _ := ioutil.ReadDir(path)
	for _, fi := range rd {
		if !fi.IsDir() {
			if strings.HasSuffix(fi.Name(), ".jpg") {
				count++
			}
		}
	}
	return count
}
func readPictureLabel(path string, dirNode *PictureDir, label string) {
	hexPicture := strings.Replace(path, global.PictureDir, global.RoutePictureDir, 1)
	rd, _ := ioutil.ReadDir(path) //read all file in this path
	for _, fi := range rd {
		if fi.IsDir() {
			var name string
			pictureNode := PictureNode{}
			pictureNode.Label = label
			if len(fi.Name()) > 60 {
				os.Remove(path+"/"+fi.Name())
				continue
			} else {
				name = fi.Name()
			}
			pictureNode.FileName = fi.Name()
			pictureNode.Url = hexPicture +"/" +fi.Name()
			pictureNode.Show = name
			pictureNode.Number = getFileNumberFromDir(path + "/" + name)
			dirNode.PictureList = append(dirNode.PictureList, pictureNode)
		}
	}
}
func readPictureDir(path string, dirNode *PictureDir) {
	rd, _ := ioutil.ReadDir(path) //read all file and dir in this path
	for _, fi := range rd {
		if fi.IsDir() {
			/* this is label */
			label := fi.Name()
			dirNode.Label = append(dirNode.Label, label)
			readPictureLabel(path + "/" + fi.Name(), dirNode, fi.Name())
		}
	}
}
func readAllPictureFiles(path string) {
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
			tmpDir := PictureDir{}
			tmpDir.Dir = fi.Name()
			readPictureDir(path + "/" + fi.Name(), &tmpDir)
/*			for i:=0; i < len(tmpDir.Label);i++ {
				fmt.Println("label", tmpDir.Label[i])
			}
			for i:=0; i < len(tmpDir.PictureList);i++ {
				fmt.Println("picture name", tmpDir.PictureList[i].FileName)
				fmt.Println("picture number", tmpDir.PictureList[i].Number)
			}*/
			PictureData = append(PictureData, tmpDir)
		}
	}
}
