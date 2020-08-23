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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	oauthv1alpha1 "oauth-proxy-operator/api/v1alpha1"

	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
)

// OAuthProxyReconciler reconciles a OAuthProxy object
type OAuthProxyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=oauth.emdinix.io,resources=oauthproxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=oauth.emdinix.io,resources=oauthproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.openshift.io,resources=deploymentconfigs,verbs=get;list;watch;create;update;patch;delete

func (r *OAuthProxyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("oauthproxy", req.NamespacedName)

	// Fetch the OAuthProxy instance
	instance := &oauthv1alpha1.OAuthProxy{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// req object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the req.
		return ctrl.Result{}, err
	}

	var result *ctrl.Result

	result, err = r.finalizer(req, instance)
	if err != nil {
		return *result, err
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {

		result, err = r.ensureServiceAccount(req, instance, r.OAuthProxyServiceAccount(instance))
		if result != nil {
			return *result, err
		}

		result, err = r.ensureService(req, instance, r.OAuthProxyService(instance))
		if result != nil {
			return *result, err
		}

		result, err = r.ensureRoute(req, instance, r.OAuthProxyRoute(instance))
		if result != nil {
			return *result, err
		}

		switch instance.Spec.ResourceKind {
		case "DeploymentConfig":
			result, err = r.handleDeploymentConfigChanges(instance)
			if err != nil {
				return *result, err
			}
		case "Deployment":
			result, err = r.handleDeploymentChanges(instance)
			if result != nil {
				return *result, err
			}
		case "DaemonSet":
			result, err = r.handleDaemonSetChanges(instance)
			if result != nil {
				return *result, err
			}
		case "StatefulSet":
			result, err = r.handleStatefulSetChanges(instance)
			if result != nil {
				return *result, err
			}
		}

		msgSuccess := fmt.Sprintf("%s %s successfully protected by oauth-proxy-operator", instance.Spec.ResourceKind, instance.Spec.ResourceName)
		if err := r.updateStatus(instance, msgSuccess); err != nil {
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func (r *OAuthProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oauthv1alpha1.OAuthProxy{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&routev1.Route{}).
		Complete(r)
}
