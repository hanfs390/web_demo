package model

import (
	"regexp"
	"strings"
	"errors"
	"fmt"
)

func CheckMacForm(mac string) (macStr string, err error){
	r, err := regexp.MatchString("[A-Fa-f0-9][A-Fa-f0-9]-[A-Fa-f0-9][A-Fa-f0-9]-[A-Fa-f0-9][A-Fa-f0-9]-[A-Fa-f0-9][A-Fa-f0-9]-[A-Fa-f0-9][A-Fa-f0-9]-[A-Fa-f0-9][A-Fa-f0-9]", mac)
	if r == true {
		macStr = strings.Replace(mac, "-", "", -1)
		return macStr, nil
	}
	r, err = regexp.MatchString("[A-Fa-f0-9][A-Fa-f0-9]:[A-Fa-f0-9][A-Fa-f0-9]:[A-Fa-f0-9][A-Fa-f0-9]:[A-Fa-f0-9][A-Fa-f0-9]:[A-Fa-f0-9][A-Fa-f0-9]:[A-Fa-f0-9][A-Fa-f0-9]", mac)
	if r == true {
		macStr = strings.Replace(mac, ":", "", -1)
		return macStr, nil
	}
	r, err = regexp.MatchString("[A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9][A-Fa-f0-9]", mac)
	if r == true {
		return mac, nil
	}
	err = errors.New("Error Mac Form")
	return "", err
}
func TransformMacStrToMacForm(macStr string) string {
	mac := macStr[:2] + ":" + macStr[2:4] + ":" + macStr[4:6] + ":" + macStr[6:8] + ":" + macStr[8:10] + ":" + macStr[10:12]
	fmt.Println(mac)
	return mac
}