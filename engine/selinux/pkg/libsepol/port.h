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

#ifndef SEPORT_H
#define SEPORT_H

#include <sepol/policydb/policydb.h>

void PortTrans(const ocontext_t *oc, uint8_t *protocol, uint16_t *lowPort, uint16_t *highPort)
{
    if ((oc == NULL) || (protocol == NULL) || (lowPort == NULL) || (highPort == NULL)) {
        return;
    }

    *protocol = oc->u.port.protocol;
    *lowPort = oc->u.port.low_port;
    *highPort = oc->u.port.high_port;
}

#endif