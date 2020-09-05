import os
import socket

from waitress import serve

from bloggulus.app import Application


if __name__ == '__main__':
    threads = os.cpu_count() * 4

    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(('0.0.0.0', 8080))
    s.listen(128)
    s.setblocking(False)

    app = Application('./web')
    serve(app, sockets=[s], threads=threads)
