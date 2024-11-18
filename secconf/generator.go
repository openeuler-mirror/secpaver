package secconf

import (
	"os"
)

type GeneratorInterface interface {
	initConfigDev(filePath string) bool
	genGeneratedCommands(file *os.File) bool
	genCheckCommands(file *os.File) bool
}

type Generator struct {
	yamlPath           string
	templatePath       string
	gConfig            []GeneratorInterface
	shellConfigFunName []string
	shellCheckFunName  []string
	ShellFuns          []string
}

var generator Generator = Generator{yamlPath: "/usr/share/secpaver/scripts/sec_conf/sec_conf.yaml",
	templatePath: "/usr/share/secpaver/scripts/sec_conf/gen_comm.sh"}

func AppendSecConfig(config GeneratorInterface, shellConfigName string, shellCheckName string) {
	generator.gConfig = append(generator.gConfig, config)
	generator.shellConfigFunName = append(generator.shellConfigFunName, shellConfigName)
	generator.shellCheckFunName = append(generator.shellCheckFunName, shellCheckName)
}

func GenGeneratedAllCommands(file *os.File) bool {
	generator.ShellFuns = generator.shellConfigFunName

	for _, config := range generator.gConfig {
		if !config.initConfigDev(generator.yamlPath) {
			return false
		}
		if !config.genGeneratedCommands(file) {
			return false
		}
	}

	return GenGeneratedCommands(generator.templatePath, file, generator)
}

func GenGeneratedAllCheck(file *os.File) bool {
	generator.ShellFuns = generator.shellCheckFunName

	for _, config := range generator.gConfig {
		if !config.initConfigDev(generator.yamlPath) {
			return false
		}
		if !config.genCheckCommands(file) {
			return false
		}
	}

	return GenGeneratedCommands(generator.templatePath, file, generator)
}
