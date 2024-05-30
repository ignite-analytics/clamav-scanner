variable "project_id" {
  description = "The project ID to deploy resources to"
}

variable "region" {
  description = "The region to deploy resources to"
}

variable "name" {
  description = "Service account ID"
}

variable "sa_display_name" {
  description = "Service Account display name"
  default     = "ClamAV Scanner"
}

variable "sa_description" {
  description = "Service Account description"
  default     = "Service account for ClamAV Scanner"
}

variable "keyring_location" {
  description = "Location of the keyring"
  default     = "europe"
}

variable "crypto_key_rotation_period" {
  description = "The rotation period of the key. Default is 90 days"
  default     = "7776000s"
}

variable "storage" {
  description = "ClamAv Scanner storage configuration"
  default = {
    log_bucket = "my-log-bucket"
    location   = "EU"
  }
}