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
Package sepath provides some of for SELinux style path.
*/
package sepath

import "gitee.com/openeuler/secpaver/common/utils"

var metaSet = []rune{'.', '^', '$', '?', '*', '+', '|', '[', ']', '(', ')', '{', '}'}

// GetFixedPrefix return the prefix of path that not have meta char
func GetFixedPrefix(path string) string {
	for i := 0; i < len(path); i++ {
		c := rune(path[i])

		if c == '\\' {
			i++
			continue
		}

		if utils.IsExistItem(c, metaSet) {
			return path[:i]
		}
	}

	return path
}
