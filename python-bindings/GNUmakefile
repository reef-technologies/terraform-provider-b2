NAME=py-terraform-provider-b2
DIR=b2_terraform
EGG_INFO=${NAME}.egg-info
SPEC=${NAME}.spec

default: build

.PHONY: deps format lint testacc clean build

deps:
	@pip install -r requirements-dev.txt

format:
	@black -Sl 100 ${DIR}

lint:
	@black --check -Sl 100 ${DIR}
	@flake8 --ignore=E501 ${DIR}
	@python ../scripts/check-headers.py '**/*.py'

testacc: build

clean:
	@rm -rf build dist ${EGG_INFO}

build:
	@pyinstaller --log-level WARN --onefile ${SPEC}
ifneq ($(origin CI), undefined)
	@echo "::set-output name=version::$(subst refs/tags/v,,${GITHUB_REF})"
endif

install: build

docs: build

all: deps lint build
