package secconf

import "testing"

func TestIMA_yaml_checkTheValidity(t *testing.T) {
	type fields struct {
		templatePath      string
		templatePathCheck string
		IMA               IMA
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{"TestDim_CheckTheValidity 1:", fields{"", "", IMA{false, true, true, []string{"/usr/bin/xc", "/usr/bin/ls"}}}, true},
		{"TestDim_CheckTheValidity 2:", fields{"", "", IMA{true, true, true, []string{"/usr/bin/xc", "/usr/bin/ls"}}}, false},
		{"TestDim_CheckTheValidity 3:", fields{"", "", IMA{true, true, true, []string{"/usr/bin/sleep", "/usr/bin/ls"}}}, true},
		{"TestDim_CheckTheValidity 4:", fields{"", "", IMA{true, true, false, []string{"/usr/bin/sleep", "/usr/bin/ls"}}}, true},
		{"TestDim_CheckTheValidity 5:", fields{"", "", IMA{true, false, false, []string{"/usr/bin/xc", "/usr/bin/ls"}}}, true},
		{"TestDim_CheckTheValidity 6:", fields{"", "", IMA{true, false, true, []string{"/usr/bin/xc", "/usr/bin/ls"}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &IMA_yaml{
				templatePath:      tt.fields.templatePath,
				templatePathCheck: tt.fields.templatePathCheck,
				IMA:               tt.fields.IMA,
			}
			if got := config.checkTheValidity(); got != tt.want {
				t.Errorf("checkTheValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}
