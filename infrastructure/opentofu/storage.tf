locals {
  storage = ["${var.name}-mirror", "${var.name}-quarantine"]
}

#  https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_key_ring
resource "google_kms_key_ring" "keyring" {
  name     = "${var.name}-keyring"
  project  = var.project_id
  location = var.keyring_location
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key
resource "google_kms_crypto_key" "crypto_key" {
  name            = "${var.name}-key"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "7776000s"
  purpose         = "ENCRYPT_DECRYPT"

  labels = {
    manager = "opentofu",
    service = var.name
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_kms_key_ring.keyring]
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/storage_project_service_account
data "google_storage_project_service_account" "gcs_account" {
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_kms_crypto_key_iam
resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"
  member        = "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"

  depends_on = [google_kms_crypto_key.crypto_key]
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket
resource "google_storage_bucket" "bucket" {
  for_each = toset(local.storage)
  name     = each.key
  location = var.storage.location

  force_destroy               = false
  public_access_prevention    = "enforced"
  uniform_bucket_level_access = true

  logging {
    log_bucket = var.storage.log_bucket
  }

  encryption {
    default_kms_key_name = google_kms_crypto_key.crypto_key.id
  }

  labels = {
    manager = "opentofu",
    service = var.name
  }

  depends_on = [google_kms_crypto_key.crypto_key]
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam#google_storage_bucket_iam_member
resource "google_storage_bucket_iam_member" "member" {
  for_each = toset(local.storage)
  bucket   = each.key
  role     = "roles/storage.objectUser"
  member   = google_service_account.service_account.member

  depends_on = [google_storage_bucket.bucket]
}
