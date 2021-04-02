Terraform := terraform

.PHONY: init kube.export

init:
	$(Terraform) init

kube.export:
	# source .env
	# brew install kops
	kops export kubecfg --admin