.POSIX:
.SUFFIXES:

.PHONY: deps
deps:
	./venv/bin/pip install -Uq wheel
	./venv/bin/pip install -Uq shiv
	./venv/bin/pip install -Uq -r requirements.txt

.PHONY: static
static: deps
	./venv/bin/python manage.py collectstatic --no-input

.PHONY: dist
dist: deps static
	./venv/bin/shiv            \
	--compressed               \
	-p '/usr/bin/env python3'  \
	-o bloggulus.pyz           \
	-e bloggulus.main:main     \
	. -r requirements.txt

.PHONY: clean
clean:
	rm -fr bloggulus.pyz bloggulus/static/
