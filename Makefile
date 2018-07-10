

UID=$(shell id -u) 
GID=$(shell id -g) 


define usercompose
version: "2.1"
services:
  base_config:
    build:
      args:
        GID: $(GID)
        UID: $(UID)
endef
export usercompose


run: .user.compose.yml
	docker-compose up --remove-orphans --build -d

.user.compose.yml:
	@echo "$$usercompose" > $@
	echo $@ created

build: containers/prod/app containers/prod/ui
	docker-compose build
	cd containers/prod/ && docker build -t smileinnovation/imannotate .

containers/prod/ui:
	docker-compose run --rm ui ng build --prod
	mv ui/dist containers/prod/ui

containers/prod/app:
	echo 'go build -ldflags "-linkmode external -extldflags -static" -o app.bin' > builder.sh
	docker-compose run --rm api sh ./builder.sh
	mv app.bin containers/prod/app
	rm builder.sh


clean:
	rm -rf ui/dist
	rm -rf containers/prod/ui containers/prod/app
	rm -f gin-bin app.bin

clean-docker:
	docker-compose down --remove-orphans


clean-volumes:
	docker-compose down -v --remove-orphans