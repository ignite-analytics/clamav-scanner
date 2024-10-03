# https://search.opentofu.org/provider/hashicorp/google/latest/docs/datasources/project
data "google_project" "project" {
  project_id = var.project_id
}

# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/pubsub_topic
resource "google_pubsub_topic" "default" {
  name   = var.pubsub.topic
  labels = var.labels
}

resource "google_pubsub_topic" "dead-letter" {
  name   = "${var.pubsub.topic}-dead-letter"
  labels = var.labels
}

# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/pubsub_subscription
resource "google_pubsub_subscription" "default" {
  name                       = var.pubsub.subscription
  topic                      = google_pubsub_topic.default.name
  ack_deadline_seconds       = 20
  retain_acked_messages      = true
  message_retention_duration = "1200s"
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead-letter.id
    max_delivery_attempts = 5
  }

  labels = var.labels

  depends_on = [google_pubsub_topic.default, google_pubsub_topic.dead-letter]
}

resource "google_pubsub_subscription" "dead-letter" {
  name                 = "${var.pubsub.subscription}-dead-letter"
  topic                = google_pubsub_topic.dead-letter.name
  ack_deadline_seconds = 20

  retry_policy {
    minimum_backoff = "1s"
    maximum_backoff = "10s"
  }

  labels = var.labels

  depends_on = [google_pubsub_topic.dead-letter]

}

# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/pubsub_topic_iam
resource "google_pubsub_topic_iam_member" "default" {
  for_each = toset([google_pubsub_topic.default.name, google_pubsub_topic.dead-letter.name])
  project  = var.project_id
  topic    = each.value
  role     = "roles/pubsub.editor"
  member   = google_service_account.service_account.member
}

resource "google_pubsub_topic_iam_member" "dead-letter" {
  project = var.project_id
  topic   = google_pubsub_topic.dead-letter.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

# https://search.opentofu.org/provider/hashicorp/google/latest/docs/resources/pubsub_subscription_iam
resource "google_pubsub_subscription_iam_member" "default" {
  for_each     = toset([google_pubsub_subscription.default.name, google_pubsub_subscription.dead-letter.name])
  project      = var.project_id
  subscription = each.value
  role         = "roles/pubsub.editor"
  member       = google_service_account.service_account.member
}

resource "google_pubsub_subscription_iam_member" "subscriber" {
  for_each     = toset(var.pubsub.subscribers)
  project      = var.project_id
  subscription = google_pubsub_subscription.default.name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${each.value}"

}
resource "google_pubsub_subscription_iam_member" "dead-letter" {
  project      = var.project_id
  subscription = google_pubsub_subscription.dead-letter.name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}
