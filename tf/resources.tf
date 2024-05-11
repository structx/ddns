
resource "kubernetes_persistent_volume_claim" "ddns" {
    metadata {
        name = "ddns-logs"
    }
    spec {
        access_modes = ["ReadWriteMany"]
        resources {
            requests = {
                storage = "1Gi"
            }
        }
        volume_name = "${kubernetes_persistent_volume.ddns.metadata.0.name}"
    }
}

resource "kubernetes_persistent_volume" "ddns" {
    metadata {
        name = "ddns-logs"
    }
    spec {
        capacity = {
          storage = "1Gi"
        }
        access_modes = [ "ReadWriteMany" ]
        persistent_volume_source {
            gce_persistent_disk {
              pd_name = "ddns-logs"
            }
        }
    }
}

resource "kubernetes_config_map" "ddns" {
    metadata {
        name = "service-config"
    }
    data = {
        "server.config.hcl" = file("./build/config/server.config.hcl")
    }
}