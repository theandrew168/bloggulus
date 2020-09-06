import os
import socket

from waitress import serve

from bloggulus.app import app


if __name__ == '__main__':
    serve(app, host='127.0.0.1', port=5000, threads=os.cpu_count() * 4)
