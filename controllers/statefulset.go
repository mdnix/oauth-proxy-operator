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

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *OAuthProxyReconciler) handleStatefulSetChanges(cr *oauthv1alpha1.OAuthProxy) (*ctrl.Result, error) {
	found := &appsv1.StatefulSet{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cr.Spec.ResourceName,
		Namespace: cr.Spec.Namespace,
	}, found)
	if err != nil {
		// The Statefulset may not have been created yet, so requeue
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	nContainers := len(found.Spec.Template.Spec.Containers)
	for i := 0; i < nContainers; i++ {
		container := found.Spec.Template.Spec.Containers[i].Name
		if container == "oauth-proxy" {
			// proxy side car already exists
			return &ctrl.Result{}, nil
		}
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.Client, found, func() error {
		sidecar, err := r.OAuthProxySideCar(cr)
		if err != nil {
			return err
		}
		found.Spec.Template.Spec.Containers = append(found.Spec.Template.Spec.Containers, sidecar)
		found.Spec.Template.Spec.ServiceAccountName = cr.Spec.ResourceName
		vol := corev1.Volume{
			Name: cr.Spec.ResourceName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: cr.Spec.ResourceName,
				},
			},
		}
		found.Spec.Template.Spec.Volumes = append(found.Spec.Template.Spec.Volumes, vol)
		return nil
	})
	if err != nil {
		delay := 5 * time.Second
		return &ctrl.Result{RequeueAfter: delay}, errors.Wrap(err, "Unable to CreateOrUpdate StatefulSet")
	}

	return &ctrl.Result{}, nil

}

func (r *OAuthProxyReconciler) finalizeStatefulSet(cr *oauthv1alpha1.OAuthProxy) (*ctrl.Result, error) {
	found := &appsv1.StatefulSet{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cr.Spec.ResourceName,
		Namespace: cr.Spec.Namespace,
	}, found)
	if err != nil {
		// The Statefulset may not have been created yet, so requeue
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	nContainers := len(found.Spec.Template.Spec.Containers)
	if nContainers > 1 {
		for i := 0; i < nContainers; i++ {
			container := found.Spec.Template.Spec.Containers[i].Name
			if container == "oauth-proxy" {
				found.Spec.Template.Spec.Containers = removeContainer(found.Spec.Template.Spec.Containers, i)
				found.Spec.Template.Spec.ServiceAccountName = "default"
			}
		}
	}
	nVolumes := len(found.Spec.Template.Spec.Volumes)
	if nVolumes > 0 {
		for i := 0; i < nVolumes; i++ {
			vol := found.Spec.Template.Spec.Volumes[i].Name
			if vol == cr.Spec.ResourceName {
				found.Spec.Template.Spec.Volumes = removeVolume(found.Spec.Template.Spec.Volumes, i)
			}
		}

	}

	if err := r.Client.Update(context.Background(), found); err != nil {
		return &ctrl.Result{}, err
	}

	// Remove OAuthProxy finalizer, release for deletion of CR
	cr.ObjectMeta.Finalizers = removeString(cr.ObjectMeta.Finalizers, finalizer)
	if err := r.Client.Update(context.Background(), cr); err != nil {
		return &ctrl.Result{}, err
	}

	return &ctrl.Result{}, nil
}
