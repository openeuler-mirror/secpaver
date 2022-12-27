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

TEST_NAME="TestPavEngineListCmd"

cleanup()
{
  rm -f policy*.zip
  rm -f $TEMP_LOG
}

test01()
{
  test_run "$TEST_NAME/$FUNCNAME"

  # check all supported engines were listed
  CLIENT engine list > $TEMP_LOG 2>&1
  for ((i=0;i<${#SUPPORT_ENGINE[@]};i++));do
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

  # check help info
  CLIENT engine list -h > $TEMP_LOG 2>&1
  is_empty "$(grep "List usable engines in secPaver server" $TEMP_LOG)"
  if [ $? = 1 ]; then
    test_fail "$TEST_NAME/$FUNCNAME"
    return
  fi

  test_pass "$TEST_NAME/$FUNCNAME"
}

test03()
{
  test_run "$TEST_NAME/$FUNCNAME"

  # check extra input
  CLIENT engine list extra_input > $TEMP_LOG 2>&1
  is_empty "$(grep "Incorrect Usage." $TEMP_LOG)"
  if [ $? = 1 ]; then
    test_fail "$TEST_NAME/$FUNCNAME"
    return
  fi

  test_pass "$TEST_NAME/$FUNCNAME"
}

test01
test02
test03

case_result

cleanup

