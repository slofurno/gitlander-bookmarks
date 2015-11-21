.PHONY: build

bookmarks: $(wildcard *.go)
	go build -o bookmarks

run: bookmarks
	./bookmarks
