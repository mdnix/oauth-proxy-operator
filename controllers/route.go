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

	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OAuthProxyReconciler) OAuthProxyRoute(cr *oauthv1alpha1.OAuthProxy) *routev1.Route {
	labels := map[string]string{
		"app": cr.Spec.ResourceName,
	}

	ro := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ResourceName,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: cr.Spec.ResourceName,
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationReencrypt,
			},
		},
	}
	ctrl.SetControllerReference(cr, ro, r.Scheme)

	return ro
}
