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
Package service implements the secPaver server service.
*/
package service

import (
	"fmt"
	"secpaver/common/server"
	"secpaver/engine"
	repo "secpaver/repository"
)

// NewServerServices creates all server service inst
func NewServerServices(repo repo.Repo, engines map[string]engine.Engine) ([]server.Service, error) {
	engSvr, err := NewEngineService(engines)
	if err != nil {
		return nil, fmt.Errorf("fail to create engine service")
	}

	prjSvr, err := NewProjectService(repo, engines)
	if err != nil {
		return nil, fmt.Errorf("fail to create project service")
	}

	polSvr, err := NewPolicyService(repo, engines)
	if err != nil {
		return nil, fmt.Errorf("fail to create policy service")
	}

	return []server.Service{
		engSvr,
		prjSvr,
		polSvr,
	}, nil
}
