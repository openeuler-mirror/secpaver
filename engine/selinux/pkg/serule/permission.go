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

package serule

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"secpaver/common/utils"
)

// GetAllPermissionsOfClass returns all valid permission of a SELinux class
func GetAllPermissionsOfClass(class string) ([]string, error) {
	classDir := filepath.Join(classRoot, class)
	exist, err := utils.PathExist(classDir)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to search %s class directory", class)
	}

	if !exist {
		return nil, fmt.Errorf("invalid selinux class %s", class)
	}

	infos, err := ioutil.ReadDir(filepath.Join(classDir, "perms"))
	if err != nil {
		return nil, fmt.Errorf("fail to read selinux perms of %s class directory", class)
	}

	var perms []string
	for _, info := range infos {
		if !info.IsDir() {
			perms = append(perms, info.Name())
		}
	}

	if len(perms) == 0 {
		return nil, fmt.Errorf("no valid permission for %s class", class)
	}

	return perms, nil
}
