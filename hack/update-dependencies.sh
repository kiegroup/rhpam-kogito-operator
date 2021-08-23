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

kogito_operator_repo="github.com/kiegroup/kogito-operator"

UPSTREAM_VERSION=$1

old_version=$(cat go.mod | grep "${kogito_operator_repo}" | head -1 | awk '{print $2}' | awk -F'-' '{print $2"-"$3}' | awk -F'.' '{print $2}')
echo "old_version = ${old_version}"

go get github.com/kiegroup/kogito-operator@${UPSTREAM_VERSION}

new_version=$(cat go.mod | grep "${kogito_operator_repo}" | head -1 | awk '{print $2}' | awk -F'-' '{print $2"-"$3}' | awk -F'.' '{print $2}')
echo "new_version = ${new_version}"

sed -i "s|${old_version}|${new_version}|g" go.mod

