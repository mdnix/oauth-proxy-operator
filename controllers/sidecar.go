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
	"fmt"
	oauthv1alpha1 "oauth-proxy-operator/api/v1alpha1"
	"strconv"

	uuid "github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
)

func (r *OAuthProxyReconciler) OAuthProxySideCar(cr *oauthv1alpha1.OAuthProxy) (corev1.Container, error) {

	secret, err := uuid.NewRandom()
	if err != nil {
		return corev1.Container{}, err
	}

	sar := fmt.Sprintf(`{"namespace":"%s","resource":"services","name":"%s","verb":"get"}`, cr.Spec.Namespace, cr.Spec.ResourceName)

	s := corev1.Container{
		Name:            "oauth-proxy",
		Image:           "openshift/oauth-proxy:latest",
		ImagePullPolicy: "IfNotPresent",
		Ports: []corev1.ContainerPort{
			{
				Name:          "public",
				ContainerPort: 8443,
			},
		},
		Args: []string{
			"--https-address=:8443",
			"--provider=openshift",
			"--openshift-service-account=" + cr.Spec.ResourceName,
			"--upstream=http://localhost:" + strconv.Itoa(int(cr.Spec.UpstreamPort)),
			"--tls-cert=/etc/tls/private/tls.crt",
			"--tls-key=/etc/tls/private/tls.key",
			"--openshift-sar=" + sar,
			"--cookie-secret=" + secret.String(),
			"--cookie-expire=1h0m0s",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      cr.Spec.ResourceName,
				MountPath: "/etc/tls/private",
			},
		},
	}

	return s, nil

}
