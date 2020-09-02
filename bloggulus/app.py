from flask import Flask

DATABASE = 'bloggulus.db'
SECRET_KEY = 'da80d8393ca7bbaee75ffa7c4c09067a437485ccf421b1239ad7d7b52b8e6d56'

app = Flask(__name__)
app.config.from_object(__name__)

@app.route("/")
def hello():
    return "Hello world from a Flask WSGI app!"

if __name__ == '__main__':
    app.run()
