all:
	go build ./cmd/filter

install: all
	mv ./filter /usr/bin/filter

clean:
	rm -f ./filter

.PHONY: all install clean
