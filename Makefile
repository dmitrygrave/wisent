NAME=wisent
OUTPUT=output/$(NAME)
PUBLIC=public
PUBLICDART=$(PUBLIC)/dart
DARTINDEX=$(PUBLICDART)/$(NAME).dart
JSOUTPUT=$(PUBLIC)/js/$(NAME).js
PUBLICSASS=$(PUBLIC)/sass/
SASSINDEX=$(PUBLICSASS)$(NAME).scss
SASSOUTPUT=$(PUBLIC)/css/$(NAME).css

.PHONY: build clean test

build: go dart sass
test: test

go: $(OUTPUT)
dart: $(JSOUTPUT)
sass: $(SASSOUTPUT)

$(OUTPUT):
	go build

$(JSOUTPUT): $(DARTINDEX)
	dart2js --out=$@ $<

$(SASSOUTPUT): $(SASSINDEX)
	sass $<:$@

test:
    @echo "Tests will be put here"

clean:
	@echo "Cleaning up files..."
	go clean
	rm -rf $(PUBLIC)/css/*
	rm -rf $(PUBLIC)/js/*
