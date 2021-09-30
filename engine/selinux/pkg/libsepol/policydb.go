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

package libsepol

/*
#cgo pkg-config: libsepol
#include <sepol/sepol.h>
#include <sepol/policydb/policydb.h>
#include <stdlib.h>
#include "port.h"
#include "context.h"
*/
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"unsafe"
)

const (
	symRoles       = 2
	symTypes       = 3
	symUsers       = 4
	typeAttribute  = 1
	oconPort       = 2
	lenTempArray   = 1 << 20
	maxPortRecords = 4096
)

// raw flag of proto
const (
	IPProtoTCP  = 6
	IPProtoUDP  = 17
	IPProtoDCCP = 33
	IPProtoSCTP = 132
)

// GetAllRoles returns the role list stored in policydb
func GetAllRoles(db *C.struct_sepol_policydb) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("policy db is nil")
	}

	if len(db.p.symtab) <= symRoles || len(db.p.sym_val_to_name) <= symRoles {
		return nil, fmt.Errorf("wrong length of sym field in policydb")
	}

	allRoles, err := goStrings(
		db.p.symtab[symRoles].nprim, db.p.sym_val_to_name[symRoles])

	if err != nil {
		return nil, errors.Wrap(err, "fail to get role list from sepolicy database")
	}

	return allRoles, nil
}

// GetAllUsers returns the user list stored in policydb
func GetAllUsers(db *C.struct_sepol_policydb) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("policy db is nil")
	}

	if len(db.p.symtab) <= symUsers || len(db.p.sym_val_to_name) <= symUsers {
		return nil, fmt.Errorf("wrong length of sym field in policydb")
	}

	allUsers, err := goStrings(
		db.p.symtab[symUsers].nprim, db.p.sym_val_to_name[symUsers])

	if err != nil {
		return nil, errors.Wrap(err, "fail to get user list from sepolicy database")
	}

	return allUsers, nil
}

// TypeInfo is the struct of type information stored in policydb
type TypeInfo struct {
	Name   string
	IsAttr bool
}

// GetAllTypesAndAttrs returns the type information list stored in policydb
func GetAllTypesAndAttrs(db *C.struct_sepol_policydb) ([]*TypeInfo, error) {
	if db == nil {
		return nil, fmt.Errorf("policy db is nil")
	}

	if len(db.p.symtab) <= symTypes || len(db.p.sym_val_to_name) <= symTypes {
		return nil, fmt.Errorf("wrong length of sym field in policydb")
	}

	all, err := goStrings(
		db.p.symtab[symTypes].nprim, db.p.sym_val_to_name[symTypes])

	if err != nil {
		return nil, errors.Wrap(err, "fail to get type list from sepolicy database")
	}

	l := len(all)
	if l > lenTempArray {
		return nil, fmt.Errorf("the number of SELinux types exceeds the upper limit")
	}

	tempSlice := (*[lenTempArray]*C.struct_type_datum)(unsafe.Pointer(db.p.type_val_to_struct))[:l:l]

	var allInfos []*TypeInfo

	// note: the type and attrs in policydb use same struct and differentiated by flavor field
	for i, s := range tempSlice {
		allInfos = append(allInfos, &TypeInfo{
			Name:   all[i],
			IsAttr: int(s.flavor) == typeAttribute,
		})
	}

	return allInfos, nil
}

// ContextID is the raw context express struct
type ContextID struct {
	User uint32
	Role uint32
	Type uint32
}

// PortRecord is the raw record for port context
type PortRecord struct {
	Protocol uint32
	LowPort  uint32
	HighPort uint32
	Context  ContextID
}

// GetAllPorts returns the port context information stored in policydb
func GetAllPorts(db *C.struct_sepol_policydb) ([]*PortRecord, error) {
	if db == nil {
		return nil, fmt.Errorf("policy db is nil")
	}

	var list []*PortRecord
	var cnt = 0

	if len(db.p.ocontexts) <= oconPort {
		return nil, fmt.Errorf("wrong length of ocontexts field in policydb")
	}

	for oc := db.p.ocontexts[oconPort]; oc != nil; oc = oc.next {
		if cnt > maxPortRecords {
			return nil, fmt.Errorf(
				"the number of port records in policydb exceeds the upper limit")
		}
		cnt++

		// note: in policydb struct, the port context record is stored by union.
		// the field in union should be extract in c program
		var tmpProto C.uint8_t
		var tmpLowPort, tmpHighPort C.uint16_t
		C.PortTrans(oc, &tmpProto, &tmpLowPort, &tmpHighPort)

		if len(oc.context) <= 1 {
			return nil, fmt.Errorf("wrong length of context field in occontext")
		}

		info := &PortRecord{
			Protocol: uint32(tmpProto),
			LowPort:  uint32(tmpLowPort),
			HighPort: uint32(tmpHighPort),
			Context: ContextID{
				User: uint32(oc.context[0].user),
				Role: uint32(oc.context[0].role),
			},
		}

		var tmpType C.uint32_t
		// note: 'type' is a keyword in golang, should be deal in c program
		C.ContextTypeTrans(&oc.context[0], &tmpType)
		info.Context.Type = uint32(tmpType)

		list = append(list, info)
	}

	return list, nil
}

func goStrings(argc C.uint, argv **C.char) ([]string, error) {
	if argv == nil {
		return nil, fmt.Errorf("the address of string array is nil")
	}

	if argc == 0 {
		return []string{}, nil
	}

	length := uint32(argc)
	if length > lenTempArray {
		return nil, fmt.Errorf("the number of CString exceeds the upper limit")
	}

	tempSlice := (*[lenTempArray]*C.char)(unsafe.Pointer(argv))[:length:length]

	result := make([]string, length, length)
	for i, s := range tempSlice {
		result[i] = C.GoString(s)
	}

	return result, nil
}
