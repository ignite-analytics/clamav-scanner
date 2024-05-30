output "service_account_email" {
  value       = google_service_account.service_account.email
  description = "Service account email"
}

output "service_account_id" {
  value       = google_service_account.service_account.id
  description = "Service account id"
}

output "service_account_unique_id" {
  value       = google_service_account.service_account.unique_id
  description = "Service account unique id"
}

output "keyring_id" {
  value       = google_kms_key_ring.keyring.id
  description = "Keyring ID"
}

output "crypto_key_id" {
  value       = google_kms_crypto_key.crypto_key.id
  description = "Crypto key ID"
}

output "storage_bucket_id" {
  value       = [for b in google_storage_bucket.bucket : b.url]
  description = "List of buckets created"
}

output "service_id" {
  value       = google_cloud_run_v2_service.default.id
  description = "Cloud Run service ID"
}

output "service_uri" {
  value       = google_cloud_run_v2_service.default.uri
  description = "Cloud Run service main URI for serving traffic"
}

output "scheduler_id" {
  value       = google_cloud_scheduler_job.job.id
  description = "Cloud Scheduler job ID"
}

output "eventarc_trigger_id" {
  value       = [for e in google_eventarc_trigger.default : e.id]
  description = "Eventarc trigger ID"
}
