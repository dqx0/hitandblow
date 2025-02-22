provider "google" {
  project = "gcp-dqx0"
  region  = "us-central1"
}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth[0].cluster_ca_certificate)
}

data "google_client_config" "default" {}

resource "google_container_cluster" "primary" {
  name              = "hitandblow-cluster"
  location          = "us-central1"      # リージョン名を指定。
  enable_autopilot  = true               # Autopilot を有効化。
}

resource "kubernetes_deployment" "hitandblow" {
  metadata {
    name      = "hitandblow"
    namespace = "default"
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "hitandblow"
      }
    }
    template {
      metadata {
        labels = {
          app = "hitandblow"
        }
      }
      spec {
        container {
          name  = "hitandblow"
          image = "dqx0/hitandblow:v1.0.0"
          image_pull_policy = "Always"
          port {
            container_port = 8080
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "hitandblow" {
  metadata {
    name      = "hitandblow-service"
    namespace = "default"
  }
  spec {
    selector = {
      app = "hitandblow"
    }
    port {
      protocol    = "TCP"
      port        = 80
      target_port = 8080
    }
    type = "LoadBalancer"
  }
}