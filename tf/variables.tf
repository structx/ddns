
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

variable "config_file_path" {
    type = string
    default = "./build/config/server.config.hcl"
    description = "local config file"
}

variable env_dservice_config_name {
    type = string
    default = "DSERVICE_CONFIG"
}

variable env_dservice_config_value {
    type = string
    default = "/local/ddns/server.config.hcl"
}