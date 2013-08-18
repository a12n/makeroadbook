.PHONY: all clean rel

all: makeroadbook

clean:
	rm -f makeroadbook

rel: makeroadbook-linix-386

makeroadbook: src/makeroadbook/main.go
	GOPATH=$(PWD) go build makeroadbook

makeroadbook-linix-386: makeroadbook
	cp $< $@
	strip $@
