# Usage

```shell

cat <<EOF > nethttp1.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: testhttp1
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 5000
EOF
kubectl apply -f nethttp1.yaml

```


```shell

cat <<EOF > nethttp2.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: testhttp2
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
    sourceAgentNodeSelector:
        matchExpressions:
          - { key: "kubernetes.io/hostname", operator: In, values: ["spiderdoctor-worker"] }
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF
kubectl apply -f nethttp2.yaml

```


```shell

cat <<EOF > nethttp3.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: testhttp3
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
  target:
    targetUrl: "http://172.19.0.6"
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 5000
EOF
kubectl apply -f nethttp3.yaml

```


