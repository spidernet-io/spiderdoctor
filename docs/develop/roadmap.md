# feature 

测试集群功能、性能巡检工具，批量化巡检，极大降低人工负载

## CRD 

### Nethttp

* 大规模集群部署后，巡检每个节点上的网络情况，具体是：本节点是否能够通过 pod ip（multus多网卡）、cluster ip、
    nodePort、loadbalancer ip、ingress ip 等所有网络渠道，访问到集群其它节点

* 集群所有的 node 上 去 压测一个集群内/外的应用地址，以查看应用的性能、集群每个角落到达该应用的连通性、给应用注入压力复现某类bug

* 给 api server 注入压力，以辅助排查 其他组件（依赖 api server）的高可用

* 生产和开发环境的心跳巡检，以qps=1为压力，每 1m 间隔巡检整个集群内 full mesh 网络的连通性、其它节点到一个应用的可用性

### NetDetecthttp

* 自动发现应用的最大性能

### Netdns

* 大规模集群部署后，测试集群中每个角落访问 dns 的连通性

* 大规模集群部署后，调试 coredns 的副本数，确认是否满足设计需求

* 测试集群外部的 DNS 服务

### Netdetectdns

* 自动发现 dns 的最大性能

### Nettcp

### Netudp

### StorageLocalDisk

* 大规模集群部署后，测试集群中每个主机上的磁盘 吞吐量 和 延时

### CpuPressure

* 给每个主机上注入 CPU 压力，以测试应用的稳定性，复现一些 bug

### MemoryPressure

* 给每个主机上注入 memory 压力，以测试应用的稳定性，复现一些 bug

### RegistryHealthy

* 检测每个节点到镜像仓库的连通性


## report

支持通过 API 获取报告

支持 pvc、本地磁盘存储

日志吐出

## metric

## 其它

如果有 job 时间重叠了，则只允许运行一个 或者 多个，避免自身 CPU 不足影响 job 的结果

中间件、etcd 等 探测


