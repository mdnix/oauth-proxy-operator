# oauth-proxy-operator

This cluster-scoped operator will protect your OpenShift Route with an oauth-proxy. It uses the integrated OpenShift Identity Providers for ease of use. Might add more OAuth providers in the future and add support to Kubernetes Ingress

The oauth-proxy is configured to use SAR (Subject Access Review). It allows access if the user can view the service of the protected route in the namespace the CustomResource was deployed.

The following kinds can be secured:

Deployment
DeploymentConfig
DaemonSet
StatefulSet

## Install

### CRD

```bash
make install
```

## Run

### Local
```bash
make run ENABLE_WEBHOOKS=false
```

### Inside Cluster

#### Build

```bash
export USERNAME=emdinix
make docker-build IMG=$USERNAME/oauth-proxy-operator:v0.0.1
```

#### Run

```bash
make deploy IMG=$USERNAME/oauth-proxy-operator:v0.0.1
```

## Use

### Create CR

```bash
cat <<EOF | kubectl apply -f -
apiVersion: oauth.emdinix.io/v1alpha1
kind: OAuthProxy
metadata:
  name: protected-deployment
spec:
  namespace: "oauth-proxy"
  resourcename: "frontend"
  resourcekind: "Deployment"
  upstreamport: 8080
EOF
```