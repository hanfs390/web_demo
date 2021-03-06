package readConf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"web_demo/global"
	"bufio"
	"strconv"
)
/* txt no include label */
type Chapters struct {
	Name string
	Start int
}
type TxtNode struct {
	FileName string
	Url string
	Dir string
	Ch []Chapters
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
func deletePreAndSufSpace(str string) string {
	strList := []byte(str)
	spaceCount, count := 0, len(strList)
	for i := 0; i <= len(strList)-1; i++ {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}

	strList = strList[spaceCount:]
	spaceCount, count = 0, len(strList)
	for i := count - 1; i >= 0; i-- {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}

	return string(strList[:count-spaceCount])
}
func isChapter(str string) string {
	if strings.Contains(str, "第") && strings.Contains(str, "章") {
		title := deletePreAndSufSpace(str)
		if len(title) < 50 {
			return title
		}
	}
	return ""
}
func CreateTxtChapters (node *TxtNode) error {
	url := strings.Replace(node.Url, global.RouteTxtDir, global.TxtDir, 1)
	file, err := os.Open(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	n := 1
	flag := 0
	for scanner.Scan() {
		lineText := scanner.Text()
		title := isChapter(lineText)
		if title != "" {
			flag = 1
			chapterNode := Chapters{Name:title, Start:n}
			node.Ch = append(node.Ch, chapterNode)
		}
		n += 1
	}
	if flag == 0 {
		//no chapter got
		chapterNode := Chapters{Name:"第一章", Start:1}
		node.Ch = append(node.Ch, chapterNode)
	}
	return nil
}
func BuildTxtDirHTML(node *TxtNode) string {
	context := ""
	for i := 0; i < len(node.Ch); i++ {
		context = context + `<div class="col-sm-6 col-md-4"><div class="thumbnail"><a href=/txt/chapters?name=`+ node.FileName +`&chapters=`+ strconv.Itoa(i) +` target="view_window">`+node.Ch[i].Name+`</a></div></div>`
	}
	return `<!DOCTYPE html>
<!-- saved from url=(0038)https://v3.bootcss.com/examples/theme/ -->
<html lang="zh-CN"><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="icon" href="https://v3.bootcss.com/favicon.ico">

    <title>`+ node.FileName +`</title>

    <!-- Bootstrap core CSS -->
    <link href="./css/bootstrap.min.css" rel="stylesheet">
    <!-- Bootstrap theme -->
    <link href="./css/bootstrap-theme.min.css" rel="stylesheet">
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <link href="./css/ie10-viewport-bug-workaround.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="./css/theme.css" rel="stylesheet">

    <!-- Just for debugging purposes. Don't actually copy these 2 lines! -->
    <!--[if lt IE 9]><script src="../../assets/js/ie8-responsive-file-warning.js"></script><![endif]-->
    <script src="./js/ie-emulation-modes-warning.js"></script>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="https://cdn.bootcss.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>

  <body>
    <div class="container theme-showcase" role="main">
      <div class="row">`+ context +`
      </div>
    </div> <!-- /container -->`
}
func readChapter(url string, start int, stop int) string {
	path := strings.Replace(url, global.RouteTxtDir, global.TxtDir, 1)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	n := 1
	data := ""
	for scanner.Scan() {
		lineText := scanner.Text()
		if stop == 0 {
			if n >= start {
				data += `<p>`+lineText+`</p>`
			}
		} else {
			if n >= start && n < stop {
				data += `<p>`+lineText+`</p>`
			}
		}
		n += 1
	}
	return strings.Replace(data, " ", "&nbsp;", -1)
}
func GetWordinTxt(name string, num int) string {
	data := ""
	for i := 0; i < len(TxtData); i++ {
		for j := 0; j < len(TxtData[i].TxtList); j++ {
			if name == TxtData[i].TxtList[j].FileName {
				tempNode := TxtData[i].TxtList[j]
				stop := 0
				if len(tempNode.Ch) > (num+1) {
					/*not last chapter*/
					stop = tempNode.Ch[num+1].Start
				}
				data = readChapter(tempNode.Url, tempNode.Ch[num].Start, stop)
			}
		}
	}
	return data
}