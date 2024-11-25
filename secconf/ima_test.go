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
		{"TestIMA_CheckTheValidity 1:", fields{"", "", IMA{[]string{"/usr/bin/ls"}, []string{"/usr/bin/ls"}}}, true},
		{"TestIMA_CheckTheValidity 2:", fields{"", "", IMA{[]string{"/usr/bin/ls"}, []string{"/usr/bin/xc"}}}, false},
		{"TestIMA_CheckTheValidity 3:", fields{"", "", IMA{[]string{"/usr/bin/sleep"}, []string{"/usr/bin/cat"}}}, true},
		{"TestIMA_CheckTheValidity 4:", fields{"", "", IMA{[]string{"/usr/bin/sleep"}, []string{"/usr/bin/jl"}}}, false},
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
