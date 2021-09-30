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
Package engine implements all engines function
*/
package engine

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"plugin"
	"secpaver/api/proto/policy"
	"secpaver/api/proto/project"
	"secpaver/common/log"
	"secpaver/domain"
	"syscall"
)

// Engine is the interface that an engine plugin must implement
type Engine interface {
	GetName() string
	GetDescription() string
	Build(prjInfo *project.ProjectInfo, out string, msg chan *project.Ack) error
	GetPolicyStatus(policy *domain.Policy) (string, error)
	Install(policy *domain.Policy, msg chan *policy.Ack) error
	Uninstall(policy *domain.Policy, msg chan *policy.Ack) error
}

// LoadEnginePlugin loads all engine plugins in specified directory
func LoadEnginePlugin(dir string) (map[string]Engine, error) {
	m := make(map[string]Engine, 2)

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("fail to read plugin directory")
	}

	for _, info := range infos {
		if (info.IsDir()) || (filepath.Ext(info.Name()) != ".so") {
			continue
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if (!ok) || (stat.Uid != 0) {
			_ = fmt.Errorf("fail to check uid of %s", info.Name())
			continue
		}

		eng, err := loadEnginePlugin(filepath.Join(dir, info.Name()))
		if err != nil {
			return nil, err
		}

		if err := registerEngine(eng, m); err != nil {
			return nil, errors.Wrap(err, "fail to register engine")
		}

		log.Infof("load %s engine plugin successfully", eng.GetName())
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("no valid engine")
	}

	return m, nil
}

func loadEnginePlugin(path string) (Engine, error) {
	enginePlugin, err := plugin.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to open engine plugin %s", filepath.Base(path))
	}

	f, err := enginePlugin.Lookup("GetEngine")
	if err != nil {
		return nil, err
	}

	if ft, ok := f.(func() Engine); ok {
		return ft(), nil
	}

	return nil, fmt.Errorf("invalid engine: %s", path)
}

func registerEngine(engine Engine, m map[string]Engine) error {
	if engine == nil {
		return fmt.Errorf("nil engine")
	}

	if engine.GetName() == "" {
		return fmt.Errorf("invalid engine name")
	}

	if _, ok := m[engine.GetName()]; ok {
		return fmt.Errorf("engine %s existed", engine.GetName())
	}

	m[engine.GetName()] = engine

	return nil
}
