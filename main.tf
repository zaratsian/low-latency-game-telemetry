// Copyright 2023 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

terraform {
  required_providers {
    google = {
      source = "google"
      version = "~> 4.47"
    }
    google-beta = {
      source = "google-beta"
      version = "~> 4.47"
    }
  }
}

/******************************************************

Enable Google Cloud Services

*******************************************************/

// Get google project 
data "google_project_xxxxxxx" "project" {
}

resource "google_project_service" "gcp_services" {
  for_each = toset(var.gcp_project_services)
  project  = "${var.gcp_project_id}"
  service  = each.value

  disable_on_destroy = false
}

/******************************************************

Google Cloud GKE Resources

*******************************************************/

resource "google_container_cluster" "udp-gke-cluster" {
  name     = var.gke_config.cluster_name
  location = var.gke_config.location
  
  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  # See issue: https://github.com/hashicorp/terraform-provider-google/issues/10782
  ip_allocation_policy {}

  # Enabling Autopilot for this cluster
  enable_autopilot = true

  # Private IP Config
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
  }

  depends_on = [google_project_service.gcp_services]
}

resource "google_service_account" "app-service-account" {
  project      = var.gcp_project_id
  account_id   = var.service_account_config_app.name
  display_name = var.service_account_config_app.description
}

resource "kubernetes_service_account" "k8s-service-account" {
  metadata {
    name      = var.k8s_service_account_id
    namespace = "default"
    annotations = {
      "iam.gke.io/gcp-service-account" : "${google_service_account.app-service-account.email}"
    }
  }
}

resource "kubernetes_secret_v1" "k8s-service-account" {
  metadata {
    annotations = {
      "kubernetes.io/service-account.name" = "k8s-service-account"
    }
    name = "k8s-service-account"
  }

  type = "kubernetes.io/service-account-token"

  depends_on = [kubernetes_service_account.k8s-service-account]
}

data "google_iam_policy" "spanner-policy" {
  binding {
    role = "roles/iam.workloadIdentityUser"
    members = [
      "serviceAccount:${var.project}.svc.id.goog[default/${kubernetes_service_account.k8s-service-account.metadata[0].name}]"
    ]
  }
}

resource "google_service_account_iam_policy" "app-service-account-iam" {
  service_account_id = google_service_account.app-service-account.name
  policy_data        = data.google_iam_policy.spanner-policy.policy_data
}

/******************************************************

Google Cloud Spanner Resources

*******************************************************/

resource "google_spanner_instance" "main" {
  config           = "${var.SPANNER_CONFIG}"
  display_name     = "${var.SPANNER_INSTANCE}"
  processing_units = "${var.SPANNER_PROCESSING_UNITS}"
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.main.name
  name     = "${var.SPANNER_DATABASE}"
  version_retention_period = "3d"
  ddl = [
    "CREATE TABLE Player (PlayerId STRING(MAX) NOT NULL,PlayerName STRING(MAX),PlayerAddress JSON,LastLogin TIMESTAMP) PRIMARY KEY(PlayerId)",

    "CREATE TABLE PlayerInventoryItem (PlayerId STRING(MAX) NOT NULL,ItemId STRING(MAX) NOT NULL,ItemName STRING(MAX) NOT NULL,Amount INT64,LastUpdate TIMESTAMP) PRIMARY KEY(PlayerId, ItemId),INTERLEAVE IN PARENT Player ON DELETE CASCADE",

    "CREATE TABLE GameTelemetry (eventId STRING(MAX),LastUpdated TIMESTAMP,PlayerId STRING(MAX) NOT NULL,PlayerName STRING(MAX),Event STRING(MAX)) PRIMARY KEY(eventId)",

    "CREATE TABLE StoreInventoryItem (ItemId STRING(MAX) NOT NULL,ItemName STRING(MAX) NOT NULL,PRICE FLOAT64,LastUpdate TIMESTAMP) PRIMARY KEY(ItemId)",

    "CREATE TABLE Venues (VenueId INT64 NOT NULL,VenueName STRING(1024),VenueAddress STRING(1024),VenueFeatures JSON,DateOpened DATE) PRIMARY KEY(VenueId)",
  ]
  deletion_protection = false
}