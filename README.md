
<p align="center">
  <img src="https://user-images.githubusercontent.com/51089137/192287277-b5682293-62fe-4cc6-9d67-574aa3a95390.png" height="156" width="312" alt="logo" />
</p>


# Simply DNS webhook service for cert-manager support     [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/simply-dns-webhook)](https://artifacthub.io/packages/search?repo=simply-dns-webhook)
This service can be installed side by side with cert manager and can be used to handle dns-01 challeneges provided by cert manager. All documentation on how to configure dns-01 chalanges can be found at [cert-manager.io](https://cert-manager.io/docs/configuration/acme/dns01/webhook/)

### Deploy
#### Helm chart: 
Add repo:

    helm repo add simply-dns-webhook https://runnerm.github.io/simply-dns-webhook/
Then:

    helm install my-simply-dns-webhook simply-dns-webhook/simply-dns-webhook --version 1.0.3

#### As sub-chart:
    dependencies:
        - name: simply-dns-webhook
          version: 1.0.3
          repository: https://runnerm.github.io/simply-dns-webhook/
          alias: simply-dns-webhook

### Usage:

**Credentials secret:**
You have to create the secret containing your simply.com api credential on your own, and 
it's name has to match with the secret ref name provided in the config of the cert-manager
issuer/cluster issuer.


#### Issuer/ClusterIssuer:
    apiVersion: cert-manager.io/v1
    kind: ClusterIssuer
    metadata:
        name: letsencrypt-nginx
    spec:
        acme:
            email: mks@usekeyhole.com
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

#### Secret

    apiVersion: v1
    kind: Secret
    data:
        account-name: <your_account_name>
        api-key: <your_api_key>
    metadata:
        name: simply-credentials # notice the name
        namespace: kh-networking
    type: Opaque

### cert-manager namespace:

You may override values with your own values if you choose to install cert-manager in custom namespace as follows (this is necessary for proper functioning):

    simply-dns-webhook:
        certManager:
            namespace: <cert-manager-namespace>
            serviceAccountName: <cert-manager-namespace>

### Resources:
I leave the choice of the resource constraints to you since you know what you run the service on. ;) 

    simply-dns-webhook:
        resources: 
            limits:
                cpu: 100m  
                memory: 128Mi
            requests:
                cpu: 100m
                memory: 128Mi

##### Special credits to: **Keyhole Aps**
