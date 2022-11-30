# Nethttp

## concept 

Fo this kind task, each spiderdoctor agent will send http request to specified target, and get success rate and mean delay. 
It could specify success condition to tell the result succeed or fail. 
And, more detailed report will print to spiderdoctor agent stdout, or save to disc by spiderdoctor controller.

the following is the spec of nethttp
```shell
apiVersion: v1
items:
- apiVersion: spiderdoctor.spidernet.io/v1
  kind: Nethttp
  metadata:
    generation: 1
    name: testhttp1
  spec:
    schedule:
      intervalMinute: 2
      roundNumber: 2
      startAfterMinute: 0
      timeoutMinute: 1
      sourceAgentNodeSelector:
        matchExpressions:
          - { key: "kubernetes.io/hostname", operator: In, values: ["spiderdoctor-worker"] }
    request:
      durationInSecond: 5
      perRequestTimeoutInMS: 1000
      qps: 10
    target:
      targetUser:
        targetUser:
          method: GET
          url: http://172.80.1.2
      targetAgent:
        testClusterIp: true
        testEndpoint: true
        testIPv4: true
        testIPv6: true
        testIngress: false
        testMultusInterface: false
        testNodePort: true
    success:
      meanAccessDelayInMs: 10000
      successRate: 1
  status:
    doneRound: 1
    expectedRound: 2
    finish: false
    lastRoundStatus: succeed
    history:
    - deadLineTimeStamp: "2022-11-21T06:19:53Z"
      duration: 10.117152806s
      startTimeStamp: "2022-11-21T06:18:53Z"
      endTimeStamp: "2022-11-21T06:19:03Z"
      failedAgentNodeList: []
      notReportAgentNodeList: []
      roundNumber: 2
      status: notstarted
      succeedAgentNodeList: []
```

* spec.schedule: set how to schedule the task.

      roundNumber: how many rounds it should be to run this task

      intervalMinute:  the time interval in minute, for run each round for this task

      startAfterMinute: when the start the first round

      timeoutMinute: the timeout in minute for each round, when the rask does not finish in time, it results to be failuire

      sourceAgentNodeSelector [optional]: set the node label selector, then, the spiderdoctor agent who locates on these nodes will implement the task. If not set this field, all spiderdoctor agent will execute the task

* spec.request: how each spiderdoctor agent should send the http request

    durationInSecond: for each round, the duration in second how long the http request lasts

    perRequestTimeoutInMS: timeout in ms for each http request 

    qps: qps

* spec.target: set the target of http request. it could not set targetUser and targetAgent at the same time

      targetUser [optional]: set an user-defined URL for the http request

        url: the url for http

        method: http method, must be one of GET POST PUT DELETE CONNECT OPTIONS PATCH HEAD

      targetAgent: [optional]: set the http tareget to spiderdoctor agents

        testClusterIp: send http request to the cluster ipv4 or ipv6 address of spiderdoctor agnent, according to testIPv4 and testIPv6.

        testEndpoint: send http request to other spiderdoctor agnent ipv4 or ipv6 address according to testIPv4 and testIPv6.

        testMultusInterface: whether send http request to all interfaces ip in testEndpoint case.

        testIPv4: test any IPv4 address. Notice, the 'enableIPv4' in configmap  spiderdocter must be enabled

        testIPv6: test any IPv6 address. Notice, the 'enableIPv6' in configmap  spiderdocter must be enabled

        testIngress: send http request to the ingress ipv4 or ipv6 address of spiderdoctor agnent

        testNodePort: send http request to the nodePort ipv4 or ipv6 address with each local node of spiderdoctor agnent , according to testIPv4 and testIPv6.

        >notice: when test targetAgent case, it will send http request to all targets at the same time with spec.request.qps for each one. That meaning, the actually QPS may be bigger than spec.request.qps

* spec.success: define the success condition of the task result 

    meanAccessDelayInMs: mean access delay in MS, if the actual delay is bigger than this, it results to be failure

    successRate: the success rate of all http requests. Notice, when a http response code is >=200 and < 400, it's treated as success. if the actual whole success rate is smaller than successRate, the task results to be failure

