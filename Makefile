NAME=wisent
OUTPUT=output/$(NAME)
PUBLIC=public
PUBLICDART=$(PUBLIC)/dart
DARTINDEX=$(PUBLICDART)/$(NAME).dart
JSOUTPUT=$(PUBLIC)/js/$(NAME).js
PUBLICSASS=$(PUBLIC)/sass/
SASSINDEX=$(PUBLICSASS)$(NAME).scss
SASSOUTPUT=$(PUBLIC)/css/$(NAME).css

.PHONY: all clean

all: build dart sass

build: $(OUTPUT)
dart: $(JSOUTPUT)
sass: $(SASSOUTPUT)

$(OUTPUT):
	go build

$(JSOUTPUT): $(DARTINDEX)
	dart2js --out=$@ $<

$(SASSOUTPUT): $(SASSINDEX)
	sass $<:$@

clean:
	@echo "Cleaning up files..."
	go clean
	rm -rf $(PUBLIC)/css/*
	rm -rf $(PUBLIC)/js/*
