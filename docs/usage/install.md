# install

## production env


## POC or E2E env


helm add spiderdoctor

helm  install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=false


helm  install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=true \
    --set feature.aggregateReport.controller.reportHostPath="/var/run/spiderdoctor/controller"


#===================

helm  install spiderdoctor spiderdoctor/spiderdoctor \
    -n kube-system --wait --debug \
    --set feature.enableIPv4=true --set feature.enableIPv6=true \
    --set feature.aggregateReport.enabled=true \
    --set feature.aggregateReport.controller.pvc.enabled=true \
    --set feature.aggregateReport.controller.pvc.storageClass=local \
    --set feature.aggregateReport.controller.pvc.storageRequests="100Mi" \
    --set feature.aggregateReport.controller.pvc.storageLimits="500Mi"


#===================


HELM_OPTION+=" --set spiderdoctorAgent.podAnnotations.v1\.multus-cni\.io/default-network=kube-system/k8s-pod-network " ; \
HELM_OPTION+=" --set spiderdoctorAgent.podAnnotations.k8s\.v1\.cni\.cncf\.io/networks=kube-system/ptp " ; \

