# Customize these variables
IMAGE_NAME = qreview
DOCKER_USERNAME = yourusername
TAG = latest
FULL_IMAGE_NAME = $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)


build:
	docker build --no-cache -t $(IMAGE_NAME) .
run:
	docker run --rm $(IMAGE_NAME)	

tag:
	docker tag $(IMAGE_NAME) $(FULL_IMAGE_NAME)

push: tag
	docker push $(FULL_IMAGE_NAME)

all: build push

clean:
	docker rmi $(IMAGE_NAME) $(FULL_IMAGE_NAME) || true
