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
	IMAIsEnable      bool     `default:"false" yaml:"ima_is_enable"`
	MeasureIsEnable  bool     `default:"false" yaml:"measure_is_enable"`
	AppraiseIsEnable bool     `default:"false" yaml:"appraise_is_enable"`
	ProtectedFiles   []string `yaml:"protected_files"`
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
	if ima.IMAIsEnable && (ima.MeasureIsEnable || ima.AppraiseIsEnable) {
		for _, file := range ima.ProtectedFiles {
			_, err := os.Stat(file)
			if err != nil {
				fmt.Printf("file %s does not exist!\n\n", file)
				return false
			}
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
