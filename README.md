# bloggulus
RSS aggregator powered by Python, Flask, and SQLite

## Building
Create and activate a [virtual environment](https://docs.python.org/3/library/venv.html):
```
python3 -m venv venv/
. ./venv/bin/activate
```

Download project depedencies:
```
pip install -r requirements.txt
```

Run the development server:
```
FLASK_ENV=development flask run
```
