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
	oauthv1alpha1 "oauth-proxy-operator/api/v1alpha1"
	"time"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OAuthProxyReconciler) ensureDeployment(request ctrl.Request, cr *oauthv1alpha1.OAuthProxy, d *appsv1.Deployment) (*ctrl.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      d.Name,
		Namespace: cr.Spec.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		r.Log.Info("Creating a new Deployment", "Deployment.Namespace", d.Namespace, "Deployment.Name", d.Name)
		err = r.Client.Create(context.TODO(), d)

		if err != nil {
			// Deployment failed
			r.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", d.Namespace, "Deployment.Name", d.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		r.Log.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *OAuthProxyReconciler) ensureServiceAccount(request ctrl.Request, instance *oauthv1alpha1.OAuthProxy, s *corev1.ServiceAccount) (*ctrl.Result, error) {
	found := &corev1.ServiceAccount{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Spec.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the service account
		r.Log.Info("Creating a new Service Account", "ServiceAccount.Namespace", s.Namespace, "ServiceAccount.Name", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Service Account", "ServiceAccount.Namespace", s.Namespace, "ServiceAccount.Name", s.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the service account not existing
		r.Log.Error(err, "Failed to get Service Account")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *OAuthProxyReconciler) ensureService(request ctrl.Request, instance *oauthv1alpha1.OAuthProxy, s *corev1.Service) (*ctrl.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Spec.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		r.Log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &ctrl.Result{}, err
		}
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the service not existing
		r.Log.Error(err, "Failed to get Service")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *OAuthProxyReconciler) ensureRoute(request ctrl.Request, instance *oauthv1alpha1.OAuthProxy, ro *routev1.Route) (*ctrl.Result, error) {
	found := &routev1.Route{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: ro.Name, Namespace: instance.Spec.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the route
		r.Log.Info("Creating a new Route", "Route.Namespace", ro.Namespace, "Route.Name", ro.Name)
		err = r.Client.Create(context.TODO(), ro)

		time.Sleep(350 * time.Millisecond)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Route", "Route.Namespace", ro.Namespace, "Route.Name", ro.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the route not existing
		r.Log.Error(err, "Failed to get Route")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
