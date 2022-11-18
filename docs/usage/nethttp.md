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
    durationInSecond: 5
    qps: 10
    perRequestTimeoutInSecond: 5
  success:
    successRate: 1
    meanAccessDelayInMs: 10000
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




metric introduction
```shell

        "Metrics": {
          "latencies": {
            "total": 27964545,
            "mean": 2796454,
            "50th": 2821970,
            "90th": 3102803,
            "95th": 3188759,
            "99th": 3188759,
            "max": 3188759,
            "min": 2362429
          },
          "bytes_in": {
            "total": 2357,
            "mean": 235.7
          },
          "bytes_out": {
            "total": 0,
            "mean": 0
          },
          "earliest": "2022-11-18T04:55:20.22108713Z",
          "latest": "2022-11-18T04:55:24.721276724Z",
          "end": "2022-11-18T04:55:24.723858358Z",
          "duration": 4500189594,
          "wait": 2581634,
          "requests": 10,
          "rate": 2.222128599500068,
          "throughput": 2.220854556815161,
          "success": 1,
          "status_codes": {
            "200": 10
          },
          "errors": []
        }

```
