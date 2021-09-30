#!/bin/sh
# Copyright (c) Huawei Technologies Co., Ltd. 2020-2021. All rights reserved.
# secPaver is licensed under the Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#     http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v2 for more details.

this_dir=$(pwd)
conf="$this_dir/config"

if [ -f "$conf" ];then
	echo "load config file $conf"
	. $conf
else
	echo "config file $conf not found"
	exit 1
fi

pp_file="$this_dir/$module_name.pp"
te_file="$this_dir/$module_name.te"
fc_file="$this_dir/$module_name.fc"
mod_file="$this_dir/$module_name.mod"
resourcelist_file="$this_dir/$resourcelist_name"

echo
echo "start to install SELinux policy module:"

## insert modules
semodule -i $pp_file || { echo "command failed"; exit 1; }

echo "finish"
echo
echo "start to refresh file contexts"

## restore file contexts of the resources defined in project
## the resources are listed in the resourcelist file
cat "$resourcelist_file" | while read line
do
	echo "refresh file context of $line"
	restorecon -vR $line
done

echo "finish"
