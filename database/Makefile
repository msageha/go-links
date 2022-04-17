MYSQL_IMAGE_NAME='mysql'

.PHONY: build_image
build_image:
	docker build . -f dockers/Dockerfile -t $(MYSQL_IMAGE_NAME)

.PHONY: run_image
run_image:
	docker run -d -p 13306:3306 $(MYSQL_IMAGE_NAME)
