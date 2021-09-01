#!/usr/bin/env bats

source ./hack/env.sh

OLD_VERSION=998.0.0
NEW_VERSION=999.0.0
NEW_VERSION_MAJOR_MINOR=$(echo "${NEW_VERSION}" | awk -F. '{print $1"."$2}')
CURRENT_VERSION=$(getOperatorVersion)

function make() { 
    echo "make $@" 
}

setup() {
    # copy structure in a temp folder for the test
    dir="${BATS_TMPDIR}/${BATS_TEST_NAME}"
    rm -rf ${dir}
    mkdir -p ${dir}
    cp -r ${BATS_TEST_DIRNAME}/.. ${dir}/

    mkdir -p ${dir}/bundle/
}

teardown() {
    rm -rf "$BATS_TMPDIR/$BATS_TEST_NAME"
}

@test "check bump-version error with no version" {
    run ${BATS_TEST_DIRNAME}/bump-version.sh
    [ "$status" -eq 1 ]
    [[ "${output}" =~ "Please inform the new version. Use X.X.X" ]]
}

@test "check csv file is set correctly, no option" {
    # export -f make

    dir="${BATS_TMPDIR}/${BATS_TEST_NAME}"
    cd ${dir}
    run hack/bump-version.sh ${NEW_VERSION}
    [ "$status" -eq 0 ]

    # Check csv file
    [[ "${output}" =~ "Version bumped from ${CURRENT_VERSION} to ${NEW_VERSION}" ]]
    [[ "${output}" =~ "Set version ${NEW_VERSION} with bundle suffix 1 and replacing version (if not empty) " ]]
    # fine tune results
    csv_file=$(cat $(getCsvFile))
    [[ "${csv_file}" =~ "containerImage: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${csv_file}" =~ "name: rhpam-kogito-operator.v${NEW_VERSION}-1" ]]
    [[ "${csv_file}" =~ "operated-by: rhpam-kogito-operator.${NEW_VERSION}-1" ]]
    [[ "${csv_file}" =~ "version: ${NEW_VERSION}-1" ]]
    bundle_csv_file=$(cat $(getBundleCsvFile))
    [[ "${bundle_csv_file}" =~ "containerImage: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${bundle_csv_file}" =~ "name: rhpam-kogito-operator.v${NEW_VERSION}-1" ]]
    [[ "${bundle_csv_file}" =~ "image: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${bundle_csv_file}" =~ "operated-by: rhpam-kogito-operator.${NEW_VERSION}-1" ]]
    [[ "${bundle_csv_file}" =~ "version: ${NEW_VERSION}-1" ]]
}

@test "check csv file is set correctly with bundle suffix" {
    # export -f make

    dir="${BATS_TMPDIR}/${BATS_TEST_NAME}"
    cd ${dir}
    run hack/bump-version.sh ${NEW_VERSION} -b 3
    [ "$status" -eq 0 ]

    # Check csv file
    [[ "${output}" =~ "Version bumped from ${CURRENT_VERSION} to ${NEW_VERSION}" ]]
    [[ "${output}" =~ "Set version ${NEW_VERSION} with bundle suffix 3 and replacing version (if not empty) " ]]
    # fine tune results
    csv_file=$(cat $(getCsvFile))
    [[ "${csv_file}" =~ "containerImage: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${csv_file}" =~ "name: rhpam-kogito-operator.v${NEW_VERSION}-3" ]]
    [[ "${csv_file}" =~ "operated-by: rhpam-kogito-operator.${NEW_VERSION}-3" ]]
    [[ "${csv_file}" =~ "version: ${NEW_VERSION}-3" ]]
    bundle_csv_file=$(cat $(getBundleCsvFile))
    [[ "${bundle_csv_file}" =~ "containerImage: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${bundle_csv_file}" =~ "name: rhpam-kogito-operator.v${NEW_VERSION}-3" ]]
    [[ "${bundle_csv_file}" =~ "image: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${bundle_csv_file}" =~ "operated-by: rhpam-kogito-operator.${NEW_VERSION}-3" ]]
    [[ "${bundle_csv_file}" =~ "version: ${NEW_VERSION}-3" ]]
}

@test "check csv file is set correctly with bundle suffix and replaces version" {
    # export -f make

    dir="${BATS_TMPDIR}/${BATS_TEST_NAME}"
    cd ${dir}
    run hack/bump-version.sh ${NEW_VERSION} -b 10 -r 7.11.1
    [ "$status" -eq 0 ]

    # Check csv file
    [[ "${output}" =~ "Version bumped from ${CURRENT_VERSION} to ${NEW_VERSION}" ]]
    [[ "${output}" =~ "Set version ${NEW_VERSION} with bundle suffix 10 and replacing version (if not empty) " ]]
    # fine tune results
    csv_file=$(cat $(getCsvFile))
    [[ "${csv_file}" =~ "containerImage: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${csv_file}" =~ "name: rhpam-kogito-operator.v${NEW_VERSION}-10" ]]
    [[ "${csv_file}" =~ "operated-by: rhpam-kogito-operator.${NEW_VERSION}-10" ]]
    [[ "${csv_file}" =~ "version: ${NEW_VERSION}-10" ]]
    bundle_csv_file=$(cat $(getBundleCsvFile))
    [[ "${bundle_csv_file}" =~ "containerImage: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${bundle_csv_file}" =~ "name: rhpam-kogito-operator.v${NEW_VERSION}-10" ]]
    [[ "${bundle_csv_file}" =~ "image: registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:${NEW_VERSION}" ]]
    [[ "${bundle_csv_file}" =~ "operated-by: rhpam-kogito-operator.${NEW_VERSION}-10" ]]
    [[ "${bundle_csv_file}" =~ "version: ${NEW_VERSION}-10" ]]
}