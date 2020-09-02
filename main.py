import os
import socket

from bloggulus.server import runserver


if __name__ == '__main__':
    if os.getenv('LISTEN_FDS'):
        s = socket.fromfd(3, socket.AF_INET, socket.SOCK_STREAM)
        s.setblocking(False)

        print('Serving HTTP on port:', 80)
        runserver(s)
    else:
        addr = ('', 5000)
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(addr)
        s.listen(128)
        s.setblocking(False)

        print('Serving HTTP on port:', 5000)
        runserver(s)
