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

HELM_RELEASE_NAME := trader
HELM_CHART := ./charts/trader
HELM_VALUES_PATH := ./trader-values.yaml

.PHONY: helm.diff
helm.diff:
	# helm plugin install https://github.com/databus23/helm-diff
	helm diff --allow-unreleased upgrade \
		--set-file secret.config=$(CONFIG_FILE) \
		--debug \
		--values $(HELM_VALUES_PATH) \
		$(HELM_RELEASE_NAME) \
		$(HELM_CHART)


.PHONY: helm.service helm.service.redis
helm.service:
	helm upgrade --install \
		--set-file secret.config=$(CONFIG_FILE) \
		--values $(HELM_VALUES_PATH) \
		-n $(KUBE_NS) \
		$(HELM_RELEASE_NAME) \
		$(HELM_CHART)

HELM_REDIS_RELEASE_NAME := redis
helm.service.redis:
	# helm repo add bitnami https://charts.bitnami.com/bitnami
	helm upgrade --install \
		--set cluster.enabled=false \
		--set password=${REDIS_PASSWORD} \
		-n redis \
		$(HELM_REDIS_RELEASE_NAME) bitnami/redis

HELM_POSTGRES_RELEASE_NAME := postgres
helm.service.postgres:
	# helm repo add bitnami https://charts.bitnami.com/bitnami
	helm upgrade --install \
		--set postgresqlDatabase=ptcg \
		--set postgresqlUsername=${POSTGRES_USERNAME} \
		--set postgresqlPassword=${POSTGRES_PASSWORD} \
		-n postgres \
		$(HELM_POSTGRES_RELEASE_NAME) bitnami/postgresql 
