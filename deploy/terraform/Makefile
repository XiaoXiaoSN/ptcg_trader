Terraform := terraform

.PHONY: init kube.export kube.ingress

init:
	$(Terraform) init

kube.export:
	# source .env
	# brew install kops
	kops export kubecfg --admin

kube.ingress.install:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.44.0/deploy/static/provider/cloud/deploy.yaml
