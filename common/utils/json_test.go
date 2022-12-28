/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2021. All rights reserved.
 * secPaver is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package utils

import (
	"reflect"
	"testing"
)

type testStruct struct {
	FieldInt    int               `json:"fieldInt"`
	FieldString string            `json:"fieldString"`
	FieldBool   bool              `json:"fieldBool"`
	FieldStruct *testFieldStruct  `json:"fieldStruct"`
	Extend      map[string][]byte `json:"-"`
}

type testFieldStruct struct {
	FieldInt    int    `json:"fieldInt"`
	FieldString string `json:"fieldString"`
	FieldBool   bool   `json:"fieldBool"`
}

var testJSON = []byte(`{
	"fieldInt": 1,
	"fieldString": "test",
	"fieldBool": true,
	"fieldStruct": {
		"fieldInt": 1,
		"fieldString": "test",
		"fieldBool": true
	},
	"fieldExtend1": "extend1",
	"fieldExtend2": 1,
	"fieldExtend3": true
}`)

var testJSONRaw = []byte(`{"fieldBool":false,"fieldExtend1":"extend1","fieldExtend2":1,"fieldExtend3":true,"fieldInt":1,"fieldString":"test","fieldStruct":{"fieldBool":true,"fieldInt":1,"fieldString":"test"}}`)

var validObject = &testStruct{}
var i = 0
var invalidObject1 = &i
var invalidObject2 *testStruct
var completeObject = &testStruct{
	FieldInt:    1,
	FieldString: "test",
	FieldBool:   false,
	FieldStruct: &testFieldStruct{
		FieldInt:    1,
		FieldString: "test",
		FieldBool:   true,
	},
	Extend: map[string][]byte{
		"fieldExtend1": []byte(`"extend1"`),
		"fieldExtend2": []byte(`1`),
		"fieldExtend3": []byte(`true`),
	},
}

// test function
func TestUnmarshalJSONWithExtendParams(t *testing.T) {
	type args struct {
		data []byte
		s    interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		{
			args: args{
				data: testJSON,
				s:    validObject,
			},
			want: map[string][]byte{
				"fieldExtend1": []byte(`"extend1"`),
				"fieldExtend2": []byte(`1`),
				"fieldExtend3": []byte(`true`),
			},
		},
		{
			args: args{
				data: testJSON,
				s:    invalidObject1,
			},
			wantErr: true,
		},
		{
			args: args{
				data: testJSON,
				s:    invalidObject2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalJSONWithExtendParams(tt.args.data, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSONWithExtendParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSONWithExtendParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestMarshalJSONWithExtendParams(t *testing.T) {
	type args struct {
		s       interface{}
		extData map[string][]byte
		extKey  string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				s:       completeObject,
				extData: completeObject.Extend,
				extKey:  "Extend",
			},
			want: testJSONRaw,
		},
		{
			args: args{
				s:       invalidObject1,
				extData: completeObject.Extend,
				extKey:  "Extend",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalJSONWithExtendParams(tt.args.s, tt.args.extData, tt.args.extKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSONWithExtendParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSONWithExtendParams() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}

