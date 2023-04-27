
# environment
- Kubenetes: `v1.25.4`
- container runtime: `containerd 1.6.12`
- OS: `CentOS Linux 8`
- kernel: `4.18.0-348.7.1.el8_5.x86_64`

| Node     | Role          | CPU  | Memory |
| -------- | ------------- | ---- | ------ |
| master1  | control-plane | 4C   | 8Gi    |
| master2  | control-plane | 4C   | 8Gi    |
| master3  | control-plane | 4C   | 8Gi    |
| worker4  |               | 3C   | 8Gi    |
| worker5  |               | 3C   | 8Gi    |
| worker6  |               | 3C   | 8Gi    |
| worker7  |               | 3C   | 8Gi    |
| worker8  |               | 3C   | 8Gi    |
| worker9  |               | 3C   | 8Gi    |
| worker10 |               | 3C   | 8Gi    |

# Nethttp

In a pod with a CPU of 1

| client       | time | requests | qps     | Memory |
|--------------|------|----------|---------|--------|
| spiderdoctor | 1m   | 67346    | 1122.43 | 2Gi    |
| wrk          | 1m   | 53612    | 892.85  | 2Mb    |

| client       | time | requests | qps      | Memory |
|--------------|------|----------|----------|--------|
| spiderdoctor | 5m   | 272403   | 908.01   | 2.6Gi  |
| wrk          | 5m   | 265551   | 884.92   | 5Mb    |

| client       | time | requests  | qps     | Memory |
|--------------|------|-----------|---------|--------|
| spiderdoctor | 10m  | 467254    | 778.75  | 4.5Gi  |
| wrk          | 10m  | 565638    | 942.58  | 5Mb    |


# Netdns

In a pod with a CPU of 1

| client       | time | requests | qps       | Memory |
|--------------|------|----------|-----------|--------|
| spiderdoctor | 1m   | 1855511  | 30,925.18 | 23Mb   |
| dnsperf      | 1m   | 1728086  |28800.406  | 8Mb    |

| client       | time | requests | qps       | Memory |
|--------------|------|----------|-----------|--------|
| spiderdoctor | 5m   | 9171699  | 30,572.33 | 100Mb  |
| dnsperf      | 5m   | 8811137  | 29370.34  | 8Mb    |

| client       | time | requests  | qps       | Memory |
|--------------|------|-----------|-----------|--------|
| spiderdoctor | 10m  | 18561282  | 30,935.47 | 173Mb  |
| dnsperf      | 10m  | 17260779  | 28767.666 | 8Mb    |
