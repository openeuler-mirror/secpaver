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

TEST_NAME="TestPavProjectExportCmd"

cleanup()
{
  rm -f *.zip
  rm -f $TEMP_LOG
}

test01()
{
  test_run "$TEST_NAME/$FUNCNAME"

  # import a valid test project
  import_project "$TEST_DIR/testdata/project/valid_all.zip"
  if [ $? = 1 ]; then
    test_fail "$TEST_NAME/$FUNCNAME"
    return
  fi

  # export
  CLIENT project export valid_all "$(pwd)" > $TEMP_LOG 2>&1
  is_empty "$(grep "Finish" $TEMP_LOG)"
  if [ $? = 1 ]; then
    echo "fail to export test project"
    test_fail "$TEST_NAME/$FUNCNAME"
    delete_project "valid_all"
    return
  fi

  # delete after test
  delete_project "valid_all"
  if [ $? = 1 ]; then
    test_fail "$TEST_NAME/$FUNCNAME"
    return
  fi

  test_pass "$TEST_NAME/$FUNCNAME"
}

test02()
{
  test_run "$TEST_NAME/$FUNCNAME"

  # export an invalid project
  project_name="project_$(echo "$RANDOM")"
  CLIENT project export $project_name . > $TEMP_LOG 2>&1
  is_empty "$(grep "not found" $TEMP_LOG)"
  if [ $? = 1 ]; then
      test_fail "$TEST_NAME/$FUNCNAME"
      return
  fi

  test_pass "$TEST_NAME/$FUNCNAME"
}

zip_project "valid_all"

test01
test02

rm -f "$TEST_DIR"/testdata/project/*.zip

case_result

cleanup
