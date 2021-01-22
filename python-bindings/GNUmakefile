SPEC=py-terraform-provider-b2.spec

default: build

.PHONY: deps format testacc clean build

deps:
	pip install -r requirements.txt

format:
	black -Sl 100 b2_terraform

testacc: build

clean:
	rm -rf build dist py-terraform-provider-b2.egg-info

build:
	pyinstaller --log-level WARN --onefile ${SPEC}

install: build

all: deps build