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

if [[ -d test/examples/logs ]]
then
  echo '---- Copying examples logs ----'
  cp -r test/examples/logs test/logs
fi

echo '---- Removing files ----'
rm -rf hack/run-tests.sh
rm -rf test/Makefile
rm -rf test/features
rm -rf test/scripts
rm -rf test/examples

exit ${testsStatus}