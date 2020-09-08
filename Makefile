.POSIX:
.SUFFIXES:

# References:
# https://shiv.readthedocs.io/en/latest/django.html
# https://lincolnloop.com/blog/single-file-python-django-deployments/
# https://www.youtube.com/watch?v=Jzf8gTLN1To
# https://www.youtube.com/watch?v=n2q1SxL4-mY

.PHONY: deps
deps:
	python3 -m venv venv/
	./venv/bin/pip install -Uq wheel
	./venv/bin/pip install -Uq shiv
	./venv/bin/pip install -Uq -r requirements.txt

.PHONY: check
check: deps
	./venv/bin/python manage.py test

.PHONY: build
build: deps
	rm -fr build/
	mkdir build/
	./venv/bin/pip install -r requirements.txt --target build/
	cp -r bloggulus/ build/
	cp manage.py build/
	./venv/bin/shiv            \
	--compressed               \
	--site-packages build/     \
	-p '/usr/bin/env python3'  \
	-e manage.main             \
	-o bloggulus.pyz

.PHONY: static
static: deps
	./venv/bin/python manage.py collectstatic --no-input

.PHONY: dist
dist: build static
	rm -fr dist/
	mkdir dist/
	cp bloggulus.pyz dist/
	cp -r static/ dist/

.PHONY: clean
clean:
	rm -fr bloggulus.pyz build/ dist/ static/
