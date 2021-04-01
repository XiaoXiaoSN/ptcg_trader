CONFIG_FILE := app.yaml


##############################
# kubectl
##############################

KUBE_CONTEXT := ptcg.10oz.tw
KUBE_NS := ptcg

.PHONY: kube.activate
kube.activate:
	kubectl config use-context $(KUBE_CONTEXT)
	kubectl config set-context --current --namespace=$(KUBE_NS)


##############################
# helm
##############################

RELEASE_NAME := trader
HELM_CHART := ./charts/trader
HELM_VALUES_PATH := ./trader-values.yaml

.PHONY: helm.diff
helm.diff:
	# helm plugin install https://github.com/databus23/helm-diff
	helm diff --allow-unreleased upgrade \
		--set-file secret.config=$(CONFIG_FILE) \
		--debug \
		--values $(HELM_VALUES_PATH) \
		$(RELEASE_NAME) \
		$(HELM_CHART)
