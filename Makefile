.PHONY: all clean rel

SRCS = dist.go gpx.go main.go roadbook.go

all: makeroadbook

clean:
	rm -f makeroadbook

rel: makeroadbook-linix-386

makeroadbook: $(SRCS)
	go build -o $@ $^

makeroadbook-linix-386: makeroadbook
	cp $< $@
	strip $@
