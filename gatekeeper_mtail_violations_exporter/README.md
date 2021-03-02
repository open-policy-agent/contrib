# Gatekeeper Mtail Violations Exporter
The configurations contained here will install Gatekeeper, adding a sidecar which reads Gatekeeper's audit loop logs and exposes violations on a prometheus scrape point.


## Setup
1. Run `kustomize build . | kubectl apply -f -` to deploy

## Container Runtimes
The default configuration for these manifests is for Docker. If you're using another runtime you will need to updated the log locations 

## Example violations

```
opa_violations{constraint_kind="MyConstraint",
               constraint_name="my-constraint",
               context="pod:violating-pod",
               msgid="deployment-manifest-violates-constraint",
               resource_kind="Deployment",
               resource_name="MyDeployment",
               resource_namespace="my-namespace"}
```

## Storage Considerations
The configs provided here include prometheus scrape annotations. If you'd prefer to use a [Service Monitor](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md#related-resources) you will need to update the configurations here. If you are using another TSDB you may need to adjust these configurations according to the DB's specifications.

### Performance
Usage of this operator assumes your TSDB is able to handle as many timeseries as you have policies. 