bin     := "bin"
entry   := "cmd"
appname := "todoapi"

args ?= ""

goose:
	goose -env .env $(args) 

build:
	@go mod tidy
	@go build -o $(bin)/$(appname) ./$(entry)/main.go

run: build
	./$(bin)/$(appname)
