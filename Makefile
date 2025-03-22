.PHONY: mock
mock:
	@mockgen -source=./aweb/internal/service/user.go -package=svcmocks -destination=./aweb/internal/service/mocks/user.mock.go
	@mockgen -source=./aweb/internal/service/code.go -package=svcmocks -destination=./aweb/internal/service/mocks/code.mock.go
	@mockgen -source=./aweb/internal/service/sms/types.go -package=smsmocks -destination=./aweb/internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./aweb/internal/repository/code.go -package=repomocks -destination=./aweb/internal/repository/mocks/code.mock.go
	@mockgen -source=./aweb/internal/repository/user.go -package=repomocks -destination=./aweb/internal/repository/mocks/user.mock.go
	@mockgen -source=./aweb/internal/repository/dao/user.go -package=daomocks -destination=./aweb/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./aweb/internal/repository/cache/user.go -package=cachemocks -destination=./aweb/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./aweb/internal/repository/cache/code.go -package=cachemocks -destination=./aweb/internal/repository/cache/mocks/code.mock.go
	@mockgen -source=./aweb/pkg/limiter/types.go -package=limitermocks -destination=./aweb/pkg/limiter/mocks/limiter.mock.go
	@go mod tidy