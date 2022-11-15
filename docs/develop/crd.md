# develop new crd

1. define CRD in pkg/k8s/apis/spiderdoctor.spidernet.io/v1/xx_types.go
   add role to pkg/k8s/apis/spiderdoctor.spidernet.io/v1/rbac.go

2. make update_openapi_sdk

3. add crd to MutatingWebhookConfiguration and ValidatingWebhookConfiguration in charts/templates/tls.yaml 

4. implement the interface pkg/pluginManager/types in pkg/plugins/xxxx
   register your interface in pkg/pluginManager/types/manager.go
