
resource "kubernetes_deployment" "ddns" {
    metadata {
        name = "ddns-deployment"
        labels = {
            test = "DaisyApp"
        }
    }

    spec {
        replicas = 1

        selector {
            match_labels = {
              test = "DaisyApp"
            }
        }

        template {
            metadata {
              labels = {
                test = "DaisyApp"
              }
            }

            spec {
                container {
                    image = "trevatk/daisy:v0.0.1"
                    name = "daisy"

                    env {
                      name = var.env_dservice_config_name
                      value = var.env_dservice_config_value
                    }

                    volume_mount {
                      name = "config-volume"
                      mount_path = "/local/ddns"
                    }

                    resources {
                        limits = {
                            cpu = "0.3"
                            memory = "512Mi"
                        }
                    }

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = var.container_http_port
                        }

                        initial_delay_seconds = 3
                        period_seconds = 3
                    }
                }

                volume {
                  name = "config-volume"
                  config_map {
                    name = kubernetes_config_map.ddns.metadata.0.name
                  }
                }
            }
        }
    }
}