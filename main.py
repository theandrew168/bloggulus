import os
import socket

from waitress import serve

from bloggulus.app import app


if __name__ == '__main__':
    s = socket.create_server(('127.0.0.1', 8080), backlog=128, reuse_port=True)
    serve(app, sockets=[s], threads=os.cpu_count() * 4)
