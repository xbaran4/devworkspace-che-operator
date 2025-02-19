//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CheManagerSpec holds the configuration of the Che controller.
// +k8s:openapi-gen=true
type CheManagerSpec struct {
	// GatewayHost is the full host name used to expose devworkspace endpoints that support url rewriting reverse proxy.
	// See the GatewayDisabled attribute for a more detailed description of where and how are devworkspace endpoints
	// exposed in various configurations.
	//
	// This attribute is mandatory on Kubernetes, optional on OpenShift.
	GatewayHost string `json:"gatewayHost,omitempty"`

	// GatewayDisabled enables or disables routing of the url rewrite supporting devworkspace endpoints
	// through a common gateway (the hostname of which is defined by the GatewayHost).
	//
	// Default value is "false" meaning that the gateway is enabled.
	//
	// If set to false (i.e. the gateway is enabled), endpoints marked using the "urlRewriteSupported" attribute
	// are exposed on unique subpaths of the GatewayHost, while the rest of the devworkspace endpoints are exposed
	// on subdomains of the RoutingSuffix specified by the DevWorkspaceRouting of the devworkspace.
	//
	// If set to true (i.e. the gateway is disabled), all endpoints are deployed on subdomains of
	// the RoutingSuffix.
	GatewayDisabled bool `json:"gatewayDisabled,omitempty"`

	// GatewayImage is the docker image to use for the Che gateway.  This is only used if GatewayDisabled is false.
	// If not defined in the CR, it is taken from
	// the `RELATED_IMAGE_gateway` environment variable of the che operator
	// deployment/pod. If not defined there, it defaults to a hardcoded value.
	GatewayImage string `json:"gatewayImage,omitempty"`

	// GatewayConfigurerImage is the docker image to use for the sidecar of the Che gateway that is
	// used to configure it. This is only used when GatewayDisabled is false. If not defined in the CR,
	// it is taken from the `RELATED_IMAGE_gateway_configurer` environment variable of the che
	// operator deployment/pod. If not defined there, it defaults to a hardcoded value.
	GatewayConfigurerImage string `json:"gatewayConfigurerImage,omitempty"`

	// Name of a secret that will be used to setup ingress/route TLS certificate.
	// When the field is empty string, the default cluster certificate will be used.
	// The same secret is assumed to exist in the same namespace as the CheManager CR and is used for both
	// the gateway and all devworkspace endpoints.
	// In case of the devworkspace endpoints, the secret is copied to the namespace of the devworkspace.
	//
	// The secret has to be of type "tls".
	//
	// +optional
	TlsSecretName string `json:"tlsSecretName,omitempty"`

	// K8s contains the configuration specific only to Kubernetes
	K8s CheManagerSpecK8s `json:"k8s,omitempty"`
}

// CheManagerSpecK8s contains the configuration options specific to Kubernetes only.
type CheManagerSpecK8s struct {
	// IngressAnnotations are the annotations to be put on the generated ingresses. This can be used to
	// configure the ingress class and the ingress-controller-specific behavior for both the gateway
	// and the ingresses created to expose the Devworkspace component endpoints.
	// When not specified, this defaults to:
	//
	//     kubernetes.io/ingress.class:                       "nginx"
	//     nginx.ingress.kubernetes.io/proxy-read-timeout:    "3600",
	//     nginx.ingress.kubernetes.io/proxy-connect-timeout: "3600",
	//     nginx.ingress.kubernetes.io/ssl-redirect:          "true"
	//
	// +optional
	IngressAnnotations map[string]string `json:"ingressAnnotations,omitempty"`
}

type GatewayPhase string

const (
	GatewayPhaseInitializing = "Initializing"
	GatewayPhaseEstablished  = "Established"
	GatewayPhaseInactive     = "Inactive"
)

type ManagerPhase string

const (
	ManagerPhaseActive          = "Active"
	ManagerPhaseInactive        = "Inactive"
	ManagerPhasePendingDeletion = "PendingDeletion"
)

// +k8s:openapi-gen=true
type CheManagerStatus struct {
	// GatewayPhase specifies the phase in which the singlehost gateway deployment currently is.
	// If the manager routing is not singlehost, this is "Inactive"
	GatewayPhase GatewayPhase `json:"gatewayPhase,omitempty"`

	// GatewayHost is the resolved host of the ingress/route, on which the gateway is accessible.
	GatewayHost string `json:"gatewayHost,omitempty"`

	// Phase is the phase in which the manager as a whole finds itself in.
	Phase ManagerPhase `json:"phase,omitempty"`

	// Message contains further human-readable info for why the manager is in the phase it currently is.
	Message string `json:"message,omitempty"`
}

// CheManager is the configuration of the CheManager layer of Devworkspace.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=chemanagers,scope=Namespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CheManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CheManagerSpec   `json:"spec,omitempty"`
	Status CheManagerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CheManagerList is the list type for Che
type CheManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CheManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CheManager{}, &CheManagerList{})
}
