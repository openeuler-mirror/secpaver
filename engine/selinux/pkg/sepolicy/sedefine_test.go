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

package sepolicy

import "testing"

// test function
func TestSeDefine_Text(t *testing.T) {
	type fields struct {
		RoleTypeDefine map[string][]string
		TypeAttrDefine map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{RoleTypeDefine: map[string][]string{
				"test_r": {"type_1", "type_2"},
			}},
			want: "role test_r types type_1;\nrole test_r types type_2;\n",
		},
		{
			fields: fields{TypeAttrDefine: map[string][]string{
				"test_t": {"attr_1", "attr_2"},
			}},
			want: "type test_t, attr_1, attr_2;\n",
		},
		{
			fields: fields{TypeAttrDefine: map[string][]string{
				"test_t": {"", "attr_1"},
			}},
			want: "type test_t, attr_1;\n",
		},
		{
			fields: fields{TypeAttrDefine: map[string][]string{
				"test_t": {""},
			}},
			want: "type test_t;\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := &SeDefine{
				RoleTypeDefine: tt.fields.RoleTypeDefine,
				TypeAttrDefine: tt.fields.TypeAttrDefine,
			}
			if got := def.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

