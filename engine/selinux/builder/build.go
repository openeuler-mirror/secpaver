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
Package builder implements the selinux policy builder.
*/
package builder

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	pb "secpaver/api/proto/project"
	"secpaver/common/ack"
	"secpaver/common/global"
	"secpaver/common/log"
	"secpaver/common/project"
	"secpaver/common/utils"
	"secpaver/engine/selinux/pkg/secontext"
	"secpaver/engine/selinux/pkg/sehandle"
	"secpaver/engine/selinux/pkg/sepolicy"
	"secpaver/engine/selinux/pkg/serule"
	"strings"
)

const (
	fileAttribute   = "file_type"
	domainAttribute = "domain"
)

var unconfinedAttrs = []string{
	"corenet_unconfined_type",
	"dbusd_unconfined",
	"devices_unconfined_type",
	"files_unconfined_type",
	"filesystem_unconfined_type",
	"kern_unconfined",
	"selinux_unconfined_type",
	"sepgsql_unconfined_type",
	"storage_unconfined_type",
	"unconfined_domain_type",
	"unconfined_usertype",
	"xserver_unconfined_type",
}

type fClass = secontext.FileClass

// Builder is the selinux policy builder
type Builder struct {
	conf     *config
	pcHandle projectContextHandle
}

// NewBuilder is the constructor of Builder
func NewBuilder() *Builder {
	return &Builder{}
}

// Build generate selinux policy modules with project build information
func (b *Builder) Build(prjInfo *pb.ProjectInfo, out string, msg chan *pb.Ack) error {
	// 1. do prepare work
	if err := b.buildPrepare(prjInfo); err != nil {
		return err
	}

	var allPolicies []*sepolicy.Policy
	var basePolicy *sepolicy.Policy
	var err error
	var ver = prjInfo.GetMeta().GetVersion()

	// 2. create base policy
	if b.conf.getPolicy().getMonolithic() { // for single policy building
		basePolicy, err = sepolicy.CreateSePolicy(singlePolicyName(prjInfo.GetMeta().GetName()), ver)
		if err != nil {
			return errors.Wrap(err, "fail to create policy")
		}

		allPolicies = append(allPolicies, basePolicy)

	} else { // for muti policy building
		basePolicy, err = sepolicy.CreateSePolicy(basePolicyName(prjInfo.GetMeta().GetName()), ver)
		if err != nil {
			return errors.Wrap(err, "fail to create policy")
		}

		allPolicies = append(allPolicies, basePolicy)
	}

	// 2. add all extra rules to base policy
	for _, rule := range b.conf.getExtraRules() {
		r, err := serule.ParseRule(rule)
		if err != nil {
			return errors.Wrap(err, "fail to parse extra rules")
		}

		if err := basePolicy.AddRulesWithHandle(b.pcHandle.systemContextHandle, r); err != nil {
			return errors.Wrap(err, "fail to add extra rules to policy")
		}
	}

	// 3. add all file context definition to base policy
	if err := b.addFileContextDefToPolicy(basePolicy, b.pcHandle.fileItems); err != nil {
		return err
	}

	// 4. add all process context definition to base policy
	if err := b.addProcessContextDefToPolicy(basePolicy, b.pcHandle.appItems); err != nil {
		return err
	}

	// 5. add spec rules to policy
	for _, spec := range prjInfo.GetSpecs() {
		if b.conf.getPolicy().getMonolithic() { // for single policy building
			if err := b.addSpecRulesToPolicy(basePolicy, spec, msg); err != nil {
				return err
			}

		} else { // for muti policy building
			modPolicy, err := sepolicy.CreateSePolicy(modPolicyName(prjInfo.GetMeta().GetName(), spec.GetName()), ver)
			if err != nil {
				return errors.Wrap(err, "fail to create policy")
			}

			if err := b.addSpecRulesToPolicy(modPolicy, spec, msg); err != nil {
				return err
			}

			allPolicies = append(allPolicies, modPolicy)
		}
	}

	// 6. deal rules conflict and export policy to files
	for _, p := range allPolicies {
		dealAndReportTypeConflict(p, msg)

		if err := b.generatePolicyFiles(p, filepath.Join(out, p.Name)); err != nil {
			return errors.Wrapf(err, "fail to generate %s policy files", p.Name)
		}
	}

	return nil
}

func (b *Builder) buildPrepare(prjInfo *pb.ProjectInfo) error {
	// 1. check project valid
	if err := project.CheckProject(prjInfo); err != nil {
		return errors.Wrap(err, "fail to check project data")
	}

	// 2. parse policy config file
	b.conf = nil
	seConf := prjInfo.GetMeta().GetSelinux().GetConfig()
	if seConf != "" {
		if data, ok := prjInfo.GetExtends()[seConf]; ok {
			config, err := parseProjectSelinuxConfigFile(data)
			if err != nil {
				return errors.Wrap(err, "fail to parse policy config file")
			}

			b.conf = config
		} else {
			return fmt.Errorf("fail to find %s file in project struct", seConf)
		}
	}

	// 3. expand all macros and groups in project files
	if err := project.RegularProject(prjInfo); err != nil {
		return err
	}

	// 4. create system policy handle
	handle, err := sehandle.HandleCreate()
	if err != nil {
		return errors.Wrap(err, "fail to create system policy handle")
	}

	// 5. create project context handle and init all context information
	b.pcHandle = projectContextHandle{
		systemContextHandle: handle,
	}

	return b.pcHandle.initProject(prjInfo)
}

func (b *Builder) addSpecRulesToPolicy(policy *sepolicy.Policy, spec *pb.SpecInfo, msg chan *pb.Ack) error {
	var allRules []serule.Rule

	for _, info := range spec.GetApplicationList() {
		app := info.GetApplication()
		appContextInfo := b.pcHandle.getAppItemByPath(app.GetPath())
		if appContextInfo == nil {
			return fmt.Errorf("undefine application %s", app.GetPath())
		}

		fileContextInfo := appContextInfo.file
		if fileContextInfo == nil {
			return fmt.Errorf("undefine application file %s", app.GetPath())
		}

		// 1. add default rules
		allRules = append(allRules, getDefaultRules(appContextInfo.domain)...)

		// 2. add unconfined rules
		if app.GetIsUnconfined() {
			allRules = append(allRules, getUnconfinedRules(appContextInfo.domain)...)
		}

		// 3. add domain transition rules
		transRules, _ := serule.CreateDomainAutoTransRule(
			domainAttribute, fileContextInfo.context.Type, appContextInfo.domain)
		allRules = append(allRules, transRules...)

		// 4. add permission rules
		for _, permission := range info.GetPermissionList() {
			rules, err := b.getRules(appContextInfo, permission, msg)
			if err != nil {
				return errors.Wrapf(err, "fail to gen rules for application %s", app.GetPath())
			}

			allRules = append(allRules, rules...)
		}

		if err := policy.AddRulesWithHandle(b.pcHandle.systemContextHandle, allRules...); err != nil {
			return errors.Wrapf(err, "fail to add rules to policy of application %s", app.GetPath())
		}
	}

	return nil
}

func (b *Builder) addFileContextDefToPolicy(policy *sepolicy.Policy, fItems []*fileItem) error {
	for _, fItem := range fItems {
		// 1. add file context define
		if !fItem.hasFileContextDefined {
			policy.AddFileContext(
				secontext.NewFileContext(fItem.sePath, fItem.class, fItem.context))
		}

		// 2. add file type define or require
		if fItem.hasTypeDefined {
			policy.AddTypeRequire(fItem.context.Type)
			continue
		}

		// 3. if file type has not defined, add definition
		if !fItem.isPrivate {
			policy.AddTypeAttrDefine(fItem.context.Type, fileAttribute)
			continue
		}

		// 4. add private file definition and base rules
		policy.AddTypeDefine(fItem.context.Type)
		rules, err := serule.CreateBaseFileTypeRules(fItem.context.Type, fItem.class)
		if err != nil {
			return errors.Wrapf(err, "fail to create base file rules for private file %s", fItem.path)
		}

		if err := policy.AddRulesWithHandle(b.pcHandle.systemContextHandle, rules...); err != nil {
			return errors.Wrapf(err, "fail to add base file rules for private file %s", fItem.path)
		}
	}

	return nil
}

func (b *Builder) addProcessContextDefToPolicy(policy *sepolicy.Policy, aItems []*applicationItem) error {
	for _, aItem := range aItems {
		// 1. add process domain define
		policy.AddTypeAttrDefine(aItem.domain, domainAttribute)

		// 2. add unconfined attributes
		if aItem.isUnconfined {
			for _, attr := range unconfinedAttrs {
				policy.AddTypeAttrDefine(aItem.domain, attr)
			}
		}

		// 3. add permissive domain
		if aItem.isPermissive {
			policy.AddPermissiveDomain(aItem.domain)
		}

		// 4. let all roles can visit this domain
		for _, role := range b.pcHandle.systemContextHandle.GetRoles() {
			policy.AddRoleTypeDefine(role, aItem.domain)
		}
	}

	return nil
}

