package secconf

import (
	"os"
)

type SecureBootYaml struct {
	templatePath      string
	templatePathCheck string
	SecureBoot        `yaml:"secure_boot"`
}

type SecureBoot struct {
	SecureBootIsEnable bool `default:"false" yaml:"secure_boot_is_enable"`
	AntiRollback       bool `default:"false" yaml:"anti_rollback"`
	Verbose            bool `default:"false" yaml:"verbose"`
}

func init() {
	secureBoot := SecureBootYaml{}
	secureBoot.templatePath = "/usr/share/secpaver/scripts/sec_conf/gen/gen_secure_boot"
	secureBoot.templatePathCheck = "/usr/share/secpaver/scripts/sec_conf/check/check_secure_boot"
	var config GeneratorInterface = &secureBoot
	AppendSecConfig(config, "secure_boot_comm", "secure_boot_comm")
}

func (config *SecureBootYaml) initConfigDev(filePath string) bool {
	return InitConfigDev(filePath, config)
}

func (config *SecureBootYaml) genGeneratedCommands(file *os.File) bool {
	return GenGeneratedCommands(config.templatePath, file, config)
}

func (config *SecureBootYaml) genCheckCommands(file *os.File) bool {
	return GenGeneratedCommands(config.templatePathCheck, file, config)
}
