package secconf

import (
	"os"
)

type ModsignYaml struct {
	templatePath      string
	templatePathCheck string
	Modsign           `yaml:"modsign"`
}

type Modsign struct {
	Enable bool `default:"false" yaml:"enable"`
}

func init() {
	modsign := ModsignYaml{}
	modsign.templatePath = "/usr/share/secpaver/scripts/sec_conf/gen/gen_modsign"
	modsign.templatePathCheck = "/usr/share/secpaver/scripts/sec_conf/check/check_modsign"
	var config GeneratorInterface = &modsign
	AppendSecConfig(config, "modsign_comm", "modsign_comm")
}

func (config *ModsignYaml) initConfigDev(filePath string) bool {
	return InitConfigDev(filePath, config)
}

func (config *ModsignYaml) genGeneratedCommands(file *os.File) bool {
	return GenGeneratedCommands(config.templatePath, file, config)
}

func (config *ModsignYaml) genCheckCommands(file *os.File) bool {
	return GenGeneratedCommands(config.templatePathCheck, file, config)
}
