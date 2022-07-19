terraform {
  required_providers {
    scaffolding = {
      source  = "local/hashicorp/scaffolding"
      version = "0.1.0"
    }
  }
}

resource "scaffolding_example" "example" {
  configurable_attribute = "some-value"
}
