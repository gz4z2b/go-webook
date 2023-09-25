.PHONY: docker
docker: 
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/repository/interface.go -package=repomocks -destination=./internal/repository/mocks/userRepo.mock.go
	@mockgen -source=./internal/repository/dao/interface.go -package=daomocks -destination=./internal/repository/dao/mocks/userDao.mock.go
	@mockgen -source=./internal/repository/cache/interface.go -package=cachemocks -destination=./internal/repository/cache/mocks/userCache.mock.go
	@mockgen -package=redismocks -destination=./internal/repository/cache/mocks/redismocks/redis.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy

	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f gz4z2b/webook:v0.0.1 
	@docker build -t gz4z2b/webook:v0.0.1 .
