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
	DimIsEnable      bool        `default:"false" yaml:"dim_is_enable"`
	BaselineIsEnable bool        `default:"false" yaml:"baseline_is_enable"`
	CoreFunc         CoreFunc    `yaml:"core_func"`
	MonitorFunc      MonitorFunc `yaml:"monitor_func"`
	ProtectedProcess []string    `yaml:"protected_process"`
	ProtectedModules []string    `yaml:"protected_modules"`
	ProtectKernel    bool        `default:"false" yaml:"protect_kernel"`
}

type CoreFunc struct {
	MeasureLogCapacity uint64 `default:"100000" yaml:"measure_log_capacity"`
	MeasureHash        string `default:"sha256" yaml:"measure_hash"`
	MeasurePcr         uint16 `default:"0" yaml:"measure_pcr"`
	MeasureSchedule    uint64 `default:"0" yaml:"measure_schedule"`
	MeasureInterval    uint64 `default:"0" yaml:"measure_interval"`
	Signature          bool   `default:"false" yaml:"signature"`
}

type MonitorFunc struct {
	ModuleType         string `default:"default" yaml:"module_type"`
	MeasureLogCapacity uint64 `default:"100000" yaml:"measure_log_capacity"`
	MeasureHash        string `default:"sha256" yaml:"measure_hash"`
	MeasurePcr         uint16 `default:"0" yaml:"measure_pcr"`
}

func init() {
	dim := Dim_yaml{}
	dim.templatePath = "/usr/share/secpaver/scripts/sec_conf/gen/gen_dim"
	dim.templatePathCheck = "/usr/share/secpaver/scripts/sec_conf/check/check_dim"
	dim.CoreFunc.MeasureHash = "sha256"
	dim.MonitorFunc.MeasureHash = "sha256"
	dim.MonitorFunc.MeasureLogCapacity = 100000
	dim.CoreFunc.MeasureLogCapacity = 100000
	var config GeneratorInterface = &dim
	AppendSecConfig(config, "dim_comm", "dim_comm")
}

func (config *Dim_yaml) checkTheValidity() bool {
	dim := config.Dim
	coreHash := false
	monitorHash := false
	hashSupport := []string{"sha256", "sm2"}

	if dim.DimIsEnable == false {
		return true
	}
	if dim.CoreFunc.MeasureLogCapacity < 100 ||
		dim.CoreFunc.MeasurePcr > 128 ||
		dim.CoreFunc.MeasureSchedule > 1000 ||
		dim.CoreFunc.MeasureInterval > 525600 ||
		dim.MonitorFunc.MeasureLogCapacity < 100 ||
		dim.MonitorFunc.MeasurePcr > 128 {
		return false
	}
	if len(dim.ProtectedProcess) > 1000 || len(dim.ProtectedModules) > 1000 {
		fmt.Println("The policy setting is too long. more than 1000 line")
		return false
	}
	for i := 0; i < len(hashSupport); i++ {
		if hashSupport[i] == dim.CoreFunc.MeasureHash {
			coreHash = true
		}
		if hashSupport[i] == dim.MonitorFunc.MeasureHash {
			monitorHash = true
		}
	}
	if monitorHash != true || coreHash != true {
		return false
	}
	if len(dim.ProtectedProcess) > 0 {
		for _, file := range dim.ProtectedProcess {
			if file == "" {
				fmt.Println("Protected_process is empty")
				return false
			}
		}
	}
	if len(dim.ProtectedModules) > 0 {
		for _, file := range dim.ProtectedModules {
			if file == "" {
				fmt.Println("protected_modules is empty")
				return false
			}
		}
	}
	if len(dim.ProtectedModules) == 0 &&
		len(dim.ProtectedProcess) == 0 &&
		dim.ProtectKernel == false {
		fmt.Println("please configure the dim policy")
		return false
	}
	return true
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
