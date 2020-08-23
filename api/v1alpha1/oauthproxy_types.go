/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OAuthProxySpec defines the desired state of OAuthProxy
type OAuthProxySpec struct {
	// Namespace where Custom Resource should take actions
	Namespace string `json:"namespace"`
	// Resource to protect by OAuthProxy
	ResourceName string `json:"resourcename"`
	// Resource Kind can be Deployment, DeyplomentConfig, DaemonSet or StatefulSet
	ResourceKind string `json:"resourcekind"`
	// Defines the upstream port which will be protected
	UpstreamPort uint16 `json:"upstreamport"`
}

// OAuthProxyStatus defines the observed state of OAuthProxy
type OAuthProxyStatus struct {
	Message string `json:"message"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=oauthproxies,scope=Cluster

// OAuthProxy is the Schema for the oauthproxies API
type OAuthProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OAuthProxySpec   `json:"spec,omitempty"`
	Status OAuthProxyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OAuthProxyList contains a list of OAuthProxy
type OAuthProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OAuthProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OAuthProxy{}, &OAuthProxyList{})
}
