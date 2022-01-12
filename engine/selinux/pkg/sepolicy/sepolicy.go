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

/*
Package sepolicy is a tool to parse and generate selinux policy.
*/
package sepolicy

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/sehandle"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
)

// Policy is the selinux policy profile model
type Policy struct {
	Requires          *SeRequire
	Defines           *SeDefine
	AvcRules          []*serule.AvcRule
	TypeRules         []*serule.TypeRule
	PermissiveDomains []string
	FileContexts      []*secontext.FileContext
	Name              string
	Ver               string
}

// CreateSePolicy returns a policy profile model with the specified name and version
func CreateSePolicy(name, ver string) (*Policy, error) {
	if name == "" || ver == "" {
		return nil, fmt.Errorf("blank policy name or version")
	}

	if err := utils.CheckVersion(ver); err != nil {
		return nil, fmt.Errorf("fail to check policy version")
	}

	if err := utils.CheckUnsafeArg(name); err != nil {
		return nil, fmt.Errorf("fail to check policy name")
	}

	if utils.IsExistItem(name, reserveWords) {
		return nil, fmt.Errorf("policy name %s is a SELinux reserver word", name)
	}

	return &Policy{
		Name:     name,
		Ver:      ver,
		Requires: NewSeRequire(),
		Defines:  NewSeDefine(),
	}, nil
}

// AddFileContext add file context definition
func (p *Policy) AddFileContext(fc ...*secontext.FileContext) {
	p.FileContexts = append(p.FileContexts, fc...)
}

// AddTypeRequire add a type require statement to SePolicy struct
func (p *Policy) AddTypeRequire(tp string) {
	// check if the type has been defined in current policy
	if _, ok := p.Defines.TypeAttrDefine[tp]; !ok {
		p.Requires.AddTypeRequire(tp)
	}
}

// AddRoleRequire add a role require statement to SePolicy struct
func (p *Policy) AddRoleRequire(role string) {
	p.Requires.AddRoleRequire(role)
}

// AddAttrRequire add a attribute require statement to SePolicy struct
func (p *Policy) AddAttrRequire(attr string) {
	p.Requires.AddAttrRequire(attr)
}

// AddClassRequire add a class require statement to SePolicy struct
func (p *Policy) AddClassRequire(cls string, actions []string) error {
	return p.Requires.AddClassRequire(cls, actions)
}

// AddRoleTypeDefine adds a role statement with type definition to SePolicy struct
func (p *Policy) AddRoleTypeDefine(role, tp string) {
	// add role require first
	p.Requires.AddRoleRequire(role)
	p.Defines.AddRoleTypeDefine(role, tp)
}

// AddTypeAttrDefine adds a type statement with attribute definition to SePolicy struct
func (p *Policy) AddTypeAttrDefine(tp, attr string) {
	// add attribute require first
	p.Requires.AddAttrRequire(attr)

	p.Defines.AddTypeAttrDefine(tp, attr)
	// if the type has been added in require list, remove it
	p.Requires.RemoveTypeRequire(tp)
}

// AddTypeDefine adds a type definition to SePolicy struct
func (p *Policy) AddTypeDefine(tp string) {
	p.Defines.AddTypeDefine(tp)
	// if the type has been added in require list, remove it
	p.Requires.RemoveTypeRequire(tp)
}

// AddPermissiveDomain adds a permissive domain
func (p *Policy) AddPermissiveDomain(domain string) {
	// add type require first
	p.AddTypeRequire(domain)

	if !utils.IsExistItem(domain, p.PermissiveDomains) {
		p.PermissiveDomains = append(p.PermissiveDomains, domain)
	}
}

// AddObjectRequireWithHandle add a object(type or attribute) require information in SePolicy
func (p *Policy) AddObjectRequireWithHandle(handle sehandle.Handle, obj string) {
	if handle == nil {
		return
	}

	if handle.AttrHasDefined(obj) {
		p.AddAttrRequire(obj)
	} else {
		p.AddTypeRequire(obj)
	}
}

// AddRulesWithHandle adds a rule to SePolicy struct with the require information
func (p *Policy) AddRulesWithHandle(handle sehandle.Handle, rules ...serule.Rule) error {
	for _, rule := range rules {
		if rule == nil {
			continue
		}

		switch rule.(type) {
		case *serule.AvcRule:
			if err := p.AddAvcRulesWithHandle(handle, rule.(*serule.AvcRule)); err != nil {
				return errors.Wrapf(err, "fail to add avc rule %s", rule.Text())
			}

		case *serule.TypeRule:
			if err := p.AddTypeRulesWithHandle(handle, rule.(*serule.TypeRule)); err != nil {
				return errors.Wrapf(err, "fail to add avc rule %s", rule.Text())
			}

		default:
			return fmt.Errorf("invalid rule type")
		}
	}

	return nil
}

// AddAvcRulesWithHandle adds a avc rule to SePolicy struct with the require information
func (p *Policy) AddAvcRulesWithHandle(handle sehandle.Handle, rules ...*serule.AvcRule) error {
loop:
	for _, rule := range rules {
		if rule == nil {
			continue
		}

		if err := p.AddClassRequire(rule.Class, rule.Actions); err != nil {
			return errors.Wrapf(err, "fail to add class require for %s", rule.Text())
		}

		p.AddObjectRequireWithHandle(handle, rule.Subject)
		p.AddObjectRequireWithHandle(handle, rule.Object)

		for _, r := range p.AvcRules {
			if r.Merge(rule) {
				continue loop
			}
		}

		p.AvcRules = append(p.AvcRules, rule)
	}

	return nil
}

// AddTypeRulesWithHandle adds a type rule to SePolicy struct with the require information
func (p *Policy) AddTypeRulesWithHandle(handle sehandle.Handle, rules ...*serule.TypeRule) error {
loop:
	for _, rule := range rules {
		if rule == nil {
			continue
		}

		p.AddObjectRequireWithHandle(handle, rule.Subject)
		p.AddObjectRequireWithHandle(handle, rule.Object)
		p.AddObjectRequireWithHandle(handle, rule.Target)

		for _, r := range p.TypeRules {
			if r.Text() == rule.Text() {
				continue loop
			}
		}

		p.TypeRules = append(p.TypeRules, rule)
	}

	return nil
}

// TeText generate a string of SELinux TE file
func (p *Policy) TeText() string {
	var bt bytes.Buffer

	// write module information
	bt.WriteString(genModuleInfo(p.Name, p.Ver))

	// write requires
	bt.WriteString(p.Requires.Text())
	bt.WriteString("\n")

	// write defines
	bt.WriteString(p.Defines.Text())
	bt.WriteString("\n")

	// write domain permissive statement
	for _, domain := range p.PermissiveDomains {
		bt.WriteString(genDomainPermissive(domain))
	}
	bt.WriteString("\n")

	// write rules
	for _, rule := range p.AvcRules {
		bt.WriteString(rule.Text())
	}
	bt.WriteString("\n")

	for _, rule := range p.TypeRules {
		bt.WriteString(rule.Text())
	}
	bt.WriteString("\n")

	return bt.String()
}

// FcText generate a string of SELinux FC file
func (p *Policy) FcText() string {
	var bt bytes.Buffer

	for _, fc := range p.FileContexts {
		bt.WriteString(fc.Text())
	}

	return bt.String()
}

// DealTypeConflict deals the type_transition in policy and returns the modification of type
func (p *Policy) DealTypeConflict() map[string]string {
	m := make(map[string]string, len(p.TypeRules))
	for i := range p.TypeRules {
		for j := i + 1; j < len(p.TypeRules); j++ {
			if p.TypeRules[i].CheckConflict(p.TypeRules[j]) == nil {
				continue
			}

			m[p.TypeRules[j].Target] = p.TypeRules[i].Target
			p.replaceType(p.TypeRules[j].Target, p.TypeRules[i].Target)
		}
	}

	return m
}

func (p *Policy) replaceType(old, new string) {
	for _, rule := range p.TypeRules {
		if rule.Subject == old {
			rule.Subject = new
		}

		if rule.Object == old {
			rule.Object = new
		}

		if rule.Target == old {
			rule.Target = new
		}
	}

	for _, rule := range p.AvcRules {
		if rule.Subject == old {
			rule.Subject = new
		}

		if rule.Object == old {
			rule.Object = new
		}
	}

	for _, ctx := range p.FileContexts {
		if ctx.Context.Type == old {
			ctx.Context.Type = new
		}
	}
}

func genModuleInfo(name, ver string) string {
	return fmt.Sprintf("module %s %s;\n\n", name, ver)
}

func genDomainPermissive(domain string) string {
	return fmt.Sprintf("permissive %s;\n", domain)
}
