IMAGE_VERSION=v0.7
REPO_URL=172.16.16.172:12380

vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

check: vet

format:
	#go get golang.org/x/tools/cmd/goimports
	find . -name '*.go' | grep -Ev 'vendor|thrift_gen' | xargs goimports -w

build: check

all: dev api gateway consumer pi benchmark

dev: check
	@>&2 echo "Great, all tests passed."

check: format vet

gateway:
	sh ./scripts/build_gateway.sh

consumer:
	sh ./scripts/build_consumer.sh

api:
	sh ./scripts/build_api.sh

pi:
	sh ./scripts/build_sample_pi.sh

benchmark:
	sh ./scripts/build_sample_benchmark.sh

docker: docker-gateway docker-consumer docker-pi docker-benchmark docker-api buildsucc

docker-gateway: gateway
	@docker build -f docker/gateway.Dockerfile .  -t $(REPO_URL)/cudgx/gateway:$(IMAGE_VERSION)

docker-consumer: consumer
	@docker build -f docker/consumer.Dockerfile  .  -t $(REPO_URL)/cudgx/consumer:$(IMAGE_VERSION)

docker-api: api
	@docker build -f docker/api.Dockerfile  .  -t $(REPO_URL)/cudgx/api:$(IMAGE_VERSION)

docker-pi: pi
	@docker build -f docker/pi.Dockerfile . -t $(REPO_URL)/cudgx/sample-pi:$(IMAGE_VERSION)

docker-benchmark: benchmark
	@docker build -f docker/benchmark.Dockerfile . -t $(REPO_URL)/cudgx/sample-benchmark:$(IMAGE_VERSION)


docker-push: docker push-gateway push-consumer push-pi push-api push-benchmark


push-gateway: docker-gateway
	docker push $(REPO_URL)/cudgx/gateway:$(IMAGE_VERSION)

push-consumer: docker-consumer
	docker push $(REPO_URL)/cudgx/consumer:$(IMAGE_VERSION)

push-api: docker-api
	docker push $(REPO_URL)/cudgx/api:$(IMAGE_VERSION)

push-pi: docker-pi
	docker push $(REPO_URL)/cudgx/sample-pi:$(IMAGE_VERSION)

push-benchmark: docker-benchmark
	docker push $(REPO_URL)/cudgx/sample-benchmark:$(IMAGE_VERSION)

# Quick start
# Pull images from dockerhub and run
docker-run-linux:
	sh ./run-for-linux.sh

docker-run-mac:
	sh ./run-for-mac.sh

docker-container-stop:
	docker ps -aq | xargs docker stop
	docker ps -aq | xargs docker rm

docker-image-rm:
	docker image prune --force --all

# Immersive experience
# Compile and run by docker-compose
docker-compose-start:
	docker-compose up -d

docker-compose-stop:
	docker-compose down

docker-compose-build:
	docker-compose build

#USE make TARGET version=xx override version
version ?= latest

docker-tag:
	docker tag cudgx_api:latest galaxyfuture/cudgx-api:${version}
	docker tag cudgx_gateway:latest galaxyfuture/cudgx-gateway:${version}
	docker tag cudgx_consumer:latest galaxyfuture/cudgx-consumer:${version}
	docker tag cudgx_sample_pi:latest galaxyfuture/cudgx-sample-pi:${version}
	docker tag cudgx_sample_benchmark:latest galaxyfuture/cudgx-sample-benchmark:${version}

docker-push-hub:
	docker push galaxyfuture/cudgx-api:${version}
	docker push galaxyfuture/cudgx-gateway:${version}
	docker push galaxyfuture/cudgx-consumer:${version}
	docker push galaxyfuture/cudgx-sample-pi:${version}
	docker push galaxyfuture/cudgx-sample-benchmark:${version}

docker-hub-all: docker-compose-build docker-tag docker-push-hub






