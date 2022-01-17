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

current_dir=$(pwd)
tmp_dir=$(mktemp -d)
kogito_operator_repo="github.com/kiegroup/kogito-operator"
local_dep=false

kogito_operator_dir=${tmp_dir}/kogito-operator
echo "Using temp dir ${kogito_operator_dir}"

echo '---- Retrieving kogito-operator repo/hash from `go.mod` file ----'
mod_ref=$(cat go.mod | grep "${kogito_operator_repo} =>")

mod_path=$(echo "${mod_ref}" | awk -F' ' '{print $3}')

echo '---- Retrieving kogito-operator repository ----'
cd ${tmp_dir}
if [[ ${mod_path} =~ ^/.* ]]; then
  echo "Copy local reference of kogito-operator"
  cp -r ${mod_path} .
else
  echo "Checkout git repository"

  kogito_operator_hash=$(echo "${mod_ref}" | awk -F'-' '{print $5}')
  echo "Got kogito-operator ${mod_path}@${kogito_operator_hash}"

  git clone https://${mod_path}.git &> /dev/null

  cd ${kogito_operator_dir}
  git reset --hard ${kogito_operator_hash}
fi

echo '---- Retrieving kogito-operator testing file(s) ----'

cd ${current_dir}

cp -r ${kogito_operator_dir}/hack/run-tests.sh hack/run-tests.sh
cp -r ${kogito_operator_dir}/hack/clean-stuck-namespaces.sh hack/clean-stuck-namespaces.sh
cp -r ${kogito_operator_dir}/hack/clean-crds.sh hack/clean-crds.sh
cp -r ${kogito_operator_dir}/hack/clean-crds.sh hack/clean-cluster-operators.sh
cp -r ${kogito_operator_dir}/test/Makefile test/Makefile
cp -r ${kogito_operator_dir}/test/features test/features
cp -r ${kogito_operator_dir}/test/examples test/examples
cp -r ${kogito_operator_dir}/test/scripts test/scripts