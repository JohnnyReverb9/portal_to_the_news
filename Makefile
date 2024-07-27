BINARY_NAME=portal_to_the_news
OUTPUT_DIR=output

.PHONY: build clean

# ~ make build
build:
	mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) .

# ~ make clean
clean:
	rm -fr $(OUTPUT_DIR)/$(BINARY_NAME)