module github.com/kiegroup/rhpam-kogito-operator

go 1.14

require (
	github.com/RHsyseng/operator-utils v0.0.0-20200506183821-e3b4a2ba9c30
	github.com/cucumber/godog v0.11.0
	github.com/infinispan/infinispan-operator v0.0.0-20210621093106-4662500f4ae1
	github.com/integr8ly/grafana-operator/v3 v3.10.0
	github.com/keycloak/keycloak-operator v0.0.0-20200917060808-9858b19ca8bf
	github.com/kiegroup/kogito-operator v0.12.1-0.20210725113612-77743fbb39d7
	github.com/kiegroup/kogito-operator/api v0.0.0-20210725113612-77743fbb39d7
	github.com/kiegroup/rhpam-kogito-operator/api v0.0.0-00010101000000-000000000000
	github.com/mongodb/mongodb-kubernetes-operator v0.3.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20200321030439-57b580e57e88
	github.com/operator-framework/operator-marketplace v0.0.0-20190919183128-4ef67b2f50e9
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.46.0
	k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v12.0.0+incompatible
	knative.dev/eventing v0.18.0
	sigs.k8s.io/controller-runtime v0.6.3
)

// local modules
replace github.com/kiegroup/rhpam-kogito-operator/api => ./api

replace (
	k8s.io/api => k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.4
	k8s.io/apiserver => k8s.io/apiserver v0.18.8
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.20.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.8
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.8
	k8s.io/code-generator => k8s.io/code-generator v0.18.8
	k8s.io/component-base => k8s.io/component-base v0.18.8
	k8s.io/cri-api => k8s.io/cri-api v0.18.8
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.8
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.8
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.8
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.8
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.8
	k8s.io/kubectl => k8s.io/kubectl v0.18.8
	k8s.io/kubelet => k8s.io/kubelet v0.18.8
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.8
	k8s.io/metrics => k8s.io/metrics v0.18.8
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.8
)

replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200623075207-eb651a5bb0ad
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200623090625-83993cebb5ae
)

// fix Azure bogus https://github.com/kubernetes/client-go/issues/628
replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM

// Required by Helm
replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

// add kogito-operator/[api|client] modules
replace (
	github.com/kiegroup/kogito-operator => github.com/kiegroup/kogito-operator v0.12.1-0.20210707114619-6ae2f22b898a
	github.com/kiegroup/kogito-operator/api => github.com/kiegroup/kogito-operator/api v0.0.0-20210707114619-6ae2f22b898a
	github.com/kiegroup/kogito-operator/client => github.com/kiegroup/kogito-operator/client v0.0.0-20210707114619-6ae2f22b898a
)
