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

rules=$(grep "avc" $audit_log_file | grep -E "$Types" | audit2allow | grep allow)
IFS=$'\n'

if [ "$1" = "json" ]
then
	echo '"extraRules":['
	line=""
	for rule in $rules
	do
		if [ "$line" != "" ]
		then
			echo '    "'$line'",'
		fi

		line=$rule
	done

	if [ "$line" != "" ]
		then
			echo '    "'$line'"'
	fi

	echo ']'
else
	for rule in $rules
	do
		echo "$rule"
	done
fi
