#!/bin/bash
# Copyright 2021 Red Hat, Inc. and/or its affiliates
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script act as a proxy for the `github.com/kiegroup/kogito-operator/hack/run-tests.sh` as the content remain the same

current_dir=$(pwd)

echo '---- Retrieving kogito-operator hash from `go.mod` file ----'
kogito_operator_hash=$(cat ../go.mod | grep 'github.com/kiegroup/kogito-operator' | awk -F'-' '{print $4}')
echo "Got kogito-operator hash ${kogito_operator_hash}"

tmp_dir=$(mktemp -d)

echo '---- Retrieving kogito-operator testing file ----'
cd ${tmp_dir}
git clone https://github.com/kiegroup/kogito-operator.git &> /dev/null

kogito_operator_dir=${tmp_dir}/kogito-operator
echo ${kogito_operator_dir}

cd ${kogito_operator_dir}
git reset --hard ${kogito_operator_hash}

cd ${current_dir}

cp -r ${kogito_operator_dir}/hack ../hack-kogito-operator
cp -r ${kogito_operator_dir}/test/features features
cp -r ${kogito_operator_dir}/test/examples examples
cp -r ${kogito_operator_dir}/test/scripts scripts

echo '---- Running tests ----'
../hack-kogito-operator/run-tests.sh "$@"

if [[ -d examples/logs ]]
then
  echo '---- Copying examples logs ----'
  cp -r examples/logs logs
fi

echo '---- Removing files ----'
rm -rf ../hack-kogito-operator
rm -rf features
rm -rf scripts
rm -rf examples
rm -rf ${tmp_dir}