package secconf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"text/template"
)

func InitConfigDev(filePath string, config interface{}) bool {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("failed to read yaml, err:", err)
		return false
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		fmt.Println("failed to yaml unmarshal, err", err)
		return false
	}
	return true
}

func GenGeneratedCommands(filePath string, file *os.File, config interface{}) bool {
	tmpl, error := template.ParseFiles(filePath)
	if error != nil {
		fmt.Println("create template failed, err:", error)
		return false
	}

	if file == nil {
		error = tmpl.Execute(os.Stdout, config)
	} else {
		error = tmpl.Execute(file, config)
	}

	if error != nil {
		fmt.Println("Execute failed, err:", error)
		return false
	}

	return true
}
