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

TEMP_LOG="temp.log"

SERVICE_NAME=pavd

TEST_DIR=$(pwd)

PROJECT_DIR="/var/local/secpaver/projects"
ENGINE_DIR="/var/local/secpaver/engines"

SUPPORT_ENGINE=("selinux")

TEST_FAIL="--- FAIL    "
TEST_PASS="--- PASS    "
TEST_RUN="=== RUN     "

CLIENT()
{
  pav $*
}

is_empty()
{
  if [ "$1" = "" ]; then
    return 1
  fi
    return 0
}

test_run()
{
  echo "$TEST_RUN$1"
}

test_pass()
{
  echo "$TEST_PASS$1"
}

test_fail()
{
  SUB_CASE_FAIL=$((SUB_CASE_FAIL+1))
  echo "$TEST_FAIL$1"
  echo "log:"
  cat $TEMP_LOG
}

case_result()
{
  TOTAL_CASE=$((TOTAL_CASE+1))
  if [ "$SUB_CASE_FAIL" -eq "0" ];then
    case_pass
  else
    case_fail
  fi
}

case_pass()
{
  PASS_CASE=$((PASS_CASE+1))
  SUB_CASE_FAIL=0
}

case_fail()
{
  FAIL_CASE=$((FAIL_CASE+1))
  SUB_CASE_FAIL=0
}

start_service()
{
  systemctl restart pavd
}

import_project()
{
  CLIENT project import $1 -f > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
    if [ $? = 1 ]; then
      echo "fail to import test project"
      return 1
    fi
  return 0
}

delete_project()
{
  CLIENT project delete $1 > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
    if [ $? = 1 ]; then
      echo "fail to delete test project"
      return 1
    fi
  return 0
}

delete_policy()
{
  CLIENT policy delete $1 > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
    if [ $? = 1 ]; then
      echo "fail to delete test policy"
      return 1
    fi
  return 0
}

uninstall_policy()
{
  CLIENT policy uninstall $1 > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
    if [ $? = 1 ]; then
      echo "fail to uninstall test policy"
      return 1
    fi
  return 0
}

zip_project()
{
    cd "$TEST_DIR"/testdata/project
    project_name=$1
    if [ ! -d "$project_name" ]; then
      echo "no dir $project_name for zip"
      return 0
    fi
    rm -f "$project_name.zip"
    zip -qr "$project_name.zip" "$project_name"
    cd - >/dev/null 2>&1
    return 1
}
