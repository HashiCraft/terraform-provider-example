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
  x = 25
  y = 63
  z = 64
  #material = "minecraft:oak_log"
  material = "minecraft:stone"
}

data "example_block" "brick" {
  id = resource.example_block.brick.id
}

output "brick_coordinates" {
  value = {
    "x" = data.example_block.brick.x
    "y" = data.example_block.brick.y
    "z" = data.example_block.brick.z
  }
}
