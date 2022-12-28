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

import (
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/sehandle"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
	"os"
	"reflect"
	"testing"
	"time"
)

var testHandle sehandle.Handle

// test function
func TestCreateSePolicy(t *testing.T) {
	type args struct {
		name string
		ver  string
	}
	tests := []struct {
		want    *Policy
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{name: "demo", ver: "1.0"},
			want: &Policy{
				Requires: NewSeRequire(),
				Defines:  NewSeDefine(),
				Name:     "demo",
				Ver:      "1.0",
			},
			wantErr: false,
		},
		{
			args:    args{name: "alias", ver: "1.0"},
			wantErr: true,
		},
		{
			args:    args{name: "demo", ver: "alias"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSePolicy(tt.args.name, tt.args.ver)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSePolicy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getBasePolicy() *Policy {
	p, _ := CreateSePolicy("test", "1.0")
	return p
}

func policyHasTypeRequire(p *Policy, tp string) bool {
	return utils.IsExistItem(tp, p.Requires.TypeRequires)
}

func policyHasAttrRequire(p *Policy, attr string) bool {
	return utils.IsExistItem(attr, p.Requires.AttrRequires)
}

func policyHasTypeDefined(p *Policy, tp string) bool {
	_, ok := p.Defines.TypeAttrDefine[tp]
	return ok
}

func policyHasTypeAttrDefined(p *Policy, tp string, attrs []string) bool {
	as, ok := p.Defines.TypeAttrDefine[tp]
	if !ok {
		return false
	}

	for _, attr := range attrs {
		if attr == "" {
			continue
		}

		if !utils.IsExistItem(attr, as) {
			return false
		}
	}

	return true
}

func policyHasClassPermsRequired(p *Policy, class string, perms []string) bool {
	ps, ok := p.Requires.ClassRequires[class]
	if !ok {
		return false
	}

	for _, perm := range perms {
		if !utils.IsExistItem(perm, ps) {
			return false
		}
	}

	return true
}

// test function
func TestPolicyOperateType(t *testing.T) {
	// test add type require
	p := getBasePolicy()
	p.AddTypeRequire("test_t")
	if !policyHasTypeRequire(p, "test_t") {
		t.Errorf("fail to test add type require")
	}

	// test add type define
	p = getBasePolicy()
	p.AddTypeDefine("test_t")
	if !policyHasTypeDefined(p, "test_t") {
		t.Errorf("fail to test add type define")
	}

	// test add type attr define
	p = getBasePolicy()
	p.AddTypeAttrDefine("test_t", "attr")
	if (!policyHasTypeAttrDefined(p, "test_t", []string{"attr"})) ||
		(!policyHasAttrRequire(p, "attr")) {

		t.Errorf("fail to test add type define")
	}

	// test add type require when type has been defined
	p = getBasePolicy()
	p.AddTypeDefine("test_t")
	p.AddTypeRequire("test_t")
	if (!policyHasTypeDefined(p, "test_t")) || (policyHasTypeRequire(p, "test_t")) {
		t.Errorf("fail to test add type require when type has been defined")
	}

	// test add type define when type has been required
	p = getBasePolicy()
	p.AddTypeRequire("test_t")
	p.AddTypeDefine("test_t")
	if (!policyHasTypeDefined(p, "test_t")) || (policyHasTypeRequire(p, "test_t")) {
		t.Errorf("fail to test add type require when type has been defined")
	}

	// test add type attr define after adding type define
	p = getBasePolicy()
	p.AddTypeDefine("test_t")
	p.AddTypeAttrDefine("test_t", "attr")
	if (!policyHasTypeAttrDefined(p, "test_t", []string{"attr"})) ||
		(!policyHasAttrRequire(p, "attr")) {

		t.Errorf("test add type attr define after adding type define")
	}
}

// test function
func TestPolicyOperateRule(t *testing.T) {
	testAvcRuleType, _ := serule.ParseAvcRule("allow test_t file_t:file { read };")
	testAvcRuleAttr, _ := serule.ParseAvcRule("allow test_t file_type:file { read };")
	// test add avc rule type
	p := getBasePolicy()
	if err := p.AddRulesWithHandle(testHandle, testAvcRuleType); err != nil {
		t.Errorf("fail to add test avc rule")
	}
	if (!policyHasClassPermsRequired(p, "file", []string{"read"})) ||
		(!policyHasTypeRequire(p, "test_t")) ||
		(!policyHasTypeRequire(p, "file_t")) {

		t.Errorf("fail to test add avc rule type")
	}

	// test add avc rule attribute
	p = getBasePolicy()
	if err := p.AddRulesWithHandle(testHandle, testAvcRuleAttr); err != nil {
		t.Errorf("fail to add test avc rule")
	}
	if (!policyHasClassPermsRequired(p, "file", []string{"read"})) ||
		(!policyHasTypeRequire(p, "test_t")) ||
		(!policyHasAttrRequire(p, "file_type")) {

		t.Errorf("fail to test add avc rule type")
	}
}

var testPolicyText = `module test 1.0;

require
{
	role test_r;
	type bin_t;
	type object_t;
	type target1_t;
	type target2_t;
	type target3_t;
	type target4_t;
	type unconfined_t;
	attribute test_attr;
	class file { read write };
};

role test_r types test_t;
type test_t, test_attr;

permissive test_t;

allow unconfined_t bin_t : file { read write };

type_transition unconfined_t object_t : file target1_t;
type_transition unconfined_t object_t : file target1_t;
type_transition unconfined_t object_t : file target3_t "name";
type_transition unconfined_t object_t : file target3_t "name";

`

var testAvcRules = []*serule.AvcRule{
	{
		Prefix:  "allow",
		Subject: "unconfined_t",
		Object:  "bin_t",
		Class:   "file",
		Actions: []string{"write", "read"},
	},
	{
		Prefix:  "allow",
		Subject: "unconfined_t",
		Object:  "bin_t",
		Class:   "file",
		Actions: []string{"read"},
	},
}

var testTypeRule = []*serule.TypeRule{
	{
		Prefix:     "type_transition",
		Subject:    "unconfined_t",
		Object:     "object_t",
		Class:      "file",
		Target:     "target1_t",
		ObjectName: "",
	},
	{
		Prefix:     "type_transition",
		Subject:    "unconfined_t",
		Object:     "object_t",
		Class:      "file",
		Target:     "target2_t",
		ObjectName: "",
	},
	{
		Prefix:     "type_transition",
		Subject:    "unconfined_t",
		Object:     "object_t",
		Class:      "file",
		Target:     "target3_t",
		ObjectName: "name",
	},
	{
		Prefix:     "type_transition",
		Subject:    "unconfined_t",
		Object:     "object_t",
		Class:      "file",
		Target:     "target4_t",
		ObjectName: "name",
	},
}

// test function
func TestPolicyGen(t *testing.T) {
	p := getBasePolicy()

	p.AddTypeRequire("test_t")
	p.AddAttrRequire("test_attr")
	p.AddRoleRequire("test_r")
	p.AddRoleTypeDefine("test_r", "test_t")
	p.AddTypeAttrDefine("test_t", "test_attr")
	p.AddPermissiveDomain("test_t")

	for _, rule := range testAvcRules {
		_ = p.AddRulesWithHandle(testHandle, rule)
	}

	for _, rule := range testTypeRule {
		_ = p.AddRulesWithHandle(testHandle, rule)
	}

	p.DealTypeConflict()
	if p.TeText() != testPolicyText {
		t.Errorf("fail to test policy gen, want: %s, got %s", testPolicyText, p.TeText())
	}
}

// main test function
func TestMain(m *testing.M) {
	// handle can't be opened too frequently, maybe other testcase just open the handle
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 5)
		testHandle, _ = sehandle.HandleCreate()
		if testHandle != nil {
			break
		}
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

