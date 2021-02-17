package readConf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"web_demo/global"
)
/* txt no include label */
type TxtNode struct {
	FileName string
	Url string
	Dir string
}
type TxtDir struct {
	Dir       string
	TxtList []TxtNode
}
var TxtData []TxtDir

func readAllTxtInDir(path string, dirNode *TxtDir) {
	hexPicture := strings.Replace(path, global.TxtDir, global.RouteTxtDir, 1)
	rd, _ := ioutil.ReadDir(path) //read all file in this path
	for _, fi := range rd {
		if !fi.IsDir() {
			if !strings.HasSuffix(fi.Name(), ".txt") {
				/* not txt */
				continue;
			}
			var name = strings.Replace(fi.Name(), ".txt", "", -1)
			txtNode := TxtNode{}
			txtNode.FileName = name
			txtNode.Url = hexPicture +"/" +fi.Name()
			dirNode.TxtList = append(dirNode.TxtList, txtNode)
		}
	}
}
func readAllTxtFiles(path string) {
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
			tmpDir := TxtDir{}
			tmpDir.Dir = fi.Name()
			readAllTxtInDir(path + "/" + fi.Name(), &tmpDir)

			for i:=0; i < len(tmpDir.TxtList);i++ {
				fmt.Println("txt name", tmpDir.TxtList[i].FileName)
				fmt.Println("txt url", tmpDir.TxtList[i].Url)
			}
			TxtData = append(TxtData, tmpDir)
		}
	}
}
