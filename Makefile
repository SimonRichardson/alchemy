PATH_DISCOURSE = github.com/SimonRichardson/alchemy

.PHONY: all
all: install
	$(MAKE) clean build

.PHONY: install
install:
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/mock/mockgen
	glide install --strip-vendor

.PHONY: build
build: dist/alchemy

.PHONY: clean
clean: FORCE
	rm -f dist/alchemy

dist/alchemy:
	go build -o dist/alchemy ${PATH_DISCOURSE}/cmd/alchemy

.PHONY: build-mocks
build-mocks:
	go generate ./pkg/...

.PHONY: clean-mocks
clean-mocks: FORCE
	@ find ./pkg -type d -name 'mocks' -exec find {} -name '*.go' -delete -print \;

.PHONY: FORCE
FORCE:


.PHONY: unit-tests
unit-tests:
	docker-compose run alchemy go test -v ./pkg/...

.PHONY: integration-tests
integration-tests:
	docker-compose run alchemy go test -v -tags=integration ./pkg/...

.PHONY: coverage-tests
coverage-tests:
	@ mkdir -p bin
	docker-compose run alchemy go test -covermode=count -coverprofile=bin/coverage.out -v -tags=integration ${COVER_PKG}

.PHONY: coverage-view
coverage-view:
	go tool cover -html=bin/coverage.out
