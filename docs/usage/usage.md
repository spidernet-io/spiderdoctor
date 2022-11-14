# Usage

```shell

cat <<EOF > nethttp.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: test1
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

kubectl apply -f nethttp.yaml

```
