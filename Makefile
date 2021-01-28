test: #run unit test
	go test -race $$(go list ./... | grep -v /vendor/) -coverprofile coverage.out 
	go tool cover -html=coverage.out -o coverage.html

vendor : #get dependency
	go mod tidy
	go mod vendor
run : 
	go run cmd/e-wallet/*.go