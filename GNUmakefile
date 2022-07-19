default: testacc

# Run acceptance tests
.PHONY: testacc
arch=$(shell /bin/bash ./arch.sh)

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

name = scaffolding
organization = hashicorp
version = 0.1.0
log_level = info

build:
	go build -o bin/terraform-provider-$(name)_v$(version)

install: build clean
	mkdir -p ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)
	mv bin/terraform-provider-$(name)_v$(version) ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)/

clean:
	rm -rf examples/resources/scaffolding_example/.terraform*
	rm -rf examples/resources/scaffolding_example/terraform.tfstate*

init:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/scaffolding_example init

plan:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/scaffolding_example plan

apply:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/scaffolding_example apply -auto-approve

destroy:
	TF_LOG=$(log_level) terraform -chdir=examples/resources/scaffolding_example destroy -auto-approve
