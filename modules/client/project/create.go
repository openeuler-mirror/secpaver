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

package project

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/errdefs"
	"gitee.com/openeuler/secpaver/common/global"
	"gitee.com/openeuler/secpaver/common/utils"
	"text/template"
)

func newCreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Usage:       "Create a template project",
		UsageText:   "Pav project create <NAME> <PATH>",
		Description: "Create a template project with the specified project name",
		Action:      createAction,
	}
}

func createCheck(ctx *cli.Context) error {
	// check arg number
	if err := utils.CheckCommandArgs(ctx, 2, utils.CheckExactArgs); err != nil {
		return err
	}

	// check valid project name
	if err := utils.CheckUnsafeArg(ctx.Args()[0]); err != nil {
		return err
	}

	// check the directory path exist
	dirExists, err := utils.DirExist(ctx.Args()[1])
	if err != nil {
		return err
	}

	if !dirExists {
		return errdefs.NewDirNotFoundError(ctx.Args()[1])
	}

	return nil
}

func createAction(ctx *cli.Context) error {
	if err := createCheck(ctx); err != nil {
		return err
	}

	prjName := ctx.Args()[0]
	prjPath := filepath.Join(ctx.Args()[1], prjName)
	exist, err := utils.PathExist(prjPath)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("%s already exists", prjPath)
	}

	if err := os.MkdirAll(filepath.Join(prjPath, "specs"), global.DefaultDirPerm); err != nil {
		return err
	}

	dataMap := map[string]string{
		"Name": prjName,
	}

	if err := writeTemplateFile(prjPath, "pav.proj", metaTemplate, dataMap); err != nil {
		return err
	}

	if err := writeTemplateFile(prjPath, "resources.json", resourcesTemplate, dataMap); err != nil {
		return err
	}

	if err := writeTemplateFile(prjPath, "selinux.json", selinuxConfigTemplate, dataMap); err != nil {
		return err
	}

	specPath := filepath.Join(prjPath, "specs")
	specFilename := fmt.Sprintf("module_%s.json", prjName)

	if err := writeTemplateFile(specPath, specFilename, specTemplate, dataMap); err != nil {
		return err
	}

	fmt.Printf("Finish creating %s template project at %s\n", prjName, prjPath)

	return nil
}

func writeTemplateFile(path, filename, content string, data map[string]string) error {
	f, err := os.OpenFile(filepath.Join(path, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, global.DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("fail to create %s file", filename)
	}

	defer f.Close()

	tmpl, err := template.New(filename).Parse(content)
	if err != nil {
		return errors.Wrap(err, "fail to parse template")
	}

	return tmpl.Execute(f, data)
}

var metaTemplate = `{
    "version":"1.0",
    "name":"{{.Name}}",
    "resources":"resources.json",
    "specs":[
        "specs/module_{{.Name}}.json"
    ],
    "selinux":{
        "config":"selinux.json"
    }
}
`

var resourcesTemplate = `{
    "macroList":[
        {
            "name":"BIN",
            "value":"/usr/bin"
        }
    ],
    "groupList":[
        {
            "name":"RESOURCE",
            "resources":[
                "/resource"
            ]
        }
    ],
    "resourceList":[
        {
            "type":"exec_file",
            "path":"$(BIN)/example",
            "selinux":{
                "isSysFile":false,
                "isPrivateFile":false,
                "domain":"example_t",
                "type":"example_exec_t"
            }
        },
        {
            "type":"file",
            "path":"/resource",
            "selinux":{
                "isSysFile":false,
                "isPrivateFile":false,
                "domain":"",
                "type":""
            }
        }
    ]
}`

var specTemplate = `{
    "applicationList":[
        {
            "application":{
                "path":"$(BIN)/example",
                "isPermissive":false
            },
            "permissionList":[
                {
                    "type":"filesystem",
                    "resources":[
                        "$(RESOURCE)"
                    ],
                    "actions":[
                        "create",
                        "read",
                        "write"
                    ]
                },
                {
                    "type":"network",
                    "resources":[
                        "domain:unix,type:stream,protocol:"
                    ],
                    "actions":[
                        "accept",
                        "bind",
                        "connect",
                        "listen",
                        "receive",
                        "send",
                        "listen"
                    ]
                },
                {
                    "type":"capability",
                    "resources":[

                    ],
                    "actions":[
                        "chown",
                        "dac_override",
                        "dac_read_search"
                    ]
                }
            ]
        }
    ]
}`

var selinuxConfigTemplate = `{
    "policy":{
        "monolithic":true
    },
    "extraRules":[
    ]
}`
