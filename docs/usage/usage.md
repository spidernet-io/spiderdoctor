# Usage

```shell

cat <<EOF > example.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: test1
spec:
  schedule:
    startAfterMinute: 10
    roundNumber: 1
    intervalMinute: 60
    TimeoutMinute: 10
  request:
    testIPv4: true
    testIPv6: true
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInSecond: 5
  failureCondition:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f example.yaml

```
