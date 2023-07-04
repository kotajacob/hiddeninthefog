# hiddeninthefog
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr
GO ?= go
GOFLAGS ?= -buildvcs=false
RM ?= rm -f

all: hiddeninthefog

hiddeninthefog:
	$(GO) build $(GOFLAGS) .

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f hiddeninthefog $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/hiddeninthefog

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/hiddeninthefog

clean:
	$(RM) hiddeninthefog

.PHONY: all hiddeninthefog install uninstall clean
