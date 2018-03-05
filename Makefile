PATH_DISCOURSE = github.com/SimonRichardson/discourse

.PHONY: all
all: install
	$(MAKE) clean build

.PHONY: install
install:
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/mock/mockgen
	glide install --strip-vendor

.PHONY: build
build: dist/discourse

.PHONY: clean
clean: FORCE
	rm -f dist/discourse

dist/discourse:
	go build -o dist/discourse ${PATH_DISCOURSE}/cmd/discourse

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
	docker-compose run discourse go test -v ./pkg/...

.PHONY: integration-tests
integration-tests:
	docker-compose run discourse go test -v -tags=integration ./pkg/...

.PHONY: coverage-tests
coverage-tests:
	@ mkdir -p bin
	docker-compose run discourse go test -covermode=count -coverprofile=bin/coverage.out -v -tags=integration ${COVER_PKG}

.PHONY: coverage-view
coverage-view:
	go tool cover -html=bin/coverage.out
