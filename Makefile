all : install

clean :
	@echo ">>> Cleaning and initializing gogeo project <<<"
	@go clean
	@gofmt -w .
	@go get github.com/stretchr/testify

test : clean
	@echo ">>> Running unit tests <<<"
	@go test ./ ./geohash

test-coverage : clean
	@echo ">>> Running unit tests and calculating code coverage <<<"
	@go test ./ ./geohash -cover

install : test
	@echo ">>> Building and installing gogeo <<<"
	@go install
	@echo ">>> gogeo installed successfully! <<<"
	@echo ""
