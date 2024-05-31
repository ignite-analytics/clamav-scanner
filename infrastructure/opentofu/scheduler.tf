# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_scheduler_job
resource "google_cloud_scheduler_job" "job" {
  name             = "${var.name}-update"
  description      = "ClamAV Scanner scheduled update job"
  schedule         = var.schedule.cron
  time_zone        = var.schedule.timezone
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_v2_service.default.uri}/update"

    oidc_token {
      service_account_email = google_service_account.service_account.email
    }
  }
}
