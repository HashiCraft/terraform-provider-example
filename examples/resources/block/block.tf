terraform {
  required_providers {
    example = {
      source  = "local/hashicraft/example"
      version = "0.1.0"
    }
  }
}

provider "example" {
  url = "http://localhost:8080"
}

resource "example_block" "brick" {
  x = 24
  y = 63
  z = 64
  #material = "minecraft:oak_log"
  material = "minecraft:stone"
}
