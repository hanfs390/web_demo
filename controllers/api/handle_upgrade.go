package api

import (
	"os"
	"log"
	"io/ioutil"
	"fmt"
	"strings"
	"mime/multipart"
	"errors"
	"io"
	"ByzoroAC/models/redis_driver"
	"ByzoroAC/controllers/task"
	"time"
)

type version struct {
	Version string
	FileName string
}

var versionGP530List []version
var versionGP830List []version
func checkVersionForm(fileName string, key string) string {
	if !strings.HasPrefix(fileName, "witfios") {
		return ""
	}
	if !strings.HasSuffix(fileName, ".bin") {
		return ""
	}
	if !strings.Contains(fileName, ".w") {
		return ""
	}
	s := strings.Replace(fileName, "witfios", "", -1)
	s1 := strings.Replace(s, ".bin", "", -1)
	s2 := strings.Replace(s1, key, "", -1)

	return s2
}

func isDir(fileAddr string) bool {
	s,err:=os.Stat(fileAddr)
	if err!=nil{
		log.Println(err)
		return false
	}
	return s.IsDir()
}
func saveToDir(file *multipart.FileHeader, dir string) error {
	flag := isDir(dir)
	if flag == false {
		err := os.Mkdir(dir, os.ModePerm)
		if err !=  nil {
			return err
		}
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dir + "/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
func SaveUpgradeFile(file *multipart.FileHeader) error {
	gp530 := checkVersionForm(file.Filename, ".w")
	if gp530 != "" {
		/* GP530 version */
		for i := 0; i < len(versionGP530List); i++ {
			if versionGP530List[i].Version == gp530 {
				return errors.New("Version file has been exist: "+file.Filename)
			}
		}
		err := saveToDir(file, "./static/version/GP530")
		if err != nil {
			return err
		}
		versionGP530List = append(versionGP530List, version{Version:gp530, FileName:file.Filename})
		return nil
	}
	gp830 := checkVersionForm(file.Filename, ".n")
	if gp830 != "" {
		/* GP830 version */
		for i := 0; i < len(versionGP830List); i++ {
			if versionGP830List[i].Version == gp830 {
				return errors.New("Version file has been exist: "+file.Filename)
			}
		}
		err := saveToDir(file, "./static/version/GP830")
		if err != nil {
			return err
		}
		versionGP830List = append(versionGP830List, version{Version:gp830, FileName:file.Filename})
		return nil
	}
	return errors.New("unknown version form")
}
func RemoveUpgradeFile(model string, version string) error {
	if model == "GP530" {
		for i := 0; i < len(versionGP530List); i++ {
			if versionGP530List[i].Version == version {
				err := os.Remove("./static/version/GP530/"+versionGP530List[i].FileName)
				if err != nil {
					fmt.Println(err)
				}
				versionGP530List = append(versionGP530List[:i], versionGP530List[i+1:]...)
				return nil
			}
		}
		return errors.New("No version file: "+version)
	}
	if (model == "GP630") || (model == "GP830") {
		for i := 0; i < len(versionGP830List); i++ {
			if versionGP830List[i].Version == version {
				err := os.Remove("./static/version/GP830/"+versionGP830List[i].FileName)
				if err != nil {
					fmt.Println(err)
				}
				versionGP830List = append(versionGP830List[:i], versionGP830List[i+1:]...)
				return nil
			}
		}
		return errors.New("No version file: "+version)
	}
	return errors.New("No model type: "+model)
}
func GetUpgradeFilesName(model string) []string{
	if model == "GP530" {
		var list []string
		for i := 0; i < len(versionGP530List); i++ {
			list = append(list, versionGP530List[i].Version)
		}
		return list
	}
	if (model == "GP630") || (model == "GP830") {
		var list []string
		for i := 0; i < len(versionGP830List); i++ {
			list = append(list, versionGP830List[i].Version)
		}
		return list
	}
	return nil
}
func getNewestVersion(v []version) string {
	if len(v) == 0 {
		return ""
	}
	newest := v[0].Version
	for i := 1; i < len(v); i++ {
		for j := 0; j < len(v[i].Version); j++ {
			if newest[j] < v[i].Version[j] {
				newest = v[i].Version
				continue
			}
		}
	}
	return newest
}
func GetNewestUpgradeFilesName(model string) string {
	if model == "GP530" {
		newest := getNewestVersion(versionGP530List)
		if newest == "" {
			fmt.Println("no upgrade file")
			return ""
		}
		fmt.Println("newest", newest)
		return newest
	}
	if (model == "GP630") || (model == "GP830") {
		newest := getNewestVersion(versionGP830List)
		if newest == "" {
			fmt.Println("no upgrade file")
			return ""
		}
		return newest
	}
	return ""
}
func UpgradeSingleAp(mac string, model string, version string) error {
	if model == "GP530" {
		for i := 0; i < len(versionGP530List); i++ {
			if versionGP530List[i].Version == version {
				ap, err := HandleFindApByMac(mac)
				if err != nil {
					return err
				}
				if ap.Model != model {
					return errors.New("Model and the AP model are not same")
				}
				if ap.TargetFirmwareVer == version {
					return errors.New("the Version is already Update")
				}
				value := make(map[string]interface{})
				value["target_firmware_ver"] = version
				err = HandleUpdateApByMac(mac, value)
				if err != nil {
					return err
				}
				/* update redis */
				g := make(map[string]interface{})
				g["TargetFirmwareVer"] = "./static/version/GP530/"+versionGP530List[i].FileName
				err = redis_driver.RedisDb.BatchHashSet(mac, g)
				if err != nil {
					fmt.Println(err)
					t := task.Task{Time:time.Now().Unix(), Op:"HashSet", Key:ap.Mac, Value:g}
					task.TaskRedis = append(task.TaskRedis, t)
				}
				return nil
			}
		}
		return errors.New("No version file: "+version)
	}
	if (model == "GP630") || (model == "GP830") {
		for i := 0; i < len(versionGP830List); i++ {
			if versionGP830List[i].Version == version {
				ap, err := HandleFindApByMac(mac)
				if err != nil {
					return err
				}
				if ap.Model != model {
					return errors.New("Model and the AP model are not same")
				}
				if ap.TargetFirmwareVer == version {
					return errors.New("the Version is already Update")
				}
				value := make(map[string]interface{})
				value["target_firmware_ver"] = version
				err = HandleUpdateApByMac(mac, value)
				if err != nil {
					return err
				}
				/* update redis */
				g := make(map[string]interface{})
				g["TargetFirmwareVer"] = "./static/version/GP830/"+versionGP830List[i].FileName
				err = redis_driver.RedisDb.BatchHashSet(mac, g)
				if err != nil {
					fmt.Println(err)
					t := task.Task{Time:time.Now().Unix(), Op:"HashSet", Key:ap.Mac, Value:g}
					task.TaskRedis = append(task.TaskRedis, t)
				}
				return nil
			}
		}
		return errors.New("No version file: "+version)
	}
	return errors.New("No model type: "+model)
}

func UpgradeAps() {

}
func readGP530Version(dir string) error {
	flag := isDir(dir)
	if flag == false {
		err := os.Mkdir(dir, os.ModePerm)
		if err !=  nil {
			return err
		}
	}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, fi := range rd {
		tmp := checkVersionForm(fi.Name(), ".w")
		if tmp == "" {
			continue
		}
		tmpVersion := version{Version:tmp, FileName:fi.Name()}
		fmt.Println(tmpVersion.Version)
		versionGP530List = append(versionGP530List, tmpVersion)
	}
	return err
}
func readGP830Version(dir string) error {
	flag := isDir(dir)
	if flag == false {
		err := os.Mkdir(dir, os.ModePerm)
		if err !=  nil {
			return err
		}
	}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, fi := range rd {
		tmp := checkVersionForm(fi.Name(), ".n")
		if tmp == "" {
			continue
		}
		tmpVersion := version{Version:tmp, FileName:fi.Name()}
		fmt.Println(tmpVersion.Version)
		versionGP830List = append(versionGP530List, tmpVersion)
	}
	return err
}
func ReadDIr() error {
	err := readGP530Version("./static/version/GP530")
	if err != nil {
		return err
	}
	err = readGP830Version("./static/version/GP830")
	if err != nil {
		return err
	}
	for i := 0; i < len(versionGP530List); i++ {
		fmt.Println(versionGP530List[i].FileName, versionGP530List[i].Version)
	}
	return nil
}