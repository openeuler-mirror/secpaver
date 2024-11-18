package secconf

import (
	"os"
	"testing"
)

func TestGenGeneratedCommands(t *testing.T) {
	var file *os.File
	var err error

	type args struct {
		filePath string
		file     *os.File
		config   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"TestGenGeneratedCommands DIM", args{"gen/gen_dim", file, Dim_yaml{}}, true},
		{"TestGenGeneratedCommands IMA", args{"gen/gen_ima", file, IMA_yaml{}}, true},
		{"TestGenGeneratedCommands secure_boot", args{"gen/gen_secure_boot", file, SecureBootYaml{}}, true},
		{"TestGenGeneratedCommands modsign", args{"gen/gen_modsign", file, ModsignYaml{}}, true},
		{"TestGenGeneratedCommands DIM", args{"check/check_dim", file, Dim_yaml{}}, true},
		{"TestGenGeneratedCommands IMA", args{"check/check_ima", file, IMA_yaml{}}, true},
		{"TestGenGeneratedCommands secure_boot", args{"check/check_secure_boot", file, SecureBootYaml{}}, true},
		{"TestGenGeneratedCommands modsign", args{"check/check_modsign", file, ModsignYaml{}}, true},
	}
	for _, tt := range tests {
		file, err = createFile("test.sh")
		if err != nil {
			t.Errorf("failed to createFile")
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := GenGeneratedCommands(tt.args.filePath, tt.args.file, tt.args.config); got != tt.want {
				t.Errorf("GenGeneratedCommands() = %v, want %v", got, tt.want)
			}
		})
		file.Close()
	}
}

func TestInitConfigDev(t *testing.T) {
	type args struct {
		filePath string
		config   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"TestInitConfigDev DIM1", args{"test/test1.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev secure_boot1", args{"test/test2.yaml", &SecureBootYaml{}}, true},
		{"TestInitConfigDev modsign1", args{"test/test3.yaml", &ModsignYaml{}}, true},
		{"TestInitConfigDev 1", args{"test.yaml", &Dim_yaml{}}, false},
		{"TestInitConfigDev DIM2", args{"test/test4.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev secure_boot2", args{"test/test5.yaml", &SecureBootYaml{}}, true},
		{"TestInitConfigDev modsign2", args{"test/test6.yaml", &ModsignYaml{}}, true},
		{"TestInitConfigDev modsign3", args{"test/test7.yaml", &ModsignYaml{}}, true},
		{"TestInitConfigDev secure_boot3", args{"test/test7.yaml", &SecureBootYaml{}}, true},
		{"TestInitConfigDev DIM3", args{"test/test7.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev modsign4", args{"test/test8.yaml", &ModsignYaml{}}, true},
		{"TestInitConfigDev secure_boot4", args{"test/test8.yaml", &SecureBootYaml{}}, true},
		{"TestInitConfigDev DIM4", args{"test/test8.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev DIM5", args{"test/test9.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev DIM6", args{"test/test10.yaml", &Dim_yaml{}}, false},
		{"TestInitConfigDev DIM7", args{"test/test11.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev DIM8", args{"test/test12.yaml", &Dim_yaml{}}, true},
		{"TestInitConfigDev IMA1", args{"test/test13.yaml", &IMA_yaml{}}, true},
		{"TestInitConfigDev IMA2", args{"test/test14.yaml", &IMA_yaml{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitConfigDev(tt.args.filePath, tt.args.config); got != tt.want {
				t.Errorf("InitConfigDev() = %v, want %v", got, tt.want)
			}
		})
	}
}
