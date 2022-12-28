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

package pid

import (
	"gotest.tools/v3/fs"
	"testing"
)

// test function
func Test_PidLock_Lock(t *testing.T) {
	dir := fs.NewDir(t, "test",
		fs.WithFile("unlock.pid", ""),
		fs.WithFile("lock.pid", ""))
	defer dir.Remove()

	unlockPl, _ := CreatePidLock(dir.Join("unlock.pid"))
	lockPl, _ := CreatePidLock(dir.Join("lock.pid"))
	_ = lockPl.Lock()
	testPl, _ := CreatePidLock(dir.Join("lock.pid"))

	tests := []struct {
		fields  *Lock
		name    string
		wantErr bool
	}{
		{
			fields:  testPl,
			wantErr: true,
		},
		{
			fields:  unlockPl,
			wantErr: false,
		},
		{
			fields:  nil,
			wantErr: true,
		},
		{
			fields:  &Lock{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.Lock(); (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func Test_PidLock_Unlock(t *testing.T) {
	dir := fs.NewDir(t, "test",
		fs.WithFile("unlock.pid", ""),
		fs.WithFile("lock.pid", ""))
	defer dir.Remove()

	unlockPl, _ := CreatePidLock(dir.Join("unlock.pid"))
	lockPl, _ := CreatePidLock(dir.Join("lock.pid"))
	_ = lockPl.Lock()

	tests := []struct {
		fields  *Lock
		name    string
		wantErr bool
	}{
		{
			fields:  unlockPl,
			wantErr: false,
		},
		{
			fields:  lockPl,
			wantErr: false,
		},
		{
			fields:  nil,
			wantErr: true,
		},
		{
			fields:  &Lock{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.fields
			if err := l.Unlock(); (err != nil) != tt.wantErr {
				t.Errorf("Unlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func Test_PidLock_Release(t *testing.T) {
	dir := fs.NewDir(t, "test",
		fs.WithFile("unlock.pid", ""),
		fs.WithFile("lock.pid", ""))
	defer dir.Remove()

	unlockPl, _ := CreatePidLock(dir.Join("unlock.pid"))
	lockPl, _ := CreatePidLock(dir.Join("lock.pid"))

	tests := []struct {
		fields  *Lock
		name    string
		wantErr bool
	}{
		{
			fields:  unlockPl,
			wantErr: false,
		},
		{
			fields:  lockPl,
			wantErr: false,
		},
		{
			fields:  nil,
			wantErr: true,
		},
		{
			fields:  &Lock{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.fields
			l.Release()
		})
	}
}

