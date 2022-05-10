TARGET = bin/notes

$(TARGET) : clean ./cmd/localserver
	go build -o $(TARGET) ./cmd/localserver   

clean :
	rm -rf bin/
