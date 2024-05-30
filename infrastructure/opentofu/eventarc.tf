# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/eventarc_trigger
resource "google_eventarc_trigger" "default" {
  for_each = { for index, s in var.scan_config : s.bucket_name => s }
  name     = "${each.value.bucket_name}-scan"
  location = var.region

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

  labels = {
    manager = "opentofu"
    service = var.name
  }

  depends_on = [google_cloud_run_v2_service.default]
}
