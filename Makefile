.POSIX:
.SUFFIXES:

.PHONY: venv
venv:
	python3 -m venv venv
	. ./venv/bin/activate

.PHONY: deps
deps: venv
	pip install -U wheel
	pip install -U -r requirements.txt

.PHONY: run
run: deps
	python server.py

.PHONY: build
build: deps
	mkdir -p build/
	cp server.py build/__main__.py
	python -m pip install -U -r requirements.txt --target build
	python -m zipapp -c -p "/usr/bin/env python3" -o "bloggulus.pyz" build

.PHONY: dist
dist: build
	mkdir -p dist/
	cp -r bloggulus.pyz dist/bloggulus
	mkdir -p dist/web/
	cp -r static/* dist/web/

.PHONY: clean
clean:
	rm -fr bloggulus.pyz build/ dist/ __pycache__/
