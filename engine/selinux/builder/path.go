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
	"path/filepath"
	"strings"
)

var endWildcardChangeMap = map[string]string{
	"/*":     "/[^/]*",
	"/**":    "/.*",
	"{,/*}":  "(/[^/]*)?",
	"{,/**}": "(/.*)?",
}

// NOTE: changeWildcardToRegExp is not a strict change tool, currently
// can deal some simple cases
func changeWildcardToRegExp(path string) string {
	var suffix string
	for k, v := range endWildcardChangeMap {
		if strings.HasSuffix(path, k) {
			suffix = v
			path = strings.TrimSuffix(path, k)
		}
	}

	return dealSpecialChar(path) + suffix
}

var specialChar = []string{"[", "]", "(", ")", "^", "+", "."}

func dealSpecialChar(path string) string {
	validEscape := `\\`
	str := path

	for _, c := range specialChar {
		str = strings.Replace(str, c, "\\"+c, -1)
	}

	for idx := 0; idx < len(str); idx++ {
		b := str[idx]
		if b == '?' {
			if findNumOfEscape(str, idx-1)%len(validEscape) == 0 {
				str = str[:idx] + "[^/]" + str[idx+1:]
				idx += len("[^/]")
				continue
			}

		} else if b == '*' {
			if findNumOfEscape(str, idx-1)%len(validEscape) != 0 {
				continue
			}

			if idx < len(str)-1 && str[idx+1] == '*' {
				str = str[:idx] + ".*" + str[idx+1+1:]
				idx += len(".*")
			} else {
				str = str[:idx] + "[^/]*" + str[idx+1:]
				idx += len("[^/]*")
			}
		}
	}

	return str
}

func findNumOfEscape(str string, idx int) int {
	num := 0
	for i := idx; i >= 0; i-- {
		if str[i] == '\\' {
			num++
		} else {
			break
		}
	}
	return num
}

func getDir(path string) string {
	if strings.HasSuffix(path, "{,/*}") ||
		strings.HasSuffix(path, "{,/**}") {

		return filepath.Dir(filepath.Dir(path))
	}

	return filepath.Dir(path)
}

func getBase(path string) string {
	for str := range endWildcardChangeMap {
		if strings.HasSuffix(path, str) {
			return ""
		}
	}

	return filepath.Base(path)
}

var linkMap = map[string]string{
	"/run":                        "/var/run",
	"/run/lock":                   "/var/lock",
	"/run/systemd/system":         "/usr/lib/systemd/system",
	"/run/systemd/generator":      "/usr/lib/systemd/system",
	"/run/systemd/generator.late": "/usr/lib/systemd/system",
	"/lib":                        "/usr/lib",
	"/lib64":                      "/usr/lib",
	"/usr/lib64":                  "/usr/lib",
	"/usr/local/lib64":            "/usr/lib",
	"/usr/local/lib32":            "/usr/lib",
	"/etc/systemd/system":         "/usr/lib/systemd/system",
	"/var/lib/xguest/home":        "/home",
	"/var/named/chroot/usr/lib64": "/usr/lib",
	"/var/named/chroot/lib64":     "/usr/lib",
	"/var/home":                   "/home",
	"/home-inst ":                 "/home",
	"/home/home-inst":             "/home",
	"/var/roothome":               "/root",
	"/sbin":                       "/usr/sbin",
	"/sysroot/tmp":                "/tmp",
	"/bin":                        "/usr/bin",
}

func dealLinkPath(path string) string {
	for k, v := range linkMap {
		if strings.HasPrefix(path, k) {
			return v + strings.TrimPrefix(path, k)
		}
	}

	return path
}

func getSePath(path string) string {
	return changeWildcardToRegExp(dealLinkPath(path))
}

func getSepPaths(path string) []string {
	if strings.Contains(path, "{,/*}") {
		return []string{
			strings.Replace(path, "{,/*}", "", -1),
			strings.Replace(path, "{,/*}", "/*", -1),
		}
	}

	if strings.Contains(path, "{,/**}") {
		return []string{
			strings.Replace(path, "{,/**}", "", -1),
			strings.Replace(path, "{,/**}", "/**", -1),
		}
	}

	return []string{path}
}

func getResourceListPaths(path string, isDir bool) []string {
	path = strings.Replace(dealLinkPath(path), "/**", "/*", -1)
	sepPaths := getSepPaths(path)
	if len(sepPaths) != 1 {
		return sepPaths
	}

	if isDir {
		return []string{path, filepath.Join(path, "*")}
	}

	return []string{path}
}
