.PHONY: build-cli build-note

build-note:
	bazel build //note

build-cli: build-note
	bazel build //cli

install: build-cli
	cp bazel-bin/cli/cli_/cli ~/go/bin/note

clean:
	bazel clean
