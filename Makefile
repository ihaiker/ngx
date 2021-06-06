.PHONY: mkdocs godep test
mkdocs:
	docker-compose -f ./scripts/docker-compose.yml run --rm mkdocs

godep:
	go mod download

test: godep
	go test ./... --cover
