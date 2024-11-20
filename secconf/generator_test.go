package secconf

import (
	"os"
	"testing"
)

func testInit() {
	generator.gConfig = []GeneratorInterface{}
	generator.shellConfigFunName = []string{}
	generator.shellCheckFunName = []string{}
	generator.ShellFuns = []string{}
	dim := Dim_yaml{}
	dim.templatePath = "gen/gen_dim"
	dim.templatePathCheck = "check/check_dim"
	dim.CoreFunc.MeasureHash = "sha256"
	dim.MonitorFunc.MeasureHash = "sha256"
	dim.MonitorFunc.MeasureLogCapacity = 100000
	dim.CoreFunc.MeasureLogCapacity = 100000
	var config GeneratorInterface = &dim
	AppendSecConfig(config, "dim_comm", "dim_comm")
	secureBoot := SecureBootYaml{}
	secureBoot.templatePath = "gen/gen_secure_boot"
	secureBoot.templatePathCheck = "check/check_secure_boot"
	config = &secureBoot
	AppendSecConfig(config, "secure_boot_comm", "secure_boot_comm")
	modsign := ModsignYaml{}
	modsign.templatePath = "gen/gen_modsign"
	modsign.templatePathCheck = "check/check_modsign"
	config = &modsign
	AppendSecConfig(config, "modsign_comm", "modsign_comm")
	ima := IMA_yaml{}
	ima.templatePath = "gen/gen_ima"
	ima.templatePathCheck = "check/check_ima"
	config = &ima
	AppendSecConfig(config, "ima_comm", "ima_comm")
}

func TestGenGeneratedAllCommands(t *testing.T) {
	type args struct {
		file     *os.File
		yamlPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"TestGenGeneratedAllCommands1", args{nil, "test/test1.yaml"}, true},
		{"TestGenGeneratedAllCommands2", args{nil, "test/test2.yaml"}, true},
		{"TestGenGeneratedAllCommands3", args{nil, "test/test3.yaml"}, true},
		{"TestGenGeneratedAllCommands4", args{nil, "test/test4.yaml"}, true},
		{"TestGenGeneratedAllCommands5", args{nil, "test/test5.yaml"}, true},
		{"TestGenGeneratedAllCommands6", args{nil, "test/test6.yaml"}, true},
		{"TestGenGeneratedAllCommands7", args{nil, "test/test7.yaml"}, true},
		{"TestGenGeneratedAllCommands8", args{nil, "test/test8.yaml"}, true},
		{"TestGenGeneratedAllCommands9", args{nil, "test/test9.yaml"}, false},
		{"TestGenGeneratedAllCommands10", args{nil, "test/test10.yaml"}, false},
		{"TestGenGeneratedAllCommands11", args{nil, "test/test11.yaml"}, false},
		{"TestGenGeneratedAllCommands12", args{nil, "test/test12.yaml"}, false},
		{"TestGenGeneratedAllCommands13", args{nil, "test/test13.yaml"}, true},
		{"TestGenGeneratedAllCommands14", args{nil, "test/test14.yaml"}, false},
		{"TestGenGeneratedAllCommands15", args{nil, "test/test15.yaml"}, true},
		{"TestGenGeneratedAllCommands16", args{nil, "test/test16.yaml"}, false},
		{"TestGenGeneratedAllCommands17", args{nil, "test/test17.yaml"}, false},
	}

	generator.templatePath = "gen_comm.sh"
	for _, tt := range tests {
		testInit()

		generator.yamlPath = tt.args.yamlPath
		t.Run(tt.name, func(t *testing.T) {
			if got := GenGeneratedAllCommands(tt.args.file); got != tt.want {
				t.Errorf("GenGeneratedAllCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenGeneratedAllCheck(t *testing.T) {
	type args struct {
		file     *os.File
		yamlPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"TestGenGeneratedAllCheck1", args{nil, "test/test1.yaml"}, true},
		{"TestGenGeneratedAllCheck2", args{nil, "test/test2.yaml"}, true},
		{"TestGenGeneratedAllCheck3", args{nil, "test/test3.yaml"}, true},
		{"TestGenGeneratedAllCheck4", args{nil, "test/test4.yaml"}, true},
		{"TestGenGeneratedAllCheck5", args{nil, "test/test5.yaml"}, true},
		{"TestGenGeneratedAllCheck6", args{nil, "test/test6.yaml"}, true},
		{"TestGenGeneratedAllCheck7", args{nil, "test/test7.yaml"}, true},
		{"TestGenGeneratedAllCheck8", args{nil, "test/test8.yaml"}, true},
		{"TestGenGeneratedAllCheck9", args{nil, "test/test9.yaml"}, false},
		{"TestGenGeneratedAllCheck13", args{nil, "test/test13.yaml"}, true},
		{"TestGenGeneratedAllCheck10", args{nil, "test/test10.yaml"}, false},
		{"TestGenGeneratedAllCheck11", args{nil, "test/test11.yaml"}, false},
		{"TestGenGeneratedAllCheck12", args{nil, "test/test12.yaml"}, false},
		{"TestGenGeneratedAllCheck13", args{nil, "test/test13.yaml"}, true},
		{"TestGenGeneratedAllCheck14", args{nil, "test/test14.yaml"}, false},
		{"TestGenGeneratedAllCheck15", args{nil, "test/test15.yaml"}, true},
		{"TestGenGeneratedAllCheck16", args{nil, "test/test16.yaml"}, false},
		{"TestGenGeneratedAllCheck17", args{nil, "test/test17.yaml"}, false},
	}

	generator.templatePath = "gen_comm.sh"
	for _, tt := range tests {
		testInit()

		generator.yamlPath = tt.args.yamlPath
		t.Run(tt.name, func(t *testing.T) {
			generator.yamlPath = tt.args.yamlPath
			if got := GenGeneratedAllCheck(tt.args.file); got != tt.want {
				t.Errorf("GenGeneratedAllCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
