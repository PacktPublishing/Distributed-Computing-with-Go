
IMAGE = goophr-librarian
GOPKGS = $(shell go list ./...)
GOOS = linux
GOARCH = amd64
CGO_ENABLED = 0
BUILD_PATH = ./build
BINARY_NAME = "exec-bin"
SRC_ROOT = "."


test: lint
	@echo "RUNNING TESTS\n"
	go get ./...
	go test --cover -v $(GOPKGS)

lint:
	@echo "COMMENCING LINT CHECKS."
	go fmt ./...
	golint ./...
	go vet ./...

run.docker: test build.linux
	docker build -t $(IMAGE) .
	docker run --rm -it --name=$(IMAGE) $(IMAGE)

build.linux:
	@echo "BUILDING LINUX BINARY"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -o $(BUILD_PATH)/lin-$(BINARY_NAME) $(SRC_ROOT)

build.bin:
	@echo "BUILDING BINARY"
	go build -o $(BUILD_PATH)/$(BINARY_NAME) $(SRC_ROOT)

run: build.bin
	./$(BUILD_PATH)/$(BINARY_NAME)

