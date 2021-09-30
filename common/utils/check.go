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
Package utils provides some util functions
*/
package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// const variable definition
const (
	CheckExactArgs = iota
	CheckMinArgs

	maxArgLength       = 50
	maxPathLength      = 150
	maxFileSize        = 3 * 1024 * 1024
	maxFileSizeInBzip2 = 10 * 1024 * 1024
	maxFileNumberInZip = 50
	maxFileSizeInZip   = 1 * 1024 * 1024
	maxVersionLength   = 8
)

var (
	unsafeArgRegexp   = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_.-]*$")
	versionRegexp     = regexp.MustCompile("^[\\d+][.\\d+]*$")
	unsafePathRegexp  = regexp.MustCompile(`^[./a-zA-Z0-9_-]+$`)
	selinuxTypeRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
)

// CheckCommandArgs checks command args number
func CheckCommandArgs(context *cli.Context, expected, checkType int) error {
	if context == nil {
		return fmt.Errorf("nil cil context")
	}

	var err error
	cmdName := context.Command.Name
	switch checkType {
	case CheckExactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf("%s: %q requires exactly %d argument(s)", os.Args[0], cmdName, expected)
		}
	case CheckMinArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q requires a minimum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	default:
		return fmt.Errorf("invalid check type %d", checkType)
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		_ = cli.ShowCommandHelp(context, cmdName)
		return err
	}

	return nil
}

// CheckUnsafeArg checks if a parameter is safe
// 1. the length of parameter should be less than MaxArgLength
// 2. the parameter characters can only contain a-z A-Z 0-9 . _ -
// 3. the first char of parameter should be a-z A-Z _
func CheckUnsafeArg(arg string) error {
	if len(arg) > maxArgLength {
		return fmt.Errorf("the length of parameter should be less than %d", maxArgLength)
	}

	if !unsafeArgRegexp.MatchString(arg) {
		return fmt.Errorf("%s is an invalid arg", arg)
	}

	return nil
}

// CheckZipFileName checks if a zip file name is legal
func CheckZipFileName(file string) error {
	ext := filepath.Ext(file)
	if strings.ToLower(ext) != ".zip" {
		return fmt.Errorf("%s should end with .zip", file)
	}
	prefix := strings.TrimSuffix(file, ext)

	if err := CheckUnsafeArg(prefix); err != nil {
		return errors.Wrap(err, "invalid zip file name")
	}

	return nil
}

// CheckPort checks if a port is legal
func CheckPort(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > math.MaxUint16 {
		return 0, fmt.Errorf("%s is an invalid port", port)
	}

	return p, nil
}

// CheckUnsafePath checks if a path is safe
// 1. the length of path should less than MaxPathLength
// 2. the path can't contains '..'
// 3. the characters of path can only contain . / a-z A-Z 0-9 _ -
func CheckUnsafePath(path string) error {
	if len(path) > maxPathLength {
		return fmt.Errorf("invalid path %s, file path should be shorter than %d", path, maxPathLength)
	}

	if strings.Contains(path, "..") {
		return fmt.Errorf("unsafe path %s, path can't contains '..'", path)
	}

	if !unsafePathRegexp.MatchString(path) {
		return fmt.Errorf("unsafe path %s, path can only contain '. / a-z A-Z 0-9 _ -'", path)
	}

	return nil
}

// CheckValidSelinuxType check if an type is a valid SELinux type
func CheckValidSelinuxType(str string) error {
	if len(str) > maxArgLength {
		return fmt.Errorf("the length of SELinux type should be less than %d", maxArgLength)
	}

	if !selinuxTypeRegexp.MatchString(str) {
		errStr := fmt.Sprintf(
			"error: invalid selinux type %s, the type can only contain character, number and '_'", str)

		return fmt.Errorf(errStr)
	}

	return nil
}

// CheckVersion check if an type is a valid version string
func CheckVersion(str string) error {
	if len(str) > maxVersionLength {
		return fmt.Errorf("the length of version should be less than %d", maxVersionLength)
	}

	if !versionRegexp.MatchString(str) {
		return fmt.Errorf("error: invalid version %s", str)
	}

	return nil
}
