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
	FLASK_APP=bloggulus.app   \
	FLASK_ENV=development     \
	flask run

.PHONY: build
build:
	mkdir -p build/
	cp -r bloggulus/ build/
	cp bloggulus/__main__.py build/
	python3 -m pip install -U -r requirements.txt --target build
	python3 -m zipapp -c -p "/usr/bin/env python3" -o "bloggulus.pyz" build

.PHONY: dist
dist: build
	mkdir -p dist/
	mv bloggulus.pyz dist/bloggulus
	cp -r bloggulus/static/ dist/

.PHONY: clean
clean:
	rm -fr bloggulus.pyz build/ dist/ __pycache__/ bloggulus/__pycache__/
