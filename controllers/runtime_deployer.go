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

package controllers

import (
	"reflect"

	api "github.com/kiegroup/kogito-operator/apis"
	"github.com/kiegroup/kogito-operator/core/framework/util"
	"github.com/kiegroup/kogito-operator/core/infrastructure"
	"github.com/kiegroup/kogito-operator/core/manager"
	"github.com/kiegroup/kogito-operator/core/operator"
	"github.com/kiegroup/rhpam-kogito-operator/internal"

	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"github.com/kiegroup/kogito-operator/core/framework"
	monv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "k8s.io/api/apps/v1"
)

const (
	envVarExternalURL = "KOGITO_SERVICE_URL"
	envVarNamespace   = "NAMESPACE"
)

// RuntimeDeployerHandler ...
type RuntimeDeployerHandler interface {
	OnGetComparators(comparator compare.ResourceComparator)
	OnDeploymentCreate(deployment *v1.Deployment) error
}

type runtimeDeployerHandler struct {
	operator.Context
	instance       api.KogitoRuntimeInterface
	runtimeHandler manager.KogitoRuntimeHandler
}

// NewRuntimeDeployerHandler ...
func NewRuntimeDeployerHandler(context operator.Context, instance api.KogitoRuntimeInterface, runtimeHandler manager.KogitoRuntimeHandler) RuntimeDeployerHandler {
	return &runtimeDeployerHandler{
		Context:        context,
		instance:       instance,
		runtimeHandler: runtimeHandler,
	}
}

func (d *runtimeDeployerHandler) OnGetComparators(comparator compare.ResourceComparator) {
	comparator.SetComparator(
		framework.NewComparatorBuilder().
			WithType(reflect.TypeOf(monv1.ServiceMonitor{})).
			WithCustomComparator(framework.CreateServiceMonitorComparator()).
			Build())
}

// onDeploymentCreate hooks into the infrastructure package to add additional capabilities/properties to the deployment creation
func (d *runtimeDeployerHandler) OnDeploymentCreate(deployment *v1.Deployment) error {
	// NAMESPACE service discovery
	framework.SetEnvVar(envVarNamespace, d.instance.GetNamespace(), &deployment.Spec.Template.Spec.Containers[0])
	// external URL
	if d.instance.GetStatus().GetExternalURI() != "" {
		framework.SetEnvVar(envVarExternalURL, d.instance.GetStatus().GetExternalURI(), &deployment.Spec.Template.Spec.Containers[0])
	}
	// sa
	deployment.Spec.Template.Spec.ServiceAccountName = infrastructure.RuntimeServiceAccountName
	// metering labels for product operator
	util.AppendToStringMap(internal.GetMeteringLabels(), deployment.Spec.Template.Labels)
	// istio
	if d.instance.GetRuntimeSpec().IsEnableIstio() {
		framework.AddIstioInjectSidecarAnnotation(&deployment.Spec.Template.ObjectMeta)
	}
	return nil
}
