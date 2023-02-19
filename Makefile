HOST=0.0.0.0
PORT=8081

build:
	go build -race -o bin/ ./...
test:
	go test ./...
docker:
	docker build -t companies-micro .
test-create-company: data/test/company.create.json
	curl -XPOST -H "Content-type: application/json" -d \@data/test/company.create.json  '${HOST}:${PORT}/companies'
test-update-company: data/test/company.update.json
	curl -XPATCH -H "Content-type: application/json" -d \@data/test/company.create.json  '${HOST}:${PORT}/companies/f250e314-fbc0-4ef8-80cf-902ff8d27e50'
