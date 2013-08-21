.PHONY: all clean distclean rel

GOHOSTARCH = $(shell go env GOHOSTARCH)
GOHOSTOS = $(shell go env GOHOSTOS)

BIN = makeroadbook
REL = makeroadbook-$(GOHOSTOS)-$(GOHOSTARCH)
SRCS = dist.go gpx.go main.go roadbook.go

all: $(BIN)

clean:
	rm -f $(BIN)

distclean: clean
	rm -f $(REL)

rel: $(REL)

$(BIN): $(SRCS)
	go build -o $@ $^

$(REL): $(BIN)
	cp $< $@
	strip $@
