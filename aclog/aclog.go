/*
   golang log rotate example
   E-Mail : Mike_Zhang@live.com
*/

package aclog

import (
	"fmt"
	"log"
	"os"
	"ByzoroAC/conf"
	"strconv"
)

const (
	debugLevel = 0
	infoLevel = 1
	warningLevel = 2
	errorLevel = 3
)

type localConf struct {
	logFilesCount int
	logFilePath string
	maxFileLines int
	logLevel int
}

var logger *log.Logger
var logFile *os.File
var lineCount int
var logConf localConf

func doRotate(fPrefix string) {
	var preFileName string
	for j := logConf.logFilesCount; j >= 1; j-- {
		curFileName := fmt.Sprintf("%s_%d.log", fPrefix, j)

		k := j - 1
		if k == 0 {
			preFileName = fmt.Sprintf("%s.log", fPrefix)
		} else {
			preFileName = fmt.Sprintf("%s_%d.log", fPrefix, k)
		}
		_, err := os.Stat(curFileName)
		if err == nil {
			os.Remove(curFileName)
			fmt.Println("remove : ", curFileName)
		}
		_,err = os.Stat(preFileName)
		if err  == nil {
			fmt.Println("rename : ", preFileName, " => ", curFileName)
			err = os.Rename(preFileName, curFileName)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
func newLogger(fileName string) (*log.Logger, *os.File) {
	_, err := os.Stat(logConf.logFilePath)
	if err != nil {
		fmt.Println(err)
		os.Mkdir(logConf.logFilePath, os.ModePerm)
	}
	logFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("open file error!", err)
		return nil, nil
	} else {
		logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return logger, logFile
}

func Init() {
	//ready
	if conf.GlobalConf["LogFilePath"] != "" {
		logConf.logFilePath = conf.GlobalConf["LogFilePath"]
	} else {
		logConf.logFilePath = "./acLog"
	}
	if conf.GlobalConf["MaxFileLines"] != "" {
		logConf.maxFileLines, _= strconv.Atoi(conf.GlobalConf["MaxFileLines"])
	} else {
		logConf.maxFileLines = 10000
	}
	if conf.GlobalConf["LogFilesCount"] != "" {
		logConf.logFilesCount, _= strconv.Atoi(conf.GlobalConf["LogFilesCount"])
	} else {
		logConf.logFilesCount = 5
	}
	if conf.GlobalConf["LogLevel"] != "" {
		logConf.logLevel, _= strconv.Atoi(conf.GlobalConf["LogLevel"])
	} else {
		logConf.logLevel = 0
	}


	logger, logFile = newLogger(logConf.logFilePath + "/acLog.log")
	if (logger == nil) || (logFile == nil) {
		return
	}

}
func print(data string) {
	lineCount++
	if lineCount > logConf.maxFileLines {
		logFile.Close()
		doRotate(logConf.logFilePath + "/acLog")
		lineCount = 0
		os.Remove(logConf.logFilePath + "/acLog.log")
		logger,logFile = newLogger(logConf.logFilePath + "/acLog.log")
	}
	logger.Println(data)
}
func Debug(format string, a ...interface{}) {
	if debugLevel >= logConf.logLevel {
		head := "[ DEBUG ]"
		payload := fmt.Sprintf(format, a...)
		data := head + payload
		print(data)
	}
}
func Info(format string, a ...interface{}) {
	if infoLevel >= logConf.logLevel {
		head := "[ INFO ]"
		payload := fmt.Sprintf(format, a...)
		data := head + payload
		print(data)
	}
}
func Warning(format string, a ...interface{}) {
	if warningLevel >= logConf.logLevel {
		head := "[ WARNING ]"
		payload := fmt.Sprintf(format, a...)
		data := head + payload
		print(data)
	}
}
func Error(format string, a ...interface{}) {
	if errorLevel >= logConf.logLevel {
		head := "[ Error ]"
		payload := fmt.Sprintf(format, a...)
		data := head + payload
		print(data)
	}
}
func SetPrintLevel(level int) {
	logConf.logLevel = level
}