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
	python main.py

.PHONY: build
build:
	mkdir -p build/
	cp main.py build/__main__.py
	cp -r bloggulus/ build/
	python3 -m pip install -U -r requirements.txt --target build
	python3 -m zipapp -c -p "/usr/bin/env python3" -o "bloggulus.pyz" build

.PHONY: dist
dist: build
	mkdir -p dist/
	cp -r bloggulus.pyz dist/bloggulus
	mkdir -p dist/web/
	cp -r static/* dist/web/

.PHONY: clean
clean:
	rm -fr bloggulus.pyz build/ dist/ __pycache__/
