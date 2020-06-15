
NAME := kubectl-utlz 
VERSION    := 1.0.0 

linux: | local
darwin: | local
local:
ifneq ($(MAKECMDGOALS), darwin)
ifneq ($(MAKECMDGOALS), linux)
	$(error Valid local build targets are "linux" and "darwin")
endif
endif
	GOOS=$(MAKECMDGOALS) GOARCH=amd64 CGO_ENABLED=0 go build -o ./$(NAME) ./main.go
