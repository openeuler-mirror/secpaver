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

TEST_NAME="TestPavPolicyExportCmd"

cleanup()
{
  rm -f $TEMP_LOG
  rm -f *.zip
}

test01()
{
  test_run "$TEST_NAME/$FUNCNAME"

  CLIENT project build --engine selinux -d "$TEST_DIR/testdata/project/valid_all" > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
  if [ $? = 1 ]; then
    echo "fail to build valid_all project to SELinux policy"
    test_fail "$TEST_NAME/$FUNCNAME"
    return
  fi

  CLIENT policy export valid_all_public_selinux . > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
  if [ $? = 1 ]; then
    echo "fail to export valid_all_public_selinux policy"
    test_fail "$TEST_NAME/$FUNCNAME"
    return
  fi

  test_pass "$TEST_NAME/$FUNCNAME"
}

test02()
{
  test_run "$TEST_NAME/$FUNCNAME"

  CLIENT policy export "policy_$(echo "$RANDOM")_selinux" . > $TEMP_LOG 2>&1
  is_empty "$(grep "not found" $TEMP_LOG)"
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
