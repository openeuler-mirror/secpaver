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
Package pid implements the pid file lock function
*/
package pid

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/global"
	"strconv"
	"syscall"
)

// Lock is the pid file lock to ensure that the process runs in a single instance
type Lock struct {
	lock *os.File
	file string
}

// CreatePidLock creates a pid lock file with specified filepath
func CreatePidLock(path string) (*Lock, error) {
	if path == "" {
		return nil, fmt.Errorf("empty pid file path")
	}

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, global.PidFilePerm)
	if err != nil {
		return nil, fmt.Errorf("fail to open pid file %s", filepath.Base(path))
	}

	return &Lock{
		file: path,
		lock: fd,
	}, nil
}

// Lock locks the pid lock
func (l *Lock) Lock() error {
	if l == nil || l.lock == nil {
		return fmt.Errorf("nil pid lock")
	}

	if err := syscall.Flock(int(l.lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return errors.Wrap(err, "fail to flock pid file")
	}

	pid := strconv.Itoa(os.Getpid())
	if err := l.lock.Truncate(0); err != nil {
		return errors.Wrap(err, "fail to truncate file")
	}

	if _, err := l.lock.Seek(0, 0); err != nil {
		return errors.Wrap(err, "fail to seek file")
	}

	if _, err := l.lock.WriteString(pid); err != nil {
		return errors.Wrap(err, "fail to write pid file")
	}

	return nil
}

// Unlock unlocks the pid lock
func (l *Lock) Unlock() error {
	if l == nil || l.lock == nil {
		return fmt.Errorf("nil pid lock")
	}

	if err := syscall.Flock(int(l.lock.Fd()), syscall.LOCK_UN); err != nil {
		return errors.Wrap(err, "fail to unflock pid file")
	}

	return nil
}

// Release releases the pid lock
func (l *Lock) Release() {
	if l == nil || l.lock == nil {
		return
	}

	_ = l.lock.Close()
	_ = os.Remove(l.file)
}
