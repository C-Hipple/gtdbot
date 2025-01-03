echo "Running go vet"
go vet

echo "\nRunning gofmt check"
gofmt -s -l .

echo "\nRunning go test"
go test ./...