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

package test

import (
	"testing"

	"github.com/cucumber/godog"
	kogitoTest "github.com/kiegroup/kogito-operator/test"
	kogitoFramework "github.com/kiegroup/kogito-operator/test/framework"
	kogitoSteps "github.com/kiegroup/kogito-operator/test/steps"
	"github.com/kiegroup/rhpam-kogito-operator/meta"
	"github.com/kiegroup/rhpam-kogito-operator/test/steps"
)

func TestMain(m *testing.M) {
	if err := kogitoFramework.InitKubeClient(meta.GetRegisteredSchema()); err != nil {
		panic(err)
	}

	kogitoTest.PreRegisterStepsHook = func(ctx *godog.ScenarioContext, d *kogitoSteps.Data) {
		data := &steps.Data{Data: d}
		data.RegisterAllSteps(ctx)
	}

	kogitoTest.AfterScenarioHook = func(scenario *godog.Scenario, d *kogitoSteps.Data) error {
		data := &steps.Data{Data: d}
		return data.AfterScenario(scenario)
	}

	kogitoTest.ExecuteTests()
}
