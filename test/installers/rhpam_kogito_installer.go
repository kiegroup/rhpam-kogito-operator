// Copyright 2021 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package installers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kiegroup/kogito-operator/core/client/kubernetes"
	"github.com/kiegroup/kogito-operator/test/config"
	"github.com/kiegroup/kogito-operator/test/framework"
	"github.com/kiegroup/kogito-operator/test/installers"
	v1 "github.com/kiegroup/rhpam-kogito-operator/api/v1"
	"github.com/kiegroup/rhpam-kogito-operator/version"
)

const (
	openShiftInternalRegistryURL             = "image-registry.openshift-image-registry.svc:5000"
	rhpmaKogitoOperatorPullImageSecretPrefix = "rhpam-kogito-operator-controller-manager-dockercfg"
)

var (
	// rhpamKogitoYamlClusterInstaller installs RHPAM Kogito operator cluster wide using YAMLs
	rhpamKogitoYamlClusterInstaller = installers.YamlClusterWideServiceInstaller{
		InstallClusterYaml:               installRhpamKogitoUsingYaml,
		InstallationNamespace:            rhpamKogitoNamespace,
		WaitForClusterYamlServiceRunning: waitForRhpamKogitoOperatorUsingYamlRunning,
		GetAllClusterYamlCrsInNamespace:  getRhpamKogitoCrsInNamespace,
		UninstallClusterYaml:             uninstallRhpamKogitoUsingYaml,
		ClusterYamlServiceName:           rhpamKogitoServiceName,
		CleanupClusterYamlCrsInNamespace: cleanupRhapmKogitoCrsInNamespace,
	}

	// rhpamKogitoOlmNamespacedInstaller installs RHPAM Kogito in the namespace using OLM
	rhpamKogitoOlmNamespacedInstaller = installers.OlmNamespacedServiceInstaller{
		SubscriptionName:                   rhpamKogitoOperatorSubscriptionName,
		Channel:                            rhpamKogitoOperatorSubscriptionChannel,
		Catalog:                            framework.CustomKogitoOperatorCatalog,
		InstallationTimeoutInMinutes:       5,
		GetAllNamespacedOlmCrsInNamespace:  getRhpamKogitoCrsInNamespace,
		CleanupNamespacedOlmCrsInNamespace: cleanupRhapmKogitoCrsInNamespace,
	}

	// rhpamKogitoOlmClusterWideInstaller installs RHPAM Kogito cluster wide using OLM
	rhpamKogitoOlmClusterWideInstaller = installers.OlmClusterWideServiceInstaller{
		SubscriptionName:                    rhpamKogitoOperatorSubscriptionName,
		Channel:                             rhpamKogitoOperatorSubscriptionChannel,
		Catalog:                             framework.CustomKogitoOperatorCatalog,
		InstallationTimeoutInMinutes:        5,
		GetAllClusterWideOlmCrsInNamespace:  getRhpamKogitoCrsInNamespace,
		CleanupClusterWideOlmCrsInNamespace: cleanupRhapmKogitoCrsInNamespace,
	}

	rhpamKogitoNamespace            = "rhpam-kogito-operator-system"
	rhpamKogitoServiceName          = "RHPAM Kogito operator"
	rhpamKogitoOperatorTimeoutInMin = 5
	rhpamKogitoImageStreamName      = "rhpam-kogito-operator"

	rhpamKogitoOperatorSubscriptionName    = "rhpam-kogito-operator"
	rhpamKogitoOperatorSubscriptionChannel = "7.x"
)

// GetRhpamKogitoInstaller returns RHPAM Kogito installer
func GetRhpamKogitoInstaller() (installers.ServiceInstaller, error) {
	if config.IsOperatorInstalledByYaml() {
		if config.IsOperatorNamespaced() {
			return nil, errors.New("Installing namespace scoped RHPAM Kogito operator using YAML files is not supported")
		}
		return &rhpamKogitoYamlClusterInstaller, nil
	}

	if config.IsOperatorInstalledByOlm() {
		if config.IsOperatorNamespaced() {
			return &rhpamKogitoOlmNamespacedInstaller, nil
		}
		return &rhpamKogitoOlmClusterWideInstaller, nil
	}

	return nil, errors.New("No RHPAM Kogito operator installer available for provided configuration")
}

