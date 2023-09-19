# KEDA Sample for CloudEvents


A simple Docker container written in **go** that will receive CloudEvent through http.


```
$ curl -X POST 'http://localhost:8899'
```

And one endpoint to get received CloudEvents.

```
$ curl -X GET  http://localhost:8899/getCloudEvent/{eventreason}
```
Sample response:
```
{"specversion":"1.0","id":"6a00b3ae-3503-4e69-8415-b2166512a30a","source":"/cluster-sample/cloudevent-test-ns/keda","type":"/cluster-sample/cloudevent-test-ns/workload/cloudevent-test-so","datacontenttype":"application/json","time":"2023-09-19T08:39:16.500738389Z","data":{"reason":"ScaledObjectCheckFailed","message":"ScaledObject doesn't have correct scaleTargetRef specification"}}
```

Empty response:
```
"Empty"
```

