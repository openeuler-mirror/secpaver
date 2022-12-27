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

## init all varible for test report
export TOTAL_CASE=0
export PASS_CASE=0
export FAIL_CASE=0
export SUB_CASE_FAIL=0

source ./lib/test_lib.sh

cd ./testdata/project
zip -qr valid_all.zip valid_all/
cd -

## clean
rm -rf /var/run/secpaver/*

## run pavd
start_service
sleep 2

## run test scripts
test_files=$(ls ./src)

echo "-------------------------------------------"
for file in $test_files; do
  chmod 750 ./src/$file
  . ./src/$file
  echo "-------------------------------------------"
  sleep 1
done

echo
echo "===== TOTAL: $TOTAL_CASE PASS: $PASS_CASE FAIL: $FAIL_CASE ====="

exit $FAIL_CASE

