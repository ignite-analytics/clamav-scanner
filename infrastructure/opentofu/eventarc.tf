# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/eventarc_trigger
resource "google_eventarc_trigger" "trigger" {
  for_each = { for index, s in var.scan_config : s.bucket_name => s }
  name     = "${each.value.bucket_name}-scan"
  location = var.region

  service_account = google_service_account.service_account.email

  matching_criteria {
    attribute = "type"
    value     = "google.cloud.storage.object.v1.finalized"
  }

  matching_criteria {
    attribute = "bucket"
    value     = each.value.bucket_name
  }

  destination {
    cloud_run_service {
      service = google_cloud_run_v2_service.default.name
      region  = var.region
    }
  }

  labels = var.labels

  depends_on = [google_cloud_run_v2_service.default]
}

# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/storage_bucket_iam#google_storage_bucket_iam_member
resource "google_storage_bucket_iam_member" "trigger" {
  for_each = { for index, s in var.scan_config : s.bucket_name => s }
  bucket   = each.key
  role     = "roles/storage.objectUser"
  member   = google_service_account.service_account.member

  depends_on = [google_eventarc_trigger.trigger]
}
