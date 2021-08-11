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
	kogitoExecutor "github.com/kiegroup/kogito-operator/test/pkg/executor"
	kogitoFramework "github.com/kiegroup/kogito-operator/test/pkg/framework"
	kogitoSteps "github.com/kiegroup/kogito-operator/test/pkg/steps"
	v1 "github.com/kiegroup/rhpam-kogito-operator/api/v1"
	"github.com/kiegroup/rhpam-kogito-operator/meta"
	"github.com/kiegroup/rhpam-kogito-operator/test/pkg/steps"

	imgv1 "github.com/openshift/api/image/v1"
	olmapiv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
)

func TestMain(m *testing.M) {
	if err := kogitoFramework.InitKubeClient(meta.GetRegisteredSchema()); err != nil {
		panic(err)
	}

	kogitoExecutor.PreRegisterStepsHook = func(ctx *godog.ScenarioContext, d *kogitoSteps.Data) {
		data := &steps.Data{Data: d}
		data.RegisterAllSteps(ctx)
		data.RegisterLogsKubernetesObjects(&imgv1.ImageStreamList{}, &v1.KogitoRuntimeList{}, &v1.KogitoBuildList{}, &olmapiv1alpha1.ClusterServiceVersionList{})
	}
	kogitoExecutor.DisableLogsKogitoCommunityObjects()

	kogitoExecutor.ExecuteBDDTests(nil)
}
