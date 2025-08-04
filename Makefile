up:
	docker compose build; docker compose up

down:
	docker compose down

restart:
	docker kill golang-store-microservices-app-1; docker-compose -f ./docker-compose.yml build; docker-compose -f ./docker-compose.yml up -d

nuke:
	docker compose down -v; docker rm -vf $$(docker ps -aq); docker rmi -f $$(docker images -aq); docker image prune -f; docker volume prune -f; docker system prune -f

trigger-start:
	stripe listen --forward-to localhost:8000/stripe/webhook/payment-confirmation




trigger-meta:
	stripe trigger checkout.session.completed \
		--add checkout_session:metadata.itemsizeqty_1=2 \
		--add checkout_session:metadata.itemsizeqty_2=1 \



trigger-meta-2:
	stripe trigger checkout.session.completed \
		--add checkout_session:metadata.itemsizeqty_101=2 \
		--add checkout_session:metadata.itemsizeqty_102=1 \
		--add checkout_session:metadata.itemsizeqty_103=5 \
		--add checkout_session:metadata.itemsizeqty_104=3 \
		--add checkout_session:metadata.itemsizeqty_105=4 \
		--add checkout_session:metadata.itemsizeqty_106=2 \
		--add checkout_session:metadata.itemsizeqty_107=1 \
		--add checkout_session:metadata.itemsizeqty_108=6 \
		--add checkout_session:metadata.itemsizeqty_109=2 \
		--add checkout_session:metadata.itemsizeqty_110=8
