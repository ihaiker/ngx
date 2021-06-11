.PHONY: mkdocs godep test
mkdocs:
	docker-compose -f ./scripts/docker-compose.yml run --rm mkdocs build -c

godep:
	go mod download

test: godep
	go test ./... --cover

build: test
	sh -c ./build.sh

clean:
	@rm -rf bin
