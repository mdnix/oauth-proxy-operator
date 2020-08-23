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

package controllers

import (
	oauthv1alpha1 "oauth-proxy-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OAuthProxyReconciler) OAuthProxyService(cr *oauthv1alpha1.OAuthProxy) *corev1.Service {
	labels := map[string]string{
		"app": cr.Spec.ResourceName,
	}
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ResourceName,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"service.alpha.openshift.io/serving-cert-secret-name": cr.Spec.ResourceName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       443,
				TargetPort: intstr.FromInt(8443),
				Name:       "oauth-proxy",
			}},
			ClusterIP: "None",
		},
	}
	ctrl.SetControllerReference(cr, s, r.Scheme)
	return s

}
