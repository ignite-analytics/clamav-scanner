terraform {
  required_providers {
    google = {
      source  = "opentofu/google"
      version = ">= 6.0, < 7.0"
    }
    google-beta = {
      source  = "opentofu/google-beta"
      version = ">= 6.0, < 7.0"
    }
  }
}

# Provider required variables
provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}
