.PHONY: build

all: dev

dev:
	bee run -downdoc=true -gendoc=true

start: build
	docker-compose up --build

build:
	docker build . -t endaaman/api.endaaman.me

push: build
	docker push endaaman/api.endaaman.me

pull:
	docker pull endaaman/api.endaaman.me
