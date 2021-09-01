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

package framework

import (
	"fmt"

	"github.com/kiegroup/kogito-operator/api"
	"github.com/kiegroup/kogito-operator/core/infrastructure"

	"github.com/kiegroup/kogito-operator/test/pkg/framework"
	kogitoFramework "github.com/kiegroup/kogito-operator/test/pkg/framework"
	v1 "github.com/kiegroup/rhpam-kogito-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetKogitoInfraResourceStub Get basic KogitoInfra stub with all needed fields initialized
func GetKogitoInfraResourceStub(namespace, name, targetResourceType, targetResourceName string) (*v1.KogitoInfra, error) {
	infraResource, err := parseKogitoInfraResource(targetResourceType)
	if err != nil {
		return nil, err
	}
	infraResource.SetName(targetResourceName)

	return &v1.KogitoInfra{
		ObjectMeta: framework.NewObjectMetadata(namespace, name),
		Spec: v1.KogitoInfraSpec{
			Resource: *infraResource,
		},
		Status: v1.KogitoInfraStatus{
			Conditions: &[]apiv1.Condition{
				{
					LastTransitionTime: apiv1.Now(),
				},
			},
		},
	}, nil
}

// Converts infra resource from name to InfraResource struct
func parseKogitoInfraResource(targetResourceType string) (*v1.InfraResource, error) {
	switch targetResourceType {
	// case infrastructure.InfinispanKind:
	// 	return &v1.InfraResource{APIVersion: infrastructure.InfinispanAPIVersion, Kind: infrastructure.InfinispanKind}, nil
	case infrastructure.KafkaKind:
		return &v1.InfraResource{APIVersion: infrastructure.KafkaAPIVersion, Kind: infrastructure.KafkaKind}, nil
	// case infrastructure.KeycloakKind:
	// 	return &v1.InfraResource{APIVersion: infrastructure.KeycloakAPIVersion, Kind: infrastructure.KeycloakKind}, nil
	// case infrastructure.MongoDBKind:
	// 	return &v1.InfraResource{APIVersion: infrastructure.MongoDBAPIVersion, Kind: infrastructure.MongoDBKind}, nil
	// case infrastructure.KnativeEventingBrokerKind:
	// 	return &v1.InfraResource{APIVersion: infrastructure.KnativeEventingAPIVersion, Kind: infrastructure.KnativeEventingBrokerKind}, nil
	default:
		return nil, fmt.Errorf("Unknown KogitoInfra target resource type %s", targetResourceType)
	}
}

// GetKogitoInfraResource retrieves the KogitoInfra resource
func GetKogitoInfraResource(namespace, name string) (api.KogitoInfraInterface, error) {
	infraResource := &v1.KogitoInfra{}
	if exists, err := kogitoFramework.GetObjectWithKey(types.NamespacedName{Name: name, Namespace: namespace}, infraResource); err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("Error while trying to look for KogitoInfra %s: %v ", name, err)
	} else if !exists {
		return nil, nil
	}
	return infraResource, nil
}
