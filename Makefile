.POSIX:
.SUFFIXES:

.PHONY: deps
deps:
	python3 -m venv venv  &&  \
	. ./venv/bin/activate &&  \
	pip install -U wheel  &&  \
	pip install -U -r requirements.txt

.PHONY: run
run: deps
	. ./venv/bin/activate &&  \
	PYRAMID_RELOAD_ASSETS=1   \
	python dev.py

.PHONY: build
build:
	mkdir -p build/
	cp prod.py build/__main__.py
	cp -r bloggulus/ build/
	python3 -m pip install -U -r requirements.txt --target build
	python3 -m zipapp -c -p "/usr/bin/env python3" -o "bloggulus.pyz" build

.PHONY: dist
dist: build
	mkdir -p dist/
	mv bloggulus.pyz dist/bloggulus
	cp -r web/ dist/

.PHONY: clean
clean:
	rm -fr bloggulus.pyz build/ dist/ __pycache__/ bloggulus/__pycache__/
