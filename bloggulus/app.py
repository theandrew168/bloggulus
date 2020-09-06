import os

from flask import Flask, render_template
import jinja2

DATABASE = os.getenv('BLOGGULUS_DATABASE') or 'bloggulus.sqlite3'
SECRET_KEY = os.getenv('BLOGGULUS_SECRET_KEY') or 'bloggulus_secret_key'

app = Flask(__name__)
app.config.from_object(__name__)
app.jinja_loader = jinja2.PackageLoader('bloggulus', 'templates')


@app.route('/')
def hello_world():
    return 'Hello, World!'


@app.route('/base')
def base():
    return render_template('base.html')
