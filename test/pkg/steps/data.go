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

package steps

import (
	"github.com/cucumber/godog"
	"github.com/kiegroup/kogito-operator/test/pkg/framework"
	"github.com/kiegroup/kogito-operator/test/pkg/steps"
	v1 "github.com/kiegroup/rhpam-kogito-operator/api/v1"
)

// Data contains all data needed by Gherkin steps to run
type Data struct {
	*steps.Data
}

// RegisterAllSteps register all steps available to the test suite
func (data *Data) RegisterAllSteps(ctx *godog.ScenarioContext) {
	registerKogitoBuildSteps(ctx, data)
	registerKogitoDeployFilesSteps(ctx, data)
	registerKogitoRuntimeSteps(ctx, data)
	registerOpenShiftSteps(ctx, data)
	registerOperatorSteps(ctx, data)
}

// AfterScenario executes some actions on data after a scenario is finished
func (data *Data) AfterScenario(scenario *godog.Scenario) error {
	error := framework.OperateOnNamespaceIfExists(data.Namespace, func(namespace string) error {
		if err := framework.LogKubernetesObjects(data.Namespace, &v1.KogitoRuntimeList{}, &v1.KogitoBuildList{}); err != nil {
			framework.GetMainLogger().Error(err, "Error logging Kubernetes objects", "namespace", namespace)
		}
		return nil
	})

	if error != nil {
		return error
	}

	return nil
}
