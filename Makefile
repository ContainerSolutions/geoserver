IMAGE=containersol/geoserver

build:
	docker run --rm -v "${PWD}":/usr/src/myapp -w /usr/src/myapp -e CGO_ENABLED=0 -e GOOS=linux golang:1.7.0 go build -a -installsuffix cgo -o geoserver .
	docker build -t ${IMAGE} .	

push: build
	docker push ${IMAGE}

all: build
