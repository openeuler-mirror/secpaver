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

package secontext

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// RestoreconPath run 'restorecon' command to restore file's context
func RestoreconPath(path string, recurse bool) (string, error) {
	if isAbs := filepath.IsAbs(path); !isAbs {
		return "", fmt.Errorf("invalid path %s, restorecon path must be absolute", path)
	}

	var cmd *exec.Cmd
	if recurse {
		cmd = exec.Command("restorecon", "-iR", path)
	} else {
		cmd = exec.Command("restorecon", "-i", path)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("fail to run restorecon command")
	}

	return string(out), nil
}
