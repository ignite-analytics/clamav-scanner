# clamav-scanner

ClamAV Scanner Helm Chart

![Version: 0.1.4](https://img.shields.io/badge/Version-0.1.4-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.11](https://img.shields.io/badge/AppVersion-0.5.11-informational?style=flat-square)

A Helm chart for deploying ClamAV Scanner to GKE cluster.

## Prerequisites

Before instaling the chart, you need to have the following:

- A running Kubernetes cluster with at least one node pool configured with Workload Identity Federation for GKE enabled.
Create new or updating existing cluster and node pool. See docs [here](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity).
- IAM service account with at least `roles/storage.objectUser` role assigned to it for `mirrorBucket`, `quarantineBucket` and all buckets that you want to scan with ClamAV.
`roles/storage.objectUser` makes sure that the service account can upload files to it's own buckets and move files between scan bucket and quarantine without granting more powerful roles like `objectAdmin`.

```console
gsutil iam ch serviceAccount:$SA_EMAIL:objectUser gs://$BUCKET_NAME
```

- IAM policy that gives Kubernetes service account the permission to impersonate the service account with `roles/iam.workloadIdentityUser` role.

```console
gcloud iam service-accounts add-iam-policy-binding $SA@$PROJECT_ID.iam.gserviceaccount.com \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:$PROJECT_ID.svc.id.goog[NAMESPACE/KSA_NAME]"
```

- IAM policy that gives the service account Pub/Sub Subscriber role to pull messages from `event-forwarder` pod.

```console
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SA_EMAIL" \
  --role="roles/pubsub.subscriber" \
  --condition=None
```

- IAM policy that gives the service account Monitoring Metric Writer role to export metrics from `event-forwarder` pod to Stackdriver.

```console
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SA_EMAIL" \
  --role="roles/monitoring.metricWriter" \
  --condition=None
```

- Route Eventarc events to GKE. See docs [here](https://cloud.google.com/eventarc/docs/gke/route-trigger-cloud-storage#gcloud)
- Eventarc trigger that listens for new objects in the bucket and triggers the GKE service:

```console
gcloud eventarc triggers create $TRIGGER_NAME \
  --location=$TRIGGER_LOCATION \
  --destination-gke-cluster=$GKE_CLUSTER \
  --destination-gke-location=$GKE_LOCATION \
  --destination-gke-namespace=$GKE_NAMESPACE \
  --destination-gke-service=$KSA_NAME \
  --destination-gke-path=/scan \
  --event-filters="type=google.cloud.storage.object.v1.finalized" \
  --event-filters="bucket=$BUCKET_NAME" \
  --service-account=$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com
```

- Add annotation to the GKE service account to allow impersonation in `values.yaml`:

```yaml
serviceAccount:
  annotations:
    iam.gke.io/gcp-service-account: $SA_NAME@$PROJECT_ID.iam.gserviceaccount.com
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| autoscaling.enabled | bool | `true` |  |
| autoscaling.maxReplicas | int | `5` |  |
| autoscaling.minReplicas | int | `1` |  |
| autoscaling.targetMemoryUtilizationPercentage | int | `80` |  |
| extraConfigMap | object | `{}` | Key/value pairs to be exposed as environment variables |
| fullnameOverride | string | `""` |  |
| image.digest | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/ignite-analytics/clamav-scanner"` |  |
| imagePullSecrets | list | `[]` |  |
| livenessProbe.failureThreshold | int | `1` |  |
| livenessProbe.initialDelaySeconds | int | `0` |  |
| livenessProbe.periodSeconds | int | `10` |  |
| livenessProbe.successThreshold | int | `1` |  |
| livenessProbe.timeoutSeconds | int | `1` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext.fsGroup | int | `102` |  |
| podSecurityContext.fsGroupChangePolicy | string | `"OnRootMismatch"` |  |
| podSecurityContext.runAsGroup | int | `102` |  |
| podSecurityContext.runAsUser | int | `100` |  |
| podSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| poddisruptionbudget.enabled | bool | `true` |  |
| poddisruptionbudget.minAvailable | int | `1` |  |
| readinessProbe.failureThreshold | int | `1` |  |
| readinessProbe.initialDelaySeconds | int | `0` |  |
| readinessProbe.periodSeconds | int | `20` |  |
| readinessProbe.successThreshold | int | `1` |  |
| readinessProbe.timeoutSeconds | int | `1` |  |
| replicaCount | int | `1` |  |
| resources.limits.cpu | string | `"1"` |  |
| resources.limits.memory | string | `"6Gi"` |  |
| resources.requests.cpu | string | `"1"` |  |
| resources.requests.memory | string | `"4Gi"` |  |
| securityContext.allowPrivilegeEscalation | bool | `false` |  |
| securityContext.capabilities.drop[0] | string | `"ALL"` |  |
| securityContext.privileged | bool | `false` |  |
| securityContext.readOnlyRootFilesystem | bool | `false` |  |
| securityContext.runAsNonRoot | bool | `true` |  |
| securityContext.runAsUser | int | `100` |  |
| service.port | int | `80` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `"clamav-scanner"` |  |
| startupProbe.failureThreshold | int | `30` |  |
| startupProbe.initialDelaySeconds | int | `10` |  |
| startupProbe.periodSeconds | int | `10` |  |
| startupProbe.successThreshold | int | `1` |  |
| startupProbe.timeoutSeconds | int | `10` |  |
| tolerations | list | `[]` |  |
| update.affinity | object | `{}` |  |
| update.image.repository | string | `"alpine/curl"` |  |
| update.image.tag | string | `"8.5.0"` |  |
| update.nodeSelector | object | `{}` |  |
| update.schedule | string | `"37 */2 * * *"` |  |
| update.tolerations[0].effect | string | `"NoSchedule"` |  |
| update.tolerations[0].key | string | `"workload"` |  |
| update.tolerations[0].operator | string | `"Equal"` |  |
| update.tolerations[0].value | string | `"disruptive"` |  |
