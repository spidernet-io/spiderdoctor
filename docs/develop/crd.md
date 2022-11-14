# develop new crd

1. define CRD in pkg/k8s/apis/spiderdoctor.spidernet.io/v1/xx_types.go

2. make update_openapi_sdk

3. add crd to MutatingWebhookConfiguration and ValidatingWebhookConfiguration in charts/templates/tls.yaml 

4. implement it in pkg/xxManager
