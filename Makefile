all:
	go build ./cmd/filter

install:
	install -D ./filter ${DESTDIR}/usr/bin/filter

clean:
	rm -f ./filter

.PHONY: all install clean
