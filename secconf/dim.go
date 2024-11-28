package secconf

import (
	"fmt"
	"os"
)

type Dim_yaml struct {
	templatePath      string
	templatePathCheck string
	Dim               `yaml:"dim"`
}

type Dim struct {
	DimIsEnable        bool     `default:"false" yaml:"enable"`
	MeasureList        []string `yaml:"measure_list"`
	BaselineIsEnable   bool     `default:"false" yaml:"auto_baseline"`
	MeasureLogCapacity uint32   `default:"100000" yaml:"log_cap"`
	MeasureHash        string   `default:"sha256" yaml:"hash"`
	CorePcr            uint16   `default:"0" yaml:"core_pcr"`
	MonitorPcr         uint16   `default:"0" yaml:"monitor_pcr"`
	MeasureSchedule    uint32   `default:"0" yaml:"schedule"`
	MeasureInterval    uint32   `default:"0" yaml:"interval"`
	Signature          bool     `default:"false" yaml:"signature"`
}

func init() {
	dim := Dim_yaml{}
	dim.templatePath = "/usr/share/secpaver/scripts/sec_conf/gen/gen_dim"
	dim.templatePathCheck = "/usr/share/secpaver/scripts/sec_conf/check/check_dim"
	dim.MeasureHash = "sha256"
	dim.MeasureLogCapacity = 100000
	var config GeneratorInterface = &dim
	AppendSecConfig(config, "dim_comm", "dim_comm")
}

func (config *Dim_yaml) checkTheValidity() bool {
	dim := config.Dim
	checkHash := false
	hashSupport := []string{"sha256", "sm3"}

	if dim.DimIsEnable == false {
		return true
	}

	if dim.MeasureLogCapacity < 100 ||
		dim.CorePcr > 128 ||
		dim.MeasureSchedule > 1000 ||
		dim.MeasureInterval > 525600 ||
		dim.MonitorPcr > 128 {
		return false
	}

	if len(dim.MeasureList) > 1000 {
		fmt.Println("The policy setting is too long. more than 1000 line")
		return false
	}

	if len(dim.MeasureList) > 0 {
		for _, file := range dim.MeasureList {
			if file == "" {
				fmt.Println("MeasureList is empty")
				return false
			}
		}
	} else {
		fmt.Println("please configure the dim policy")
		return false
	}

	for i := 0; i < len(hashSupport); i++ {
		if hashSupport[i] == dim.MeasureHash {
			checkHash = true
		}
	}

	return checkHash
}

func (config *Dim_yaml) initConfigDev(filePath string) bool {
	return InitConfigDev(filePath, config)
}

func (config *Dim_yaml) genGeneratedCommands(file *os.File) bool {
	if config.checkTheValidity() == false {
		fmt.Println("Invalid dim parameter.")
		return false
	}

	return GenGeneratedCommands(config.templatePath, file, config)
}

func (config *Dim_yaml) genCheckCommands(file *os.File) bool {
	if config.checkTheValidity() == false {
		fmt.Println("Invalid dim parameter.")
		return false
	}

	return GenGeneratedCommands(config.templatePathCheck, file, config)
}
