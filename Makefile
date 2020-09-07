.POSIX:
.SUFFIXES:

.PHONY: deps
deps:
	python3 -m venv venv  &&  \
	. ./venv/bin/activate &&  \
	pip install -Uq wheel  &&  \
	pip install -Uq -r requirements.txt

.PHONY: run
run: deps
	. ./venv/bin/activate &&  \
	FLASK_APP=bloggulus.app   \
	FLASK_ENV=development     \
	flask run

.PHONY: dist
dist:
	mkdir -p dist/
	cp -r bloggulus/ dist/
	python3 -m pip install -Uq -r requirements.txt --target dist/
	zip -rq bloggulus.zip dist/

.PHONY: clean
clean:
	rm -fr bloggulus.zip dist/ __pycache__/ bloggulus/__pycache__/
