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
	"github.com/cucumber/godog"
	"github.com/kiegroup/rhpam-kogito-operator/test/pkg/installers"
)

func registerKafkaSteps(ctx *godog.ScenarioContext, data *Data) {
	ctx.Step(`^Kafka Operator is deployed$`, data.kafkaOperatorIsDeployed)
}

func (data *Data) kafkaOperatorIsDeployed() error {
	return installers.GetAmqStreamsInstaller().Install(data.Namespace)
}
