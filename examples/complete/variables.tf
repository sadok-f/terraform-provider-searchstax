variable "ssx_username" {
  type        = string
  description = "SearchStax account email."
}

variable "ssx_pwd" {
  type        = string
  sensitive   = true
  description = "SearchStax account password."
}

variable "ssx_host" {
  type        = string
  default     = "https://app.searchstax.com/api/rest/v2"
  description = "SearchStax REST API base URL."
}

variable "account_name" {
  type        = string
  description = "SearchStax tenant account name."
}

variable "deployment_uid" {
  type        = string
  description = "Existing deployment UID to configure."
}

variable "office_cidr" {
  type        = string
  default     = "203.0.113.0/24"
  description = "CIDR allowed through the deployment IP filter."
}

variable "deployment_tags" {
  type        = list(string)
  default     = ["terraform", "managed"]
  description = "Tags applied to the deployment."
}