func (b *Builder) generatePolicyFiles(policy *sepolicy.Policy, outDir string) error {
	if err := os.RemoveAll(outDir); err != nil {
		log.Errorf("fail to remove policy output directory %s", outDir)
		return fmt.Errorf("fail to remove policy output directory")
	}

	if err := compilePolicy(policy, outDir); err != nil {
		return errors.Wrap(err, "fail to compile policy")
	}

	if err := writeResourceList(filepath.Join(outDir, "resourcelist"), b.pcHandle.fileItems); err != nil {
		return errors.Wrap(err, "fail to create resourcelist file")
	}

	var types []string
	for tp := range policy.Defines.TypeAttrDefine {
		types = append(types, tp)
	}

	if err := writeScripts(policy.Name, types, outDir); err != nil {
		return errors.Wrap(err, "fail to create policy script files")
	}

	return nil
}

func compilePolicy(policy *sepolicy.Policy, out string) error {
	teFile := filepath.Join(out, policy.Name+".te")
	fcFile := filepath.Join(out, policy.Name+".fc")

	if err := utils.WriteFile(teFile, []byte(policy.TeText()), global.DefaultFilePerm); err != nil {
		return errors.Wrapf(err, "fail to write %s file", filepath.Base(teFile))
	}

	if err := utils.WriteFile(fcFile, []byte(policy.FcText()), global.DefaultFilePerm); err != nil {
		return errors.Wrapf(err, "fail to write %s file", filepath.Base(fcFile))
	}

	modFile := strings.TrimSuffix(teFile, ".te") + ".mod"

	cmd := exec.Command("checkmodule", "-mMo", modFile, teFile)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("fail to run checkmodule command: %s", string(cmdOut))
		return fmt.Errorf("fail to run checkmodule command")
	}

	ppFile := strings.TrimSuffix(teFile, ".te") + ".pp"

	if fcFile == "" {
		cmd = exec.Command("semodule_package", "-o", ppFile, "-m", modFile)
	} else {
		cmd = exec.Command("semodule_package", "-o", ppFile, "-m", modFile, "-f", fcFile)
	}

	cmdOut, err = cmd.CombinedOutput()
	if err != nil {
		log.Errorf("fail to run semodule_package command: %s", string(cmdOut))
		return fmt.Errorf("fail to run semodule_package command")
	}

	return nil
}

func writeResourceList(path string, resources []*fileItem) error {
	var paths []string
	for _, res := range resources {
		paths = append(paths, getResourceListPaths(res.path, res.class == secontext.DirFile)...)
	}

	paths = utils.RemoveRepeatedElement(paths)

	var data bytes.Buffer
	for _, path := range paths {
		if _, err := data.WriteString(path + "\n"); err != nil {
			return err
		}
	}

	if err := utils.WriteFile(path, data.Bytes(), global.DefaultFilePerm); err != nil {
		return fmt.Errorf("fail to write file")
	}

	return nil
}

func (b *Builder) getRules(subject *applicationItem, perm *pb.Permission, msg chan *pb.Ack) ([]serule.Rule, error) {
	switch perm.Type {
	case project.RuleFileSystem:
		return b.getFileRules(subject, perm)

	case project.RuleCapability:
		return b.getCapabilityRules(subject, perm)

	case project.RuleNetWork:
		return b.getNetworkRules(subject, perm)

	default:
		sendMsg(msg, ack.LevelWarn, fmt.Sprintf("invalid permission type %s", perm.Type))
		return nil, nil
	}
}

func dealAndReportTypeConflict(p *sepolicy.Policy, msg chan *pb.Ack) {
	if m := p.DealTypeConflict(); len(m) != 0 {
		for t1, t2 := range m {
			sendMsg(msg, ack.LevelWarn, fmt.Sprintf(
				"type %s will be replaced by %s due to type_transition rule conflict", t1, t2))
		}
	}
}

func sendMsg(ch chan *pb.Ack, level, msg string) {
	if ch == nil {
		return
	}

	ch <- &pb.Ack{
		Level:  level,
		Status: msg,
	}
}

func singlePolicyName(project string) string {
	return project
}

func basePolicyName(project string) string {
	return project + "_public"
}

func modPolicyName(project, spec string) string {
	return project + "_" + utils.GetBodyFileName(spec)
}
