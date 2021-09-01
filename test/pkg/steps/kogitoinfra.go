// Copyright 2019 Red Hat, Inc. and/or its affiliates
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
	"github.com/cucumber/messages-go/v10"
	kogitoFramework "github.com/kiegroup/kogito-operator/test/pkg/framework"
	kogitoMappers "github.com/kiegroup/kogito-operator/test/pkg/steps/mappers"
	"github.com/kiegroup/rhpam-kogito-operator/test/pkg/framework"
)

/*
	DataTable for Kogito Infra:
	| config          | <key>       | <value>      |
*/

func registerKogitoInfraSteps(ctx *godog.ScenarioContext, data *Data) {
	ctx.Step(`^Install (Kafka) Kogito Infra "([^"]*)" targeting service "([^"]*)" within (\d+) (?:minute|minutes)$`, data.installKogitoInfraTargetingServiceWithinMinutes)
}

func (data *Data) installKogitoInfraTargetingServiceWithinMinutes(targetResourceType, name, targetResourceName string, timeoutInMin int) error {
	return data.installKogitoInfraTargetingServiceWithinMinutesWithConfiguration(targetResourceType, name, targetResourceName, timeoutInMin, &messages.PickleStepArgument_PickleTable{})
}

func (data *Data) installKogitoInfraTargetingServiceWithinMinutesWithConfiguration(targetResourceType, name, targetResourceName string, timeoutInMin int, table *godog.Table) error {
	infraResource, err := framework.GetKogitoInfraResourceStub(data.Namespace, name, targetResourceType, targetResourceName)
	if err != nil {
		return err
	}

	if err := kogitoMappers.MapKogitoInfraTable(table, infraResource); err != nil {
		return err
	}

	kogitoFramework.GetLogger(data.Namespace).Debug("Got kogitoInfra config", "config", infraResource.Spec.InfraProperties)
	err = kogitoFramework.InstallKogitoInfraComponent(data.Namespace, kogitoFramework.GetDefaultInstallerType(), infraResource)
	if err != nil {
		return err
	}

	return kogitoFramework.WaitForKogitoInfraResource(data.Namespace, name, timeoutInMin, framework.GetKogitoInfraResource)
}
