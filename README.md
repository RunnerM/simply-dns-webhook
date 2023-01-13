
<p align="center">
  <img src="https://user-images.githubusercontent.com/51089137/193522982-c0792104-7ecd-4e6c-a4a9-56b066b65331.png" height="156" width="312" alt="logo" />
</p>


# Simply DNS webhook service for cert-manager support     [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/simply-dns-webhook)](https://artifacthub.io/packages/search?repo=simply-dns-webhook)
This service can be installed side by side with cert manager and can be used to handle dns-01 challeneges provided by cert manager. All documentation on how to configure dns-01 chalanges can be found at [cert-manager.io](https://cert-manager.io/docs/configuration/acme/dns01/webhook/)

### Version support:
The version compatibility I have tested for can be seen below:

| cert-manager version | simply-dns-webhook version |
|----------------------|----------------------------|
| `1.9.0`              | `1.0.3`                    |
| `1.10.1`             | `1.1.0`- `1.1.1`           |



### Deploy
#### Helm chart: 
Add repo:
```shell
    helm repo add simply-dns-webhook https://runnerm.github.io/simply-dns-webhook/
```
Then:
```shell
    helm install my-simply-dns-webhook simply-dns-webhook/simply-dns-webhook --version 1.1.1
```
#### As sub-chart:
```YAML
    dependencies:
        - name: simply-dns-webhook
          version: 1.1.1
          repository: https://runnerm.github.io/simply-dns-webhook/
          alias: simply-dns-webhook
```
### Usage:

**Credentials secret:**
You have to create the secret containing your simply.com api credential on your own, and 
it's name has to match with the secret ref name provided in the config of the cert-manager
issuer/cluster issuer.


#### Issuer/ClusterIssuer:
```YAML
    apiVersion: cert-manager.io/v1
    kind: ClusterIssuer
    metadata:
        name: letsencrypt-nginx
    spec:
        acme:
            email: <your_acme_email>
            server: https://acme-v02.api.letsencrypt.org/directory
            privateKeySecretRef:
                name: letsencrypt-nginx-private-key
            solvers:
            - dns01:
                webhook:
                    groupName: com.github.runnerm.cert-manager-simply-webhook
                    solverName: simply-dns-solver
                    config:
                        secretName: simply-credentials # notice the name
            selector:
                dnsZones:
                - '<your_domain>'
```

**Credentials in config:**
You may choose to use the webhook configuration directly as shown below.
**_(use it at your own risk)_**
```diff
-              secretName: simply-credentials # notice the name
+              accountName: "<account-name>"
+              apiKey: "<api-key>"
```
#### Secret
```YAML
    apiVersion: v1
    kind: Secret
    data:
        account-name: <your_account_name>
        api-key: <your_api_key>
    metadata:
        name: simply-credentials # notice the name
        namespace: <namespace-where-cert-manager-is-installed>
    type: Opaque
```
### cert-manager namespace:

You may override values with your own values if you choose to install cert-manager in custom namespace as follows (this is necessary for proper functioning):
```YAML
    simply-dns-webhook:
        certManager:
            namespace: <cert-manager-namespace>
            serviceAccountName: <cert-manager-namespace>
```
### Resources:
I leave the choice of the resource constraints to you since you know what you run the service on. ;) 
```YAML
    simply-dns-webhook:
        resources: 
            limits:
                cpu: 100m  
                memory: 128Mi
            requests:
                cpu: 100m
                memory: 128Mi
```

## Running the test suite

Update [config](testdata/simply-dns-webhook/config.json) with your API credentials and run:

```bash
$ TEST_ZONE_NAME=example.com. make test
```

## Parameters

The following table lists the configurable parameters of the simply-dns-webhook chart, and their default values.

| Parameter                        | Description                                     | Default                                          |
|----------------------------------|-------------------------------------------------|--------------------------------------------------|
| `groupName`                      | Group name for the webhook                      | `com.github.runnerm.cert-manager-simply-webhook` |
| `certManager.namespace`          | cert-manager namespace                          | `cert-manager`                                   |
| `certManager.serviceAccountName` | cert-manager service account name               | `cert-manager`                                   |
| `image.repository`               | Docker image repository                         | `deyaeddin/cert-manager-webhook-hetzner`         |
| `image.tag`                      | Docker image tag                                | `v1.1.1`                                         |
| `image.pullPolicy`               | Docker image pull policy                        | `Always`                                         |
| `nameOverride`                   | Name override for the chart                     | `""`                                             |
| `fullnameOverride`               | Full name override for the chart                | `""`                                             |
| `service.type`                   | Service type                                    | `ClusterIP`                                      |
| `service.port`                   | Service port                                    | `443`                                            |
| `resources`                      | Pod resources                                   | Check `values.yaml` file                         |
| `nodeSelector`                   | Node selector                                   | `nil`                                            |
| `tolerations`                    | Node toleration                                 | `nil`                                            |
| `affinity`                       | Node affinity                                   | `nil`                                            |

##### Special credits to: **Keyhole Aps**
