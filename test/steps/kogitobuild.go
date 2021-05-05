// Copyright 2020 Red Hat, Inc. and/or its affiliates
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

package steps

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/kiegroup/kogito-operator/api"
	v1 "github.com/kiegroup/rhpam-kogito-operator/api/v1"
	"github.com/kiegroup/rhpam-kogito-operator/test/framework"

	"github.com/cucumber/godog"
	"github.com/kiegroup/kogito-operator/test/config"
	communityFramework "github.com/kiegroup/kogito-operator/test/framework"
	bddtypes "github.com/kiegroup/kogito-operator/test/types"
	"github.com/kiegroup/rhpam-kogito-operator/test/steps/mappers"
)

/*
	DataTable for KogitoBuild:
	| config        | native     | enabled/disabled |
	| build-request | cpu/memory | value            |
	| build-limit   | cpu/memory | value            |
*/

const defaultTimeoutToStartBuildInMin = 5

func registerKogitoBuildSteps(ctx *godog.ScenarioContext, data *Data) {
	// Deploy steps
	ctx.Step(`^Build (quarkus|springboot) example service "([^"]*)" with configuration:$`, data.buildExampleServiceWithConfiguration)
	ctx.Step(`^Build binary (quarkus|springboot) service "([^"]*)" with configuration:$`, data.buildBinaryServiceWithConfiguration)
	ctx.Step(`^Build binary (quarkus|springboot) local example service "([^"]*)" from target folder with configuration:$`, data.buildBinaryLocalExampleServiceFromTargetFolderWithConfiguration)
}

// Build service steps

func (data *Data) buildExampleServiceWithConfiguration(runtimeType, contextDir string, table *godog.Table) error {
	buildHolder, err := getKogitoBuildConfiguredStub(data.Namespace, runtimeType, filepath.Base(contextDir), table)
	if err != nil {
		return err
	}

	buildHolder.KogitoBuild.GetSpec().SetType(api.RemoteSourceBuildType)
	buildHolder.KogitoBuild.GetSpec().GetGitSource().SetURI(config.GetExamplesRepositoryURI())
	buildHolder.KogitoBuild.GetSpec().GetGitSource().SetContextDir(contextDir)
	if ref := config.GetExamplesRepositoryRef(); len(ref) > 0 {
		buildHolder.KogitoBuild.GetSpec().GetGitSource().SetReference(ref)
	}

	if err := framework.DeployKogitoBuild(data.Namespace, buildHolder); err != nil {
		return err
	}

	// In case of OpenShift the ImageStream needs to be patched to allow insecure registries
	if communityFramework.IsOpenshift() {
		if err := makeImageStreamInsecure(data.Namespace, framework.GetKogitoBuildS2IImage()); err != nil {
			return err
		}
		if err := makeImageStreamInsecure(data.Namespace, framework.GetKogitoBuildRuntimeImage(buildHolder.KogitoBuild.(*v1.KogitoBuild))); err != nil {
			return err
		}
	}

	return nil
}

func (data *Data) buildBinaryServiceWithConfiguration(runtimeType, serviceName string, table *godog.Table) error {
	buildHolder, err := getKogitoBuildConfiguredStub(data.Namespace, runtimeType, serviceName, table)
	if err != nil {
		return err
	}

	buildHolder.KogitoBuild.GetSpec().SetType(api.BinaryBuildType)

	if err := framework.DeployKogitoBuild(data.Namespace, buildHolder); err != nil {
		return err
	}

	// In case of OpenShift the ImageStream needs to be patched to allow insecure registries
	if communityFramework.IsOpenshift() {
		if err := makeImageStreamInsecure(data.Namespace, framework.GetKogitoBuildRuntimeImage(buildHolder.KogitoBuild.(*v1.KogitoBuild))); err != nil {
			return err
		}
	}

	return nil
}

func (data *Data) buildBinaryLocalExampleServiceFromTargetFolderWithConfiguration(runtimeType, serviceName string, table *godog.Table) error {
	buildHolder, err := getKogitoBuildConfiguredStub(data.Namespace, runtimeType, serviceName, table)
	if err != nil {
		return err
	}

	buildHolder.KogitoBuild.GetSpec().SetType(api.BinaryBuildType)
	buildHolder.BuiltBinaryFolder = fmt.Sprintf(`%s/%s/target`, data.KogitoExamplesLocation, serviceName)

	err = framework.DeployKogitoBuild(data.Namespace, buildHolder)
	if err != nil {
		return err
	}

	// In case of OpenShift the ImageStream needs to be patched to allow insecure registries
	if communityFramework.IsOpenshift() {
		if err := makeImageStreamInsecure(data.Namespace, framework.GetKogitoBuildRuntimeImage(buildHolder.KogitoBuild.(*v1.KogitoBuild))); err != nil {
			return err
		}
	}

	// If we don't use Kogito CLI then upload target folder using OC client
	return communityFramework.WaitForOnOpenshift(data.Namespace, fmt.Sprintf("Build '%s' to start", serviceName), defaultTimeoutToStartBuildInMin,
		func() (bool, error) {
			_, err := communityFramework.CreateCommand("oc", "start-build", serviceName, "--from-dir="+buildHolder.BuiltBinaryFolder, "-n", data.Namespace).WithLoggerContext(data.Namespace).Execute()
			return err == nil, err
		})
}

// Misc methods

// getKogitoBuildConfiguredStub Get KogitoBuildHolder initialized from table if provided
func getKogitoBuildConfiguredStub(namespace, runtimeType, serviceName string, table *godog.Table) (buildHolder *bddtypes.KogitoBuildHolder, err error) {
	kogitoBuild := framework.GetKogitoBuildStub(namespace, runtimeType, serviceName)
	kogitoRuntime := framework.GetKogitoRuntimeStub(namespace, runtimeType, serviceName, "")

	buildHolder = &bddtypes.KogitoBuildHolder{
		KogitoServiceHolder: &bddtypes.KogitoServiceHolder{KogitoService: kogitoRuntime},
		KogitoBuild:         kogitoBuild,
	}

	if table != nil {
		err = mappers.MapKogitoBuildTable(table, buildHolder)
	}

	framework.SetupKogitoBuildImageStreams(kogitoBuild)

	return buildHolder, err
}

func makeImageStreamInsecure(namespace, insecureImageTag string) error {
	// Need to wait as operator overrides image stream in initial reconciliation
	time.Sleep(time.Duration(2*config.GetLoadFactor()) * time.Second)
	return communityFramework.WaitForOnOpenshift(namespace, fmt.Sprintf("patching ImageStream pointing to %s to be insecure", insecureImageTag), 2, func() (bool, error) {
		imageStreams, err := communityFramework.GetImageStreams(namespace)
		if err != nil {
			return false, err
		}
		for _, is := range imageStreams.Items {
			imageStream := &is
			for i, tag := range imageStream.Spec.Tags {
				if tag.From != nil && strings.Contains(tag.From.Name, insecureImageTag) {
					// Image tag has to be removed and created again, to trigger image fetch
					imageStream.Spec.Tags = append(imageStream.Spec.Tags[:i], imageStream.Spec.Tags[i+1:]...)
					if err := communityFramework.UpdateObject(imageStream); err != nil {
						return false, err
					}
					tag.ImportPolicy.Insecure = true
					imageStream.Spec.Tags = append(imageStream.Spec.Tags, tag)
					if err := communityFramework.UpdateObject(imageStream); err != nil {
						return false, err
					}
					return true, nil
				}
			}
		}
		return false, nil
	})
}
