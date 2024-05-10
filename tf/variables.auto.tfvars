
variable "container_http_port" {
    type = number
    default = 8081
    description = "http server default port"
}

variable "container_grpc_port" {
    type = number
    default = 50051
    description = "gRPC server default port"
}

variable "app_name" {
    type = string
    default = "daisy"
    description = "service name"
}