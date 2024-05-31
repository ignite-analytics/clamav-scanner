# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account
resource "google_service_account" "service_account" {
  project                      = var.project_id
  account_id                   = var.name
  display_name                 = var.sa_description
  description                  = var.sa_description
  disabled                     = false
  create_ignore_already_exists = true
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam#google_project_iam_member
resource "google_project_iam_member" "service_account" {
  project = var.project_id
  role    = "roles/eventarc.eventReceiver"
  member  = google_service_account.service_account.member
}

resource "google_project_iam_member" "gcs_account" {
  project = var.project_id
  role    = "roles/pubsub.publisher"
  member  = data.google_storage_project_service_account.gcs_account.member
}
