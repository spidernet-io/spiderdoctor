# Usage

```shell

cat <<EOF > nethttp1.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: testhttp1
spec:
  schedule:
    startAfterMinute: 1
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
  request:
    testIPv4: true
    testIPv6: true
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f nethttp1.yaml


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
    testIPv4: true
    testIPv6: true
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f nethttp2.yaml


cat <<EOF > netdns1.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Netdns
metadata:
  name: testdns1
spec:
  schedule:
    startAfterMinute: 10
    roundNumber: 1
    intervalMinute: 60
    timeoutMinute: 10
  request:
    testIPv4: true
    testIPv6: true
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f netdns1.yaml

```
