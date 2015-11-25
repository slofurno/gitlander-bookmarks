.PHONY: run

bookmarks: $(wildcard *.go)
	go build -o bookmarks

run: bookmarks
	./bookmarks