* status: the status of the task
    doneRound: how many rounds have finished

    expectedRound: how many rounds the task expect

    finish: whether all rounds of this task have finished

    lastRoundStatus: the result of last round

    history:
        roundNumber: the round number

        status: the status of this round

        startTimeStamp: when this round begins

        endTimeStamp: when this round finally finished

        duration: how long the round spent

        deadLineTimeStamp: the time deadline of a round 

        failedAgentNodeList: the node list where failed spiderdoctor agent locate

        notReportAgentNodeList: the node list where uknown spiderdoctor agent locate. This means these agents have problems.

        succeedAgentNodeList: the node list where successful spiderdoctor agent locate


## example 

a quick task to test spiderdoctor agent, to verify the whole network is ok, each agent could reach each other

```shell

cat <<EOF > nethttp-test-agent.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: test-agent
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
  request:
    durationInSecond: 2
    qps: 2
    perRequestTimeoutInMS: 1000
  success:
    successRate: 1
    meanAccessDelayInMs: 10000
EOF
kubectl apply -f nethttp-test-agent.yaml

```


a detail task to test spiderdoctor agent

```shell
cat <<EOF > test-detail-agent.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: test-detail-agent
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
  target:
    targetAgent:
      testIPv4: true
      testIPv6: true
      testEndpoint: true
      testMultusInterface: true
      testClusterIp: true
      testNodePort: true
      testLoadBalancer: true
      testIngress: true
  request:
    durationInSecond: 2
    qps: 2
    perRequestTimeoutInMS: 1000
  success:
    successRate: 1
    meanAccessDelayInMs: 10000
EOF
kubectl apply -f test-detail-agent.yaml

```

test custom URL
```shell

cat <<EOF > test-custom.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: test-custom
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
  target:
    targetUser:
      url: "http://172.80.1.2"
      method: "GET"
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInMS: 1000
  success:
    successRate: 1
    meanAccessDelayInMs: 5000
EOF
kubectl apply -f test-custom.yaml

```

use the spicified spiderdoctor agents to send the http request
```shell

cat <<EOF > source-agent.yaml
apiVersion: spiderdoctor.spidernet.io/v1
kind: Nethttp
metadata:
  name: source-agent
spec:
  schedule:
    startAfterMinute: 0
    roundNumber: 2
    intervalMinute: 2
    timeoutMinute: 1
    sourceAgentNodeSelector:
        matchExpressions:
          - { key: "kubernetes.io/hostname", operator: In, values: ["spiderdoctor-worker"] }
  target:
    targetUser:
      url: "http://172.80.1.2"
      method: "GET"
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInMS: 1000
  success:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF
kubectl apply -f source-agent.yaml

```





## debug

when something wrong happen, see the log for your task with following command
```shell
#get log 
CRD_KIND="nethttp"
CRD_NAME="test1"
kubectl logs -n kube-system  spiderdoctor-agent-v4vzx | grep -i "${CRD_KIND}.${CRD_NAME}"

```


## report

when the spiderdoctor is not enabled to aggerate reports, all reports will be printed in the stdout of spiderdoctor agent.
Use the following command to get its report
```shell
kubectl logs -n kube-system  spiderdoctor-agent-v4vzx | jq 'select( .TaskName=="nethttp.testhttp1" )'
```

when the spiderdoctor is enabled to aggregate reports, all reports will be collected in the PVC or hostPath of spiderdoctor controller.


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
          "wait": 2581634, # Wait is the extra time waiting for responses from targets.
          "requests": 10, #the total number of requests executed
          "rate": 2.222128599500068, #Rate is the rate of sent requests per second.
          "throughput": 2.220854556815161, #Throughput is the rate of successful requests per second.
          "success": 1, #percentage of non-error responses
          "status_codes": {
            "200": 10
          },
          "errors": []
        }
```
