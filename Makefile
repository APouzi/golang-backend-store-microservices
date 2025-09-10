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
		--remove checkout_session:payment_intent_data.shipping \
		--add checkout_session:automatic_tax.enabled=true \
		--add checkout_session:customer_update[shipping]=auto \
		--add checkout_session:customer=cus_T0z8warRsQ8Viw \
		--add "checkout_session:shipping_address_collection[allowed_countries][0]=US" \
		--add "checkout_session:shipping_address_collection[allowed_countries][1]=CA" \
		--add checkout_session:metadata.itemsizeqty_1=2 \
		--add checkout_session:metadata.itemsizeqty_2=1 \
		--add checkout_session:metadata.itemsizeqty_3=5 \
		--add checkout_session:metadata.itemsizeqty_4=3 \
		--add checkout_session:metadata.itemsizeqty_5=4 \
		--add checkout_session:metadata.itemsizeqty_6=2 \
		--add checkout_session:metadata.itemsizeqty_7=1 \
		--add checkout_session:metadata.itemsizeqty_8=6

create-stripe-customer:
	stripe customers create   -d "address[line1]=1600 Amphitheatre Pkwy"   -d "address[city]=Mountain View"   -d "address[state]=CA"   -d "address[postal_code]=94043"   -d "address[country]=US"