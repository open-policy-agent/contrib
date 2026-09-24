# Terraform CLI v1.1.9 or higher is required for Run Task feedback in the CLI.
# Using an earlier CLI version will only show Run Task results in your
# Terraform Cloud workspace run details online.

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
  required_version = ">= 0.14.9"
}
