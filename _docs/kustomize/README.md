Kustomize Kubernetes
--------------------

Directory contains templates and files to configure a Kubernetes cluster to support Akash Provider services.

# Akash Network labeling

Key and value label pairs used by Akash services in the Kubernetes Provider implementation.

* `akash.network` key indicates a resource is apart of tooling or speicification and value is set to `true`.
* `akash.network/name` general key for Akash resources.
* `akash.network/component` key for general name of a component, eg: `akash-provider`
* `akash.network/tenant-namespace` key informs the (K8s) Namespace of a resource. Namespaces are generated new for tenants' workloads. eg: `ab6ij0o2b3vtbh1t26j3up9id2154ihgdk2apln182t6e`
* `akash.network/manifest-service` key is derived from the SDL as the application or group's name. Example value being `web` in demo-app examples.


## `networking/`

Normal Kubernetes declaration files to be `kubectl apply -f networking/`

* `akash-services` namespace declaration
* Default-deny Network policies for the `default` namespace to cripple the potential of malicious containers being run there.

## `akash-services/`

[Kustomize](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/) directory for configuring Network Policies to support `akash-services` namespace of Akash apps.
`kubectl kustomize akash-services/ | kubectl apply -f -`

## `akashd/` 

Kustomize directory to configure running the `akash` blockchain node service.

`kubectl kustomize akashd/ | kubectl apply -f -`

## `akash-provider/`

Kustomize directory to configure running the `akash-provider` instance which manages tenant workloads within the cluster.

`kubectl kustomize akash-provider/ | kubectl apply -f -`
