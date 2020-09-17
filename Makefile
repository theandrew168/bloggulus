.POSIX:
.SUFFIXES:

.PHONY: default
default: run

.PHONY: deps
deps:
	python3 -m venv venv/
	./venv/bin/pip install -Uq wheel
	./venv/bin/pip install -Uq -r requirements.txt

.PHONY: run
run: deps
	FLASK_ENV=development ./venv/bin/flask run

.PHONY: check
check: deps
	true

.PHONY: build
build: deps
	rm -fr build/ && mkdir build/
	./venv/bin/pip install -q -r requirements.txt --target build/
	cp app.py build/
	python3 -m zipapp                  \
	  --compress                       \
	  --main "app:main"                \
	  --output "bloggulus.pyz"         \
	  --python "/usr/bin/env python3"  \
	  build/

.PHONY: dist
dist: build
	rm -fr dist/ && mkdir dist/
	cp bloggulus.pyz dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: clean
clean:
	rm -fr bloggulus.pyz build/ dist/
