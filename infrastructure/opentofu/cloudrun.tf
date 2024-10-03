# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/cloud_run_v2_service
resource "google_cloud_run_v2_service" "default" {
  name     = var.name
  project  = var.project_id
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account       = google_service_account.service_account.email
    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
    containers {
      image = "${var.image.repository}:${var.image.tag}"

      resources {
        startup_cpu_boost = true
        limits = {
          cpu    = "1"
          memory = "4Gi"
        }
      }

      ports {
        container_port = var.port
      }

      startup_probe {
        initial_delay_seconds = 60
        timeout_seconds       = 10
        period_seconds        = 10
        failure_threshold     = 30
        http_get {
          path = "/health"
          port = 1337
        }
      }

      liveness_probe {
        http_get {
          path = "/health"
          port = 1337
        }
      }

      env {
        name  = "MIRROR_BUCKET"
        value = var.storage.mirror_bucket
      }

      env {
        name  = "QUARANTINE_BUCKET"
        value = var.storage.quarantine_bucket
      }

      env {
        name  = "LISTEN_ADDRESS"
        value = ":${var.port}"
      }

      env {
        name  = "PUBSUB_TOPIC"
        value = var.pubsub.topic
      }

      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
    }

    max_instance_request_concurrency = 20
    scaling {
      min_instance_count = 1
      max_instance_count = 5
    }
  }

  labels = var.labels

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/cloud_run_v2_service_iam#google_cloud_run_v2_service_iam_member
resource "google_cloud_run_v2_service_iam_member" "member" {
  project  = var.project_id
  location = var.region
  name     = var.name
  role     = "roles/run.invoker"
  member   = google_service_account.service_account.member

  depends_on = [google_cloud_run_v2_service.default]
}
