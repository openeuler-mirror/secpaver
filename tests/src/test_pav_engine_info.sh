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

source ./lib/test_lib.sh

TEST_NAME="TestPavEngineInfoCmd"

cleanup()
{
  rm -f policy*.zip
  rm -f $TEMP_LOG
}

test01()
{
  test_run "$TEST_NAME/$FUNCNAME"

  # check all supported engines
  for ((i=0;i<${#SUPPORT_ENGINE[@]};i++));do
    CLIENT engine info "${SUPPORT_ENGINE[i]}" > $TEMP_LOG 2>&1
    is_empty "$(grep "${SUPPORT_ENGINE[i]}" $TEMP_LOG)"
    if [ $? = 1 ]; then
      test_fail "$TEST_NAME/$FUNCNAME"
      return
    fi
  done

  test_pass "$TEST_NAME/$FUNCNAME"
}

test02()
{
  test_run "$TEST_NAME/$FUNCNAME"

  # check invalid engine
  CLIENT engine info invalid_engine > $TEMP_LOG 2>&1
  is_empty "$(grep "invalid engine" $TEMP_LOG)"

  if [ $? = 1 ]; then
      test_fail "$TEST_NAME/$FUNCNAME"
      return
  fi

  test_pass "$TEST_NAME/$FUNCNAME"
}

test01
test02

case_result

cleanup

