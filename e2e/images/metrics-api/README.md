# KEDA Sample for Metrics API


A simple Docker container written in **go** that will expose metrics. This metrics will be used with for a KEDA scaler [Metrics APi](https://keda.sh/docs/latest/scalers/metrics-api/)

```
$ curl http:localhost:8080/api/value
```
Sample response:
```
{
    value: 0
    success: true,
}
```

And one endpoint to update the metric.

```
$ curl --location --request POST 'http://localhost:8080/api/value/10'
```
