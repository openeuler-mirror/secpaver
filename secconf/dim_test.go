package secconf

import (
	"testing"
)

func TestCheckTheValidity(t *testing.T) {
	type fields struct {
		TemplatePath      string
		TemplatePathCheck string
		Dim               Dim
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{"TestDim_CheckTheValidity 1:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 2:", fields{"", "", Dim{false, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 3:", fields{"", "", Dim{true, false, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 4:", fields{"", "", Dim{true, true, CoreFunc{100, "sm2", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 5:", fields{"", "", Dim{true, true, CoreFunc{100, "md5", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 6:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 129, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 7:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 525601, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 8:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 520000, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 9:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 1001, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 10:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 1000, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 11:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, false}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 12:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 10, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 13:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 101, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, true},
		{"TestDim_CheckTheValidity 14:", fields{"", "", Dim{true, true, CoreFunc{99, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 15:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "md5", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 16:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 250}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 17:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, make([]string, 1500, 2000), []string{"1", "2", "3"}, false}}, false},
		{"TestDim_CheckTheValidity 18:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, make([]string, 1501, 2000), false}}, false},
		{"TestDim_CheckTheValidity 19:", fields{"", "", Dim{true, true, CoreFunc{100, "sha256", 11, 0, 0, true}, MonitorFunc{"", 100, "sha256", 12}, []string{"1", "2", "3"}, []string{"1", "2", "3"}, true}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Dim_yaml{
				//				templatePath:      tt.fields.TemplatePath,
				//				templatePathCheck: tt.fields.TemplatePathCheck,
				Dim: tt.fields.Dim,
			}
			if got := config.checkTheValidity(); got != tt.want {
				t.Errorf("checkTheValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}
