# install

## production env


## POC or E2E environment

when POC or E2E case, it could disable the controller to collect reports, so no need to install strogeClass.

the following method leads the agent just print report to console
```shell 
helm repo add spiderdoctor https://spidernet-io.github.io/spiderdoctor

helm install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=false
```

the following method leads controller collects all report to disc of local host. BTW, when the spiderdoctor controller is schedules to other nodes, the historical reports will be not migrated 
```shell 
helm repo add spiderdoctor https://spidernet-io.github.io/spiderdoctor

helm  install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=true \
    --set feature.aggregateReport.controller.reportHostPath="/var/run/spiderdoctor/controller"
```

## production environment

the following method leads the spiderdoctor controller collect report to stroage, so firstly, it should install storageClass

```shell 
helm repo add spiderdoctor https://spidernet-io.github.io/spiderdoctor

helm  install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=true \
    --set feature.aggregateReport.controller.pvc.enabled=true \
    --set feature.aggregateReport.controller.pvc.storageClass=local \
    --set feature.aggregateReport.controller.pvc.storageRequests="100Mi" \
    --set feature.aggregateReport.controller.pvc.storageLimits="500Mi"
```

## multus environment

if it is required to test all interface of agent pod, it should annotate the agent with multus annotation

```shell 
helm repo add spiderdoctor https://spidernet-io.github.io/spiderdoctor

# replace following with actual multus configuration
MULTUS_DEFAULT_CNI=kube-system/k8s-pod-network
MULTUS_ADDITIONAL_CNI=kube-system/macvlan

helm install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=false \
    --set spiderdoctorAgent.podAnnotations.v1\.multus-cni\.io/default-network=${MULTUS_DEFAULT_CNI} \
    --set spiderdoctorAgent.podAnnotations.k8s\.v1\.cni\.cncf\.io/networks=${MULTUS_ADDITIONAL_CNI}
    
```
