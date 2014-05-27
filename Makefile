
all:
	scripts/build.sh

dist:
	scripts/dist.sh

clean:
	rm bin/sphere-status-leds || true
	rm -rf .gopath || true

test:
	go test ./...

.PHONY: all	dist clean test
