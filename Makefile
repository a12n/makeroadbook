.PHONY: all clean

all: makeroadbook

clean:
	rm -f makeroadbook

makeroadbook: src/makeroadbook/main.go
	GOPATH=$(PWD) go build makeroadbook
