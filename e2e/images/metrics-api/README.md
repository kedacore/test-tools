# KEDA Sample for Metrics API


A simple Docker container written in **go** that will expose a metric value. This metrics will be used with for a KEDA scaler [Metrics APi](https://keda.sh/docs/latest/scalers/metrics-api/)

```
$ curl -X GET http:localhost:8080/api/value
```
Sample response:
```
{
    value: 0
    success: true,
}
```

### Get the metric value using basic auth:

To use basic authentication `AUTH_USERNAME` and `AUTH_PASSWORD` have to be passed as environment variables.

```
$ curl --location --request GET 'http://localhost:8080/api/basic/value' --header 'Authorization: Basic <base-64>'
```

### Get the metric value using a bearer token:

To use bearer token `AUTH_TOKEN` have to be passed as environment variables.

```
$ curl --location --request GET 'http://localhost:8080/api/token/value' --header 'Authorization: Bearer <token-here>'
```

## Update the metric

Metric value can be updated using the `api/value` endpoint:

```
$ curl --location --request POST 'http://localhost:8080/api/value/10' 
```
And also using the client:

```
$ docker exec <container-name> /client -value=13
```
