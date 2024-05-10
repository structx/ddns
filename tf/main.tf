
provider "kubernetes" {
    config_path = "~/.kube/config"
}

resource "kubernetes_namespace" "development" {
    metadata {
        name = "development"
    }
}

resource "kubernetes_deployment" "development" {
    metadata {
        name = var.app_name
        namespace = kubernetes_namespace.development.metadata.0.name
    }
    spec {
        replicas = 1
        selector {
            match_labels = {
                app = var.app_name
            }
        }
        template {
            metadata {
              labels = {
                app = var.app_name
              }
            }
            spec {
                container {
                    image = "trevatk/daisy:v0.0.1"
                    name = var.app_name
                    port {
                        container_port = var.container_http_port
                    }
                }
            }
        }
    }
}

resource "kubernetes_service" "development" {
    metadata {
        name = var.app_name
        namespace = kubernetes_namespace.development.metadata.0.name
    }
    spec {
        selector = {
            app = kubernetes_deployment.development.spec.0.template.0.metadata.0.labels.app
        }
        type = "NodePort"
        port {
            node_port = 30201
            port = var.container_http_port
            target_port = var.container_http_port
        }
    }
}