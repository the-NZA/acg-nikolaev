APP = acg
BUILD_FLAGS = -v

.PHONY: build
build:
	go build $(BUILD_FLAGS) ./cmd/$(APP)

.PHONY: run
run:
	go run $(BUILD_FLAGS) ./cmd/$(APP)

.PHONY: buildnrace
buildnrace:
	go run $(BUILD_FLAGS) -race ./cmd/$(APP)

DBPORT = 27017
DBDIR $(shell pwd)/mongodata

.PHONY: rundb
rundb:
	docker run -it -p $(DBPORT):$(DBPORT) -v $(DBDIR):/data/db --name acg_db --rm mongo

.PHONY: clean
clean:
	rm ./$(APP)

.DEFAULT_GOAL := build