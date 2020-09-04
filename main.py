import os
import socket

from waitress import serve

from bloggulus.wsgiapp import Application

# os.environ['LISTEN_FDS'] will hold number of FDs
#
# fds = [fd + 3 for fd in range(LISTEN_FDS)]
# s80 = socket.fromfd(fds[0], socket.AF_INET, socket.SOCK_STREAM)
# s443 = socket.fromfd(fds[1], socket.AF_INET, socket.SOCK_STREAM)
#
# SYSTEMD_SOCKET_FD_PORT_80 = 3
# SYSTEMD_SOCKET_FD_PORT_443 = 4


if __name__ == '__main__':
    threads = os.cpu_count() * 4

    if os.getenv('LISTEN_FDS'):
        s = socket.fromfd(3, socket.AF_INET, socket.SOCK_STREAM)
        s.setblocking(False)

        # TODO: kick off 80 -> 443 redirect thread
        # TODO: get TLS files from Let's Encrypt and ssl.wrap_socket()

        web_root = os.environ['BLOGGULUS_WEB_ROOT']

        app = Application(web_root, templates_root)
        serve(app, sockets=[s], threads=threads)
    else:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(('0.0.0.0', 8080))
        s.listen(128)
        s.setblocking(False)

        app = Application('./web')
        serve(app, sockets=[s], threads=threads, expose_tracebacks=True)
