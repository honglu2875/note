.PHONY: build-cli build-note gazelle

gazelle:
	bazel run //:gazelle

build-note:
	bazel build //note --verbose_failures

build-cli: build-note
	bazel build //cli --verbose_failures

install: build-cli
	cp bazel-bin/cli/cli_/cli ~/go/bin/note
	chmod 755 ~/go/bin/note

clean:
	bazel clean
