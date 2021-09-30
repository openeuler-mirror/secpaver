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

echo
echo "start to recompile module:"
checkmodule -Mmo "$mod_file" "$te_file" || { echo "command failed"; exit 1; }

if [ -d "$fc_file" ]; then
  semodule_package -o "$pp_file" -m "$mod_file" -f "$fc_file" || { echo "command failed"; exit 1; }
else
  semodule_package -o "$pp_file" -m ""$mod_file"" || { echo "command failed"; exit 1; }
fi

echo "finish"
