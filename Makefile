.PHONY: docker
docker: 
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f gz4z2b/webook:v0.0.1 
	@docker build -t gz4z2b/webook:v0.0.1 .