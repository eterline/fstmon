.PHONY: build run

# ========= Vars definitions =========

app = fstmon

# ========= Prepare commands =========

tidy:
	go mod tidy
	go clean

del:
	rm ./$(app)* || echo "file didn't exists"
	rm ./trace*  || echo "file didn't exists"

# ========= Compile commands =========

build:
	GOOS=linux
	GOARCH=amd64
	go build -o ./$(app) -v ./cmd/$(app)/main.go

build-prod:
	go build -ldflags="-s -w" -o ./$(app) -v ./cmd/$(app)/main.go

run: del build
	./$(app)

.DEFAULT_GOAL := run
