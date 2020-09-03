import os
import socket

from pyramid.config import Configurator
from waitress import serve

from bloggulus.handlers import hello_world


# os.environ['LISTEN_FDS'] will hold number of FDs
#
# fds = [fd + 3 for fd in range(LISTEN_FDS)]
# s80 = socket.fromfd(fds[0], socket.AF_INET, socket.SOCK_STREAM)
# s443 = socket.fromfd(fds[1], socket.AF_INET, socket.SOCK_STREAM)
#
# SYSTEMD_SOCKET_FD_PORT_80 = 3
# SYSTEMD_SOCKET_FD_PORT_443 = 4


if __name__ == '__main__':
    with Configurator() as config:
        config.add_route('hello', '/')
        config.add_view(hello_world, route_name='hello')
        app = config.make_wsgi_app()

    if os.getenv('LISTEN_FDS'):
        s = socket.fromfd(3, socket.AF_INET, socket.SOCK_STREAM)
        s.setblocking(False)

        # TODO: get TLS files from Let's Encrypt and ssl.wrap_socket()

        serve(app, sockets=[s], threads=os.cpu_count() * 4)
    else:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(('0.0.0.0', 8080))
        s.listen(128)
        s.setblocking(False)

        serve(app, sockets=[s], threads=os.cpu_count() * 4)
