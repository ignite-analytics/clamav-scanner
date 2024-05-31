Here you can find example of OpenTofu setup to deploy the service on Google Cloud Run. All you need to do
is to configure your variables in file `variables.tf` and run `tofu apply`.

## What is created

-   Service Account and adds IAM permissions for Eventarc events
-   Creates mirror and quarantine buckets. Bucket used for logging is expect to be already created
-   Creates KMS keyring and crypto key used for encrypting mirror and quarantine buckets
-   Creates Cloud Run service 2nd generation
-   Creates Cloud Scheduler job that calls `update` endpoint
-   Creates Eventarc trigger that listens for `google.cloud.storage.object.v1.finalized` events. This is created by
    populating `scan_config` in `variables.tf`. It's a list of objects where you declare name and location. That is
    the place where you register buckets to be scanned for malware. Event triggers are created in same location as the
    buckets.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.7 |
| google | >= 5.0 |
| google-beta | >= 5.0 |

## Providers

| Name | Version |
|------|---------|
| google | 5.31.1 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_cloud_run_v2_service.default](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_service) | resource |
| [google_cloud_run_v2_service_iam_member.member](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_service_iam_member) | resource |
| [google_cloud_scheduler_job.job](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_scheduler_job) | resource |
| [google_eventarc_trigger.trigger](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/eventarc_trigger) | resource |
| [google_kms_crypto_key.crypto_key](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key) | resource |
| [google_kms_crypto_key_iam_member.crypto_key](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key_iam_member) | resource |
| [google_kms_key_ring.keyring](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_key_ring) | resource |
| [google_project_iam_member.gcs_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_iam_member) | resource |
| [google_project_iam_member.service_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_iam_member) | resource |
| [google_service_account.service_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_storage_bucket.bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket) | resource |
| [google_storage_bucket_iam_member.member](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam_member) | resource |
| [google_storage_bucket_iam_member.trigger](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam_member) | resource |
| [google_storage_project_service_account.gcs_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/storage_project_service_account) | data source |

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| crypto\_key\_rotation\_period | The rotation period of the key. Default is 90 days | `string` | no |
| image | ClamAV Scanner image configuration | <pre>object({<br>    repository = string<br>    tag        = string<br>  })</pre> | yes |
| keyring\_location | Location of the keyring | `string` | yes |
| name | Service account ID | `string` | yes |
| port | Port to expose the service on | `number` | no |
| project\_id | The project ID to deploy resources to | `string` | yes |
| region | The region to deploy resources to | `string` | yes |
| sa\_description | Service Account description | `string` | no |
| sa\_display\_name | Service Account display name | `string` | no |
| scan\_config | ClamAV Scanner buckets scan configuration | <pre>list(object({<br>    bucket_name = string<br>    location    = string<br>  }))</pre> | yes |
| schedule | ClamAV Scanner scheduler configuration | <pre>object({<br>    cron     = string<br>    timezone = string<br>  })</pre> | yes |
| storage | ClamAv Scanner storage configuration | <pre>object({<br>    mirror_bucket     = string<br>    quarantine_bucket = string<br>    log_bucket        = string<br>    location          = string<br>  })</pre> | yes |

## Outputs

| Name | Description |
|------|-------------|
| crypto\_key\_id | Crypto key ID |
| eventarc\_trigger\_id | Eventarc trigger ID |
| keyring\_id | Keyring ID |
| scan\_buckets | List of buckets to scan |
| scheduler\_id | Cloud Scheduler job ID |
| service\_account\_email | Service account email |
| service\_account\_id | Service account id |
| service\_account\_unique\_id | Service account unique id |
| service\_id | Cloud Run service ID |
| service\_uri | Cloud Run service main URI for serving traffic |
| storage\_bucket\_id | List of buckets created |
<!-- END_TF_DOCS -->
