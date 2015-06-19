# assumes GNU make

DIR_BIN=../bin

TARGETS = \
	$(DIR_BIN)/fetch_and_wait \
	$(DIR_BIN)/pick_config \
	$(DIR_BIN)/merge_config

SRC_COMMON = \
	$(wildcard util/*.go)

all: $(TARGETS)
	@echo "$@ done."

fmt: 
	@find . -name '*.go' | while read line; do dirname "$$line"; done | sort -u | while read line; do echo "go fmt in $$line" 1>&2 && (cd "$$line" && go fmt); done

clean:
	/bin/rm -f $(TARGETS)
	@echo "$@ done."

$(DIR_BIN)/fetch_and_wait: fetch_and_wait.go
	go build -o $@ $<

$(DIR_BIN)/pick_config: pick_config.go $(SRC_COMMON)
	go build -o $@ $<

$(DIR_BIN)/merge_config: merge_config.go $(SRC_COMMON)
	go build -o $@ $<

