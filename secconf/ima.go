package secconf

import (
	"fmt"
	"os"
)

type IMA_yaml struct {
	templatePath      string
	templatePathCheck string
	IMA               `yaml:"ima"`
}

type IMA struct {
	MeasureList  []string `yaml:"measure_list"`
	AppraiseList []string `yaml:"appraise_list"`
}

func init() {
	ima := IMA_yaml{}
	ima.templatePath = "/usr/share/secpaver/scripts/sec_conf/gen/gen_ima"
	ima.templatePathCheck = "/usr/share/secpaver/scripts/sec_conf/check/check_ima"
	var config GeneratorInterface = &ima
	AppendSecConfig(config, "ima_comm", "ima_comm")
}

func (config *IMA_yaml) checkTheValidity() bool {
	ima := config.IMA
	for _, file := range ima.MeasureList {
		_, err := os.Stat(file)
		if err != nil {
			fmt.Printf("file %s does not exist!\n", file)
			return false
		}
	}
	for _, file := range ima.AppraiseList {
		_, err := os.Stat(file)
		if err != nil {
			fmt.Printf("file %s does not exist!\n", file)
			return false
		}
	}
	return true

}

func (config *IMA_yaml) initConfigDev(filePath string) bool {
	return InitConfigDev(filePath, config)
}

func (config *IMA_yaml) genGeneratedCommands(file *os.File) bool {
	if config.checkTheValidity() == false {
		fmt.Println("Invalid ima parameter.")
		return false
	}

	return GenGeneratedCommands(config.templatePath, file, config)
}

func (config *IMA_yaml) genCheckCommands(file *os.File) bool {
	if config.checkTheValidity() == false {
		fmt.Println("Invalid ima parameter.")
		return false
	}

	return GenGeneratedCommands(config.templatePathCheck, file, config)
}
