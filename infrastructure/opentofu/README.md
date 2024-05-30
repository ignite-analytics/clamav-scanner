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
| [google_kms_crypto_key.crypto_key](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key) | resource |
| [google_kms_crypto_key_iam_member.crypto_key](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key_iam_member) | resource |
| [google_kms_key_ring.keyring](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_key_ring) | resource |
| [google_project_iam_member.service_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_iam_member) | resource |
| [google_service_account.service_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_storage_bucket.bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket) | resource |
| [google_storage_bucket_iam_member.member](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam_member) | resource |
| [google_storage_project_service_account.gcs_account](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/storage_project_service_account) | data source |

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| crypto\_key\_rotation\_period | The rotation period of the key. Default is 90 days | `string` | no |
| keyring\_location | Location of the keyring | `string` | no |
| name | Service account ID | `any` | yes |
| project\_id | The project ID to deploy resources to | `any` | yes |
| region | The region to deploy resources to | `any` | yes |
| sa\_description | Service Account description | `string` | no |
| sa\_display\_name | Service Account display name | `string` | no |
| storage | ClamAv Scanner storage configuration | `map` | no |

## Outputs

| Name | Description |
|------|-------------|
| crypto\_key\_id | Crypto key ID |
| keyring\_id | Keyring ID |
| service\_account\_email | Service account email |
| service\_account\_id | Service account id |
| service\_account\_unique\_id | Service account unique id |
| storage\_bucket\_id | List of buckets created |
<!-- END_TF_DOCS -->
