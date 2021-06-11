export OPERATOR_NAME=quay.io/kiegroup/rhpam-kogito-operator-nightly
export IMAGE_TAG=temp-main-tests

export OPERATOR_NAME=default-route-openshift-image-registry.apps-crc.testing/openshift/rhpam-kogito-operator
export OPERATOR_NAME=image-registry.openshift-image-registry.svc:5000/openshift/rhpam-kogito-operator
export IMAGE_TAG=kogito-4858

make build-examples-images load_default_config=true local=true load_factor=3 disable_maven_native_build_container=true operator_image=${OPERATOR_NAME} operator_tag=${IMAGE_TAG} cr_deployment_only=true image_cache_mode=never runtime_application_image_registry=default-route-openshift-image-registry.apps-crc.testing runtime_application_image_namespace=openshift runtime_application_image_version=${IMAGE_TAG} container_engine=podman tags=@rhpam

make run-tests load_default_config=true local=true load_factor=3 disable_maven_native_build_container=true operator_image=${OPERATOR_NAME} operator_tag=${IMAGE_TAG} cr_deployment_only=true image_cache_mode=always runtime_application_image_registry=image-registry.openshift-image-registry.svc:5000 runtime_application_image_namespace=openshift runtime_application_image_version=${IMAGE_TAG} container_engine=podman tags=@rhpam