variable "project_id" {
  description = "The project ID to deploy resources to"
  type        = string
}

variable "region" {
  description = "The region to deploy resources to"
  type        = string
}

variable "name" {
  description = "Service account ID"
  type        = string
}

variable "sa_display_name" {
  description = "Service Account display name"
  type        = string
  default     = "ClamAV Scanner"
}

variable "sa_description" {
  description = "Service Account description"
  type        = string
  default     = "Service account for ClamAV Scanner"
}

variable "keyring_location" {
  description = "Location of the keyring"
  type        = string
}

variable "crypto_key_rotation_period" {
  description = "The rotation period of the key. Default is 90 days"
  type        = string
  default     = "7776000s"
}

variable "storage" {
  description = "ClamAv Scanner storage configuration"
  type = object({
    mirror_bucket     = string
    quarantine_bucket = string
    log_bucket        = string
    location          = string
  })
}

variable "image" {
  description = "ClamAV Scanner image configuration"
  type = object({
    repository = string
    tag        = string
  })
}

variable "port" {
  description = "Port to expose the service on"
  type        = number
  default     = 1337
}

variable "schedule" {
  description = "ClamAV Scanner scheduler configuration"
  type = object({
    cron     = string
    timezone = string
  })
}

variable "scan_config" {
  description = "ClamAV Scanner buckets scan configuration"
  type = list(object({
    bucket_name = string
    location    = string
  }))
}

variable "pubsub" {
  description = "ClamAV Scanner Pub/Sub configuration"
  type = object({
    topic        = string
    subscription = string
    subscribers  = optional(list(string))
  })
}

variable "labels" {
  description = "ClamAV Scanner labels"
  type        = map(string)
  default = {
    manager = "opentofu"
    service = "clamav-scanner"
  }
}
