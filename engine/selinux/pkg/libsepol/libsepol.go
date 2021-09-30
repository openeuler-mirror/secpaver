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
Package libsepol packages some c-libsepol api.
*/
package libsepol

/*
#cgo pkg-config: libsepol
#include <sepol/sepol.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	"unsafe"
)

// Fopen is the go wrapper function of fopen
func Fopen(path string) (*C.struct__IO_FILE, error) {
	var fp *C.struct__IO_FILE

	csPath := C.CString(path)
	defer C.free(unsafe.Pointer(csPath))

	csArg := C.CString("rb")
	defer C.free(unsafe.Pointer(csArg))

	fp, err := C.fopen(csPath, csArg)
	if err != nil {
		return nil, fmt.Errorf("fail to open file %s", filepath.Base(path))
	}

	return fp, nil
}

// Fclose is the go wrapper function of fclose
func Fclose(fp *C.struct__IO_FILE) {
	if fp != nil {
		C.fclose(fp)
	}
}

// HandleCreate is the go wrapper function of sepol_handle_create
func HandleCreate() (*C.struct_sepol_handle, error) {
	var db *C.struct_sepol_handle
	db, err := C.sepol_handle_create()
	if err != nil {
		return nil, errors.Wrap(err, "fail to create sepol handle")
	}

	return db, nil
}

// HandleDestroy is the go wrapper function of sepol_handle_destroy
func HandleDestroy(handle *C.struct_sepol_handle) {
	C.sepol_handle_destroy(handle)
}

// PolicyFileCreate is the go wrapper function of sepol_policy_file_create
func PolicyFileCreate() (*C.struct_sepol_policy_file, error) {
	var pf *C.struct_sepol_policy_file
	r, err := C.sepol_policy_file_create(&pf)
	if r < 0 {
		if err != nil {
			return nil, fmt.Errorf("fail to create sepol policy file")
		}

		return nil, errors.Wrap(err, "fail to create sepol policy file")
	}

	return pf, nil
}

// PolicyFileFree is the go wrapper function of sepol_policy_file_free
func PolicyFileFree(spf *C.struct_sepol_policy_file) {
	C.sepol_policy_file_free(spf)
}

// PolicyFileSetFp is the go wrapper function of sepol_policy_file_set_fp
func PolicyFileSetFp(spf *C.struct_sepol_policy_file, fp *C.struct__IO_FILE) {
	C.sepol_policy_file_set_fp(spf, fp)
}

// PolicyFileSetHandle is the go wrapper function of sepol_policy_file_set_handle
func PolicyFileSetHandle(spf *C.struct_sepol_policy_file, handle *C.struct_sepol_handle) {
	C.sepol_policy_file_set_handle(spf, handle)
}

// PolicydbCreate is the go wrapper function of sepol_policydb_create
func PolicydbCreate() (*C.struct_sepol_policydb, error) {
	var db *C.struct_sepol_policydb
	r, err := C.sepol_policydb_create(&db)
	if r < 0 {
		if err != nil {
			return nil, fmt.Errorf("fail to create sepol policydb")
		}

		return nil, errors.Wrap(err, "fail to create sepol policydb")
	}

	return db, nil
}

// PolicydbFree is the go wrapper function of sepol_policydb_free
func PolicydbFree(db *C.struct_sepol_policydb) {
	C.sepol_policydb_free(db)
}

// PolicydbRead is the go wrapper function of sepol_policydb_read
func PolicydbRead(db *C.struct_sepol_policydb, spf *C.struct_sepol_policy_file) error {
	r, err := C.sepol_policydb_read(db, spf)
	if r < 0 {
		if err != nil {
			return fmt.Errorf("fail to read sepol policydb")
		}

		return errors.Wrap(err, "fail to read sepol policydb")
	}

	return nil
}
