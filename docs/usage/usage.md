# Usage

```shell

cat <<EOF > example.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: test1
spec:
  schedule:
    roundNumber: 1
    interval: 60
  enabledIPv4: true
  enabledIPv6: true
EOF

kubectl apply -f example.yaml

```
