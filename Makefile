all:
	go build ./cmd/filter

install: all
	mv ./filter ${DESTDIR}/usr/bin/filter

clean:
	rm -f ./filter

.PHONY: all install clean
