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

package builder

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/global"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/utils"
	"text/template"
)

var (
	configFile  = "config"
	scriptFiles = []string{
		"install.sh",
		"uninstall.sh",
		"update.sh",
		"audit.sh",
		"restorecon.sh",
	}
)

func writeScripts(module string, types []string, outDir string) error {
	if err := writeConfigFile(module, types, outDir); err != nil {
		return err
	}

	for _, name := range scriptFiles {
		src := filepath.Join(global.DefaultScriptRoot, "selinux", name)
		dst := filepath.Join(outDir, name)

		if err := utils.CopyFile(src, dst); err != nil {
			log.Errorf("fail to write %s script file: %v", name, err)
			return fmt.Errorf("fail to write %s script file", name)
		}
	}

	return nil
}

func writeConfigFile(module string, types []string, outDir string) error {
	configTmplFile := filepath.Join(global.DefaultScriptRoot, "selinux", configFile)

	tmpl, err := template.New(configFile).ParseFiles(configTmplFile)
	if err != nil {
		return errors.Wrap(err, "fail to compile config template file")
	}

	var typeList string
	for i, tp := range types {
		if i != 0 {
			typeList += "|"
		}

		typeList += tp
	}

	data := map[string]string{
		"Module":   module,
		"TypeList": typeList,
	}

	path := filepath.Join(outDir, configFile)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, global.DefaultDirPerm)
	if err != nil {
		log.Errorf("fail to create config file %s", filepath.Base(path))
		return fmt.Errorf("fail to create config file %s", filepath.Base(path))
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
