TARGET := file-parser
CGO_ENABLED := 0
GOOS := linux
GOARCH := amd64
CMD_API_DIR := cmd/api
BIN_DIR := bin

GO_ENV_VAR := CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH)
FLAGS := -ldflags="-s -w"

all: $(TARGET)

$(TARGET):
	$(GO_ENV_VAR) go build $(FLAGS) -o $(BIN_DIR)/$(TARGET) $(CMD_API_DIR)/main.go

launch: clean all
	./$(BIN_DIR)/$(TARGET)

run:
	go run $(CMD_API_DIR)/main.go

build:
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down

clean:
	-rm $(BIN_DIR)/*
