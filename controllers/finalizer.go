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
	"context"
	"fmt"
	oauthv1alpha1 "oauth-proxy-operator/api/v1alpha1"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	finalizer = "oauthproxy.finalizers.emdinix.io"
)

func (r *OAuthProxyReconciler) finalizer(request ctrl.Request, cr *oauthv1alpha1.OAuthProxy) (*ctrl.Result, error) {

	// examine DeletionTimestamp to determine if object is under deletion
	if cr.ObjectMeta.DeletionTimestamp.IsZero() {
		// OAuthProxy not being deleted
		// Add the finalizer if not already there
		if !containsString(cr.ObjectMeta.Finalizers, finalizer) {
			cr.ObjectMeta.Finalizers = append(cr.ObjectMeta.Finalizers, finalizer)
			if err := r.Client.Update(context.Background(), cr); err != nil {
				return &ctrl.Result{}, err
			}
		}
	} else {
		// OAuthProxy being deleted
		if containsString(cr.ObjectMeta.Finalizers, finalizer) {
			// OAuthProxy finalizer present, handle pre delete

			switch cr.Spec.ResourceKind {
			case "DeploymentConfig":
				result, err := r.finalizeDeploymentConfig(cr)
				if result != nil {
					return result, err
				}
				fmt.Printf("\nDeploymentConfig will be modified\n")
			case "Deployment":
				result, err := r.finalizeDeployment(cr)
				if result != nil {
					return result, err
				}
			case "DaemonSet":
				result, err := r.finalizeDaemonSet(cr)
				if result != nil {
					return result, err
				}
			case "StatefulSet":
				result, err := r.finalizeStatefulSet(cr)
				if result != nil {
					return result, err
				}
			}
		}
	}
	// Stop reconciliation as the item is being deleted
	return &ctrl.Result{}, nil
}
