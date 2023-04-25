# netdns

```shell

cat <<EOF > netdns.yaml
apiVersion: spiderdoctor.spidernet.io/v1beta1
kind: Netdns
metadata:
  name: testdns
spec:
  schedule:
    schedule: "1 1"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    targetDns:
      testIPv4: true
      testIPv6: false
    protocol: udp
  request:
    durationInSecond: 10
    qps: 20
    perRequestTimeoutInMS: 500
    domain: "kube-dns.kube-system.svc.cluster.local"
  success:
    successRate: 1
    meanAccessDelayInMs: 10000
EOF

kubectl apply -f netdns.yaml

```

```shell

cat <<EOF > netdns1.yaml
apiVersion: spiderdoctor.spidernet.io/v1beta1
kind: Netdns
metadata:
  name: testdns1
spec:
  schedule:
    schedule: "1 2"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    protocol: udp
    targetUser:
      server: 172.18.0.1
      port: 53
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInMS: 500
    domain: "baidu.com"
  success:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f netdns1.yaml

```