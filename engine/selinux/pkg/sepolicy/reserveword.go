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

package sepolicy

var reserveWords = []string{
	"alias", "allow", "and", "attribute", "attribute_role", "auditallow", "auditdeny", "bool", "category",
	"cfalse", "class", "clone", "common", "constrain", "ctrue", "dom", "domby", "dominance", "dontaudit",
	"else", "equals", "false", "filename", "filesystem", "fscon", "fs_use_task", "fs_use_trans",
	"fs_use_xattr", "genfscon", "h1", "h2", "identifier", "if", "incomp", "inherits", "iomemcon",
	"ioportcon", "ipv4_addr", "ipv6_addr", "l1", "l2", "level", "mlsconstrain", "mlsvalidatetrans", "module",
	"netifcon", "neverallow", "nodecon", "not", "notequal", "number", "object_r", "optional", "or", "path",
	"pcidevicecon", "permissive", "pirqcon", "policycap", "portcon", "r1", "r2", "r3", "range", "range_transition",
	"require", "role", "roleattribute", "roles", "role_transition", "sameuser", "sensitivity", "sid", "source", "t1",
	"t2", "t3", "target", "true", "type", "typealias", "typeattribute", "typebounds", "type_change", "type_member",
	"types", "type_transition", "u1", "u2", "u3", "user", "validatetrans", "version_identifier", "xor", "default_user",
	"default_role", "default_type", "default_range", "low", "high", "low_high",
}
