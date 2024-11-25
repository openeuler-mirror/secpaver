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
		{name: "TestDim_CheckTheValidity 1:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 2:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: false, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 3:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: false,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 4:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sm3", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 5:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "md5", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 6:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 129, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 7:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 525601, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 8:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 525600, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 9:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 1001, MeasureInterval: 100, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 10:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 1000, MeasureInterval: 100, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 11:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: false}}, want: true},
		{name: "TestDim_CheckTheValidity 12:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 99, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 13:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 250,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 14:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: make([]string, 1500, 2000), BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: false},
		{name: "TestDim_CheckTheValidity 15:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: false, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 10000, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 16:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: false, MeasureList: make([]string, 1500, 2000), BaselineIsEnable: true,
				MeasureLogCapacity: 101, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: true},
		{name: "TestDim_CheckTheValidity 17:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: false, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 99, MeasureHash: "sha256", CorePcr: 11, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: false}}, want: true},
		{name: "TestDim_CheckTheValidity 18:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: false, MeasureList: []string{"1", "2", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 99, MeasureHash: "sha256", CorePcr: 255, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: false}}, want: true},
		{name: "TestDim_CheckTheValidity 19:", fields: fields{TemplatePath: "", TemplatePathCheck: "",
			Dim: Dim{DimIsEnable: true, MeasureList: []string{"1", "", "3"}, BaselineIsEnable: true,
				MeasureLogCapacity: 99, MeasureHash: "sha256", CorePcr: 255, MonitorPcr: 12,
				MeasureSchedule: 100, MeasureInterval: 100, Signature: true}}, want: false},
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
