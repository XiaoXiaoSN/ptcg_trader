
.PHONY: create_order stress.create_order.prod

stress.create_order:
	# brew install k6
	k6 run create_order.js

TRADER_URL := http://trader.ptcg.10oz.tw
stress.create_order.prod:
	# brew install k6
	k6 run -e TRADER_URL=$(TRADER_URL) create_order.js
