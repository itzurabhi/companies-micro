HOST=0.0.0.0
PORT=8081

build:
	go build -race -o bin/ ./...
test:
	go test -coverprofile=coverage.out ./...
coverage: coverage.out
	go tool cover -html=coverage.out
docker: extern/wait-for-it/wait-for-it.sh
	docker build -t companies-micro .
test-create-company: data/test/company.create.json
	curl -XPOST -H "Content-type: application/json" -d \@data/test/company.create.json  '${HOST}:${PORT}/companies'
test-update-company: data/test/company.update.json
	curl -XPATCH -H "Content-type: application/json" -d \@data/test/company.update.json  '${HOST}:${PORT}/companies/3c2bd8d3-6f87-4558-a0e4-d4e2a9614736'
extern/wait-for-it/wait-for-it.sh:
	git submodule update --init --recursive
login:
	@curl -XPOST -s -H "Content-Type: application/x-www-form-urlencoded"  -d 'user=john&password=doe'  '${HOST}:${PORT}/login' | jq .token
genmock:
	docker run --rm -v "${PWD}":/src -w /src vektra/mockery --inpackage --keeptree --all
