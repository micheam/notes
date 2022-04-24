TARGET = bin/notes

$(TARGET) : clean ./cmd/cli/notes
	go build -o $(TARGET) ./cmd/cli/notes

clean :
	rm -rf bin/
