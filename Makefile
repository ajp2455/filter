all:
	go build ./cmd/filter

install: all
	mv ./filter ${DESTDIR}/filter

clean:
	rm -f ./filter

.PHONY: all install clean
