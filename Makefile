PROG=imap-cleaner
DEPS=$(shell find -type f -name "*.go")

.PHONY: build clean

build: Makefile $(DEPS)
	CGO_ENABLED=0 GOAMD64=v3 go build -o $(PROG) -trimpath -buildmode=pie .
clean:
	$(RM) $(PROG)
