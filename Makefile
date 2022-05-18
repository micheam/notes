TARGET = bin/notes

$(TARGET) : clean unit-test ./cmd/localserver
	go build -o $(TARGET) ./cmd/localserver   

.PHONY : clean
clean :
	@rm -rf bin/

.PHONY : unit-test              
unit-test :
	@go test ./...
