PATH := ${GOROOT}/bin:${GOPATH}/bin:${PATH}
SHELL := env PATH=$(PATH) /bin/bash
COLOR := "\e[1;36m%s\e[0m\n"


full-check: check test-module

define NEWLINE
endef

test-module:
	@printf $(COLOR) "Start all test..."
	- go test -covermode=count  ../...
	@printf $(COLOR) "Finish all test."

check:
	@printf $(COLOR) "Start staticcheck..."
	- staticcheck -f stylish ../...
	@printf $(COLOR) "Finish staticcheck."

