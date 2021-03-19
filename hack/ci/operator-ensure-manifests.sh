#!/bin/bash
# Copyright 2020 Red Hat, Inc. and/or its affiliates
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

# Script responsible for ensuring and correcting manifests as needed.
set -e

source ./hack/ci/ensure-image.sh

tempfolder=$(mktemp -d)
echo "Temporary folder is ${tempfolder}"

version=$(getOperatorVersion)

git clone https://github.com/operator-framework/community-operators.git "${tempfolder}"
mkdir -p "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
## copy the latest manifests
cp -r bundle/manifests/ "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
cp -r bundle/metadata/ "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
cp -r bundle/tests/ "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
cp bundle.Dockerfile "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"

#Edit dockerfile with correct relative path
sed -i "s|bundle/manifests|manifests|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"
sed -i "s|bundle/metadata|metadata|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"
sed -i "s|bundle/tests|tests|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"
#replace image in target CSV
sed -i  "s|${OPERATOR_IMAGE}|${KIND_IMAGE}|g"  "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/manifests/rhpam-kogito-operator.clusterserviceversion.yaml"
#


echo "---> verify CSV updates"
cat "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/manifests/rhpam-kogito-operator.clusterserviceversion.yaml" | grep "${KIND_IMAGE}"
if grep replaces "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/manifests/rhpam-kogito-operator.clusterserviceversion.yaml"
then
  latest_released_version=$(getLatestOlmReleaseVersion)
  ##ensure correct replace field is there
  sed -i "s|replace.*|replaces: rhpam-kogito-operator.v${latest_released_version}|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/manifests/rhpam-kogito-operator.clusterserviceversion.yaml"
  #
  cat "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/manifests/rhpam-kogito-operator.clusterserviceversion.yaml" | grep "replaces: rhpam-kogito-operator.v${latest_released_version}"
fi
