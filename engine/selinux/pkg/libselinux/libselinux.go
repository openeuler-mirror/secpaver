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
Package libselinux packages some c-libselinux api.
*/
package libselinux

import "C"

// #cgo pkg-config: libselinux
// #include <selinux/selinux.h>
// #include <selinux/restorecon.h>
// #include <selinux/label.h>
// #include <selinux/get_context_list.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// FileContextPath is the go wrapper function of selinux_file_context_path
func FileContextPath() string {
	var path *C.char
	path, _ = C.selinux_file_context_path()
	return C.GoString(path)
}

// FileContextHomedirPath is the go wrapper function of selinux_file_context_homedir_path
func FileContextHomedirPath() string {
	var path *C.char
	path, _ = C.selinux_file_context_homedir_path()
	return C.GoString(path)
}

// FileContextLocalPath is the go wrapper function of selinux_file_context_local_path
func FileContextLocalPath() string {
	var path *C.char
	path, _ = C.selinux_file_context_local_path()
	return C.GoString(path)
}

// CurrentPolicyPath is the go wrapper function of selinux_current_policy_path
func CurrentPolicyPath() string {
	var path *C.char
	path, _ = C.selinux_current_policy_path()
	return C.GoString(path)
}

// GetPolicyType is the go wrapper function of selinux_getpolicytype
func GetPolicyType() (string, error) {
	var tp *C.char
	defer func() {
		if tp != nil {
			C.free(unsafe.Pointer(tp))
		}
	}()

	ret, err := C.selinux_getpolicytype(&tp)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to get SELinux policy type")
		}
	}

	return C.GoString(tp), nil
}

// StringToSecurityClass is the go wrapper function of string_to_security_class
func StringToSecurityClass(name string) uint16 {
	var class C.ushort

	csName := C.CString(name)
	defer C.free(unsafe.Pointer(csName))

	class, _ = C.string_to_security_class(csName)

	return uint16(class)
}

// SecurityComputeCreate is the go wrapper function of security_compute_create
func SecurityComputeCreate(scon, tcon string, tclass uint16) (string, error) {
	var newCon *C.char

	csScon := C.CString(scon)
	defer C.free(unsafe.Pointer(csScon))

	csTcon := C.CString(tcon)
	defer C.free(unsafe.Pointer(csTcon))

	iTclass := C.ushort(tclass)

	r, err := C.security_compute_create(csScon, csTcon, iTclass, &newCon)
	if r == 0 {
		strNewcon := C.GoString(newCon)
		C.free(unsafe.Pointer(newCon))

		return strNewcon, nil
	}

	if err == nil {
		err = fmt.Errorf("fail to compute SELinux labeling decision")
	}

	return "", err
}
