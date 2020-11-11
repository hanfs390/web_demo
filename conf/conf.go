/*
	read .conf to struct GlobalConf
*/
package conf

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

var GlobalConf map[string]string

func ReadGlobalConf() error {
	GlobalConf = make(map[string]string)
	file, err := os.Open("/home/hfs/go/src/ByzoroAC/.configure")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		vaildText := strings.Replace(lineText, " ", "", -1)
		if strings.HasPrefix(vaildText, "#") {
			continue
		}
		stringArray := strings.Split(vaildText, "=")
		if len(stringArray) != 2 {
			fmt.Printf("conf format error: %s", lineText)
			continue
		}
		GlobalConf[stringArray[0]] = stringArray[1]

		fmt.Printf("AC read configure: %s=%s\n", stringArray[0], stringArray[1])
	}
	return nil
}