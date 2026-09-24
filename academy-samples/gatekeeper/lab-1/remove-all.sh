#!/bin/sh -x
kubectl delete k8spspprivilegedcontainer.constraints.gatekeeper.sh/disallow-privileged-container
kubectl delete k8srequiredresources.constraints.gatekeeper.sh/require-container-limits-and-requests
kubectl delete ConstraintTemplate/k8spspprivilegedcontainer
kubectl delete ConstraintTemplate/k8srequiredresources
kubectl delete Pod/frontend
kubectl delete Pod/redis
kubectl delete Deployment/nginx-deployment
kubectl delete -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml