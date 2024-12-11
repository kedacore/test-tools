# KEDA Sample for using BoundServiceAccountToken trigger authentication parameter source with a metrics API scaler

A simple Docker container written in **go** that will expose metrics. Requires that the container is running in a Kubernetes cluster. These metrics will be used for a KEDA scaler [Metrics API](https://keda.sh/docs/latest/scalers/metrics-api/) but requires bearer auth that has permissions to access the GET endpoint. The server delegates auth decisions to a the k8s auth api server.

## GET /api/value

### Without auth

```bash
$ curl http:localhost:8080/api/value
Unauthorized
```

### With auth

This time with a k8s service account token that has the necessary permissions dictated by this server (*in KEDA, this token will get populated by specifying the service account name in the trigger auth parameter using the `BoundServiceAccountToken` trigger auth source*)

```bash
$ curl http://localhost:8080/api/value -H "Authorization: Bearer myk8stoken123"
{
    "value": 0
}
```

## POST /api/value/{value}

And one endpoint to update the metric. This endpoint does not require any authentication/authorization for ease of testing.

```bash
curl --location --request POST 'http://localhost:8080/api/value/10'
```

### Notes

This is the RBAC policy that is required for the service account to access the API:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{.ServiceAccountName}}
rules:
- nonResourceURLs:
  - /api/value
  verbs:
  - get
```

Note that whatever policy enforced is completely arbitrary. The example is only to illustrate that the server delegates auth decisions to the k8s auth api server and thus requires a valid permissive token to access this endpoint.
