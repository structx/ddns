
resource "kubernetes_namespace" "ddns" {
    metadata {
        name = "development"
    }
}