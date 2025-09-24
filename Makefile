.PHONY: build run

# ========= Vars definitions =========

app = fstmon
test = test

# ========= Prepare commands =========

tidy:
	go mod tidy
	go clean

del:
	rm ./$(app)* || echo "file didn't exists"
	rm ./trace*  || echo "file didn't exists"

# ========= Compile commands =========

build-test:
	GOOS=linux
	GOARCH=amd64
	go build -o ./tester -v ./cmd/$(test)/main.go

run-test: del build-test
	./tester

build:
	GOOS=linux
	GOARCH=amd64
	go build -o ./$(app) -v ./cmd/$(app)/main.go

run: del build
	./$(app)

build-prod:
	go build -ldflags="-s -w" -o ./$(app) -v ./cmd/$(app)/main.go

pack: build-prod
	rm -rf ./fstmon-app || echo ""
	mkdir ./fstmon-app
	cp init/fstmon.service ./fstmon-app/fstmon.service
	mv ./fstmon ./fstmon-app/fstmon


.DEFAULT_GOAL := run
