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
  value       = [for bucket in google_storage_bucket.bucket : bucket.url]
  description = "List of buckets created"
}