func installRhpamKogitoUsingYaml() error {
	framework.GetMainLogger().Info("Installing RHPAM Kogito operator")

	// Create namespace first so ImageStream can be placed there
	if err := framework.CreateNamespace(rhpamKogitoNamespace); err != nil {
		return err
	}

	yamlContent, err := framework.ReadFromURI(config.GetOperatorYamlURI())
	if err != nil {
		framework.GetMainLogger().Error(err, "Error while reading kogito-operator.yaml file")
		return err
	}

	// Use insecure ImageStream when deploying on OpenShift to support using insecure registries, unless the operator tag already points to internal registry
	if framework.IsOpenshift() && !strings.Contains(framework.GetOperatorImageNameAndTag(), openShiftInternalRegistryURL) {
		if err := framework.CreateInsecureImageStream(rhpamKogitoNamespace, rhpamKogitoImageStreamName, config.GetOperatorImageTag(), framework.GetOperatorImageNameAndTag()); err != nil {
			return err
		}

		rhpamKogitoInternalImageTagName := fmt.Sprintf("%s/%s/%s:%s", openShiftInternalRegistryURL, rhpamKogitoNamespace, rhpamKogitoImageStreamName, config.GetOperatorImageTag())
		yamlContent = strings.ReplaceAll(yamlContent, "quay.io/kiegroup/rhpam-kogito-operator:"+version.Version, rhpamKogitoInternalImageTagName)
	} else {
		yamlContent = strings.ReplaceAll(yamlContent, "quay.io/kiegroup/rhpam-kogito-operator:"+version.Version, framework.GetOperatorImageNameAndTag())
	}

	tempFilePath, err := framework.CreateTemporaryFile("kogito-operator*.yaml", yamlContent)
	if err != nil {
		framework.GetMainLogger().Error(err, "Error while storing adjusted YAML content to temporary file")
		return err
	}

	_, err = framework.CreateCommand("oc", "apply", "-f", tempFilePath).Execute()
	if err != nil {
		framework.GetMainLogger().Error(err, "Error while installing RHPAM Kogito operator from YAML file")
		return err
	}

	return nil
}

func waitForRhpamKogitoOperatorUsingYamlRunning() error {
	return framework.WaitForOnOpenshift(rhpamKogitoNamespace, "RHPAM Kogito operator running", rhpamKogitoOperatorTimeoutInMin,
		func() (bool, error) {
			podList, err := framework.GetPods(rhpamKogitoNamespace)
			if err != nil {
				framework.GetLogger(rhpamKogitoNamespace).Error(err, "Error while trying to retrieve RHPAM Kogito Operator pods")
				return false, nil
			}
			if len(podList.Items) != 1 {
				return false, nil
			}

			running := framework.CheckPodsAreReady(podList)

			// If not running, make sure the image pull secret is present in pod
			// If not present, delete the pod to allow its reconstruction with correct pull secret
			// Note that this is specific to Openshift
			if !running && framework.IsOpenshift() {
				for _, pod := range podList.Items {
					if !framework.CheckPodHasImagePullSecretWithPrefix(&pod, rhpmaKogitoOperatorPullImageSecretPrefix) {
						// Delete pod as it has been misconfigured (missing pull secret)
						framework.GetLogger(rhpamKogitoNamespace).Info("RHPAM Kogito Operator pod does not have the image pull secret needed. Deleting it to renew it.")
						err := framework.DeleteObject(&pod)
						if err != nil {
							framework.GetLogger(rhpamKogitoNamespace).Error(err, "Error while trying to delete RHPAM Kogito Operator pod")
							return false, nil
						}
					}
				}
			}
			return running, nil
		})
}

func uninstallRhpamKogitoUsingYaml() error {
	framework.GetMainLogger().Info("Uninstalling Kogito operator")

	output, err := framework.CreateCommand("oc", "delete", "-f", config.GetOperatorYamlURI(), "--timeout=30s").Execute()
	if err != nil {
		framework.GetMainLogger().Error(err, fmt.Sprintf("Deleting RHPAM Kogito operator failed, output: %s", output))
		return err
	}

	return nil
}

func getRhpamKogitoCrsInNamespace(namespace string) ([]kubernetes.ResourceObject, error) {
	crs := []kubernetes.ResourceObject{}

	kogitoRuntimes := &v1.KogitoRuntimeList{}
	if err := framework.GetObjectsInNamespace(namespace, kogitoRuntimes); err != nil {
		return nil, err
	}
	for i := range kogitoRuntimes.Items {
		crs = append(crs, &kogitoRuntimes.Items[i])
	}

	kogitoBuilds := &v1.KogitoBuildList{}
	if err := framework.GetObjectsInNamespace(namespace, kogitoBuilds); err != nil {
		return nil, err
	}
	for i := range kogitoBuilds.Items {
		crs = append(crs, &kogitoBuilds.Items[i])
	}

	return crs, nil
}

func cleanupRhapmKogitoCrsInNamespace(namespace string) bool {
	crs, err := getRhpamKogitoCrsInNamespace(namespace)
	if err != nil {
		framework.GetLogger(namespace).Error(err, "Error getting RHPAM Kogito CRs.")
		return false
	}

	for _, cr := range crs {
		if err := framework.DeleteObject(cr); err != nil {
			framework.GetLogger(namespace).Error(err, "Error deleting RHPAM Kogito CR.", "CR name", cr.GetName())
			return false
		}
	}
	return true
}
