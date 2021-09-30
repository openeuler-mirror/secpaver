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
Package libsemanage packages some c-libsemanage api.
*/
package libsemanage

// #cgo pkg-config: libsemanage
// #include <semanage/semanage.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// HandleCreate is the go wrapper function of semanage_handle_create
func HandleCreate() (*C.struct_semanage_handle, error) {
	sh, err := C.semanage_handle_create()
	return sh, err
}

// HandleDestroy is the go wrapper function of semanage_handle_destroy
func HandleDestroy(sh *C.struct_semanage_handle) {
	C.semanage_handle_destroy(sh)
}

// Connect is the go wrapper function of semanage_connect
func Connect(sh *C.struct_semanage_handle) error {
	ret, err := C.semanage_connect(sh)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to connect semanage")
		}
		return err
	}

	return nil
}

// Disconnect is the go wrapper function of semanage_disconnect
func Disconnect(sh *C.struct_semanage_handle) {
	C.semanage_disconnect(sh)
}

// BeginTransaction is the go wrapper function of semanage_begin_transaction
func BeginTransaction(sh *C.struct_semanage_handle) error {
	ret, err := C.semanage_begin_transaction(sh)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to begin semanage transaction")
		}
		return err
	}

	return nil
}

// Commit is the go wrapper function of semanage_commit
func Commit(sh *C.struct_semanage_handle) error {
	ret, err := C.semanage_commit(sh)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to commit semanage transaction")
		}
		return err
	}

	return nil
}

// ModuleKeyCreate is the go wrapper function of semanage_module_key_create
func ModuleKeyCreate(sh *C.struct_semanage_handle) (*C.struct_semanage_module_key, error) {
	var key *C.struct_semanage_module_key
	ret, err := C.semanage_module_key_create(sh, &key)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to create semanage module key")
		}
		return nil, err
	}

	return key, nil
}

// ModuleKeyDestroy is the go wrapper function of semanage_module_key_destroy
func ModuleKeyDestroy(sh *C.struct_semanage_handle, key *C.struct_semanage_module_key) {
	C.semanage_module_key_destroy(sh, key)
	if key != nil {
		C.free(unsafe.Pointer(key))
	}
}

// ModuleKeySetName is the go wrapper function of semanage_module_key_set_name
func ModuleKeySetName(sh *C.struct_semanage_handle, key *C.struct_semanage_module_key, name string) error {
	csName := C.CString(name)
	defer C.free(unsafe.Pointer(csName))

	ret, err := C.semanage_module_key_set_name(sh, key, csName)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to set semanage module key name")
		}
		return err
	}

	return nil
}

// ModuleGetModuleInfo is the go wrapper function of semanage_module_get_module_info
func ModuleGetModuleInfo(
	sh *C.struct_semanage_handle, key *C.struct_semanage_module_key) (*C.struct_semanage_module_info, error) {

	var info *C.struct_semanage_module_info
	ret, err := C.semanage_module_get_module_info(sh, key, &info)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to get semanage module info")
		}
		return nil, err
	}

	return info, nil
}

// ModuleInfoGetPriority is the go wrapper function of semanage_module_info_get_priority
func ModuleInfoGetPriority(sh *C.struct_semanage_handle, info *C.struct_semanage_module_info) (int, error) {
	var priority C.ushort
	ret, err := C.semanage_module_info_get_priority(sh, info, &priority)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to get semanage module info")
		}
		return 0, err
	}

	return int(priority), nil
}

// ModuleInfoGetEnabled is the go wrapper function of semanage_module_info_get_enabled
func ModuleInfoGetEnabled(sh *C.struct_semanage_handle, info *C.struct_semanage_module_info) (bool, error) {
	var enabled C.int
	ret, err := C.semanage_module_info_get_enabled(sh, info, &enabled)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to get semanage module enabled")
		}
		return false, err
	}

	return int(enabled) > 0, nil
}

// ModuleInfoDestroy is the go wrapper function of semanage_module_info_destroy
func ModuleInfoDestroy(sh *C.struct_semanage_handle, info *C.struct_semanage_module_info) {
	C.semanage_module_info_destroy(sh, info)
	C.free(unsafe.Pointer(info))
}

// ModuleRemove is the go wrapper function of semanage_module_remove
func ModuleRemove(sh *C.struct_semanage_handle, name string) error {
	csName := C.CString(name)
	defer C.free(unsafe.Pointer(csName))

	ret, err := C.semanage_module_remove(sh, csName)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to remove semanage module")
		}
		return err
	}

	return nil
}

// ModuleInstallFile is the go wrapper function of semanage_module_install_file
func ModuleInstallFile(sh *C.struct_semanage_handle, file string) error {
	csFile := C.CString(file)
	defer C.free(unsafe.Pointer(csFile))

	ret, err := C.semanage_module_install_file(sh, csFile)
	if ret < 0 {
		if err == nil {
			err = fmt.Errorf("fail to install semanage module")
		}
		return err
	}

	return nil
}

// Root is the go wrapper function of semanage_root
func Root() string {
	var path *C.char
	path, _ = C.semanage_root()

	return C.GoString(path)
}

// SetRebuild is the go wrapper function of semanage_set_rebuild
func SetRebuild(sh *C.struct_semanage_handle) {
	C.semanage_set_rebuild(sh, C.int(1))
}
