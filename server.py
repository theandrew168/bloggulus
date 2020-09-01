import os
import socket
import sys

# os.environ['LISTEN_FDS'] will hold number of FDs
#
# fds = [fd + 3 for fd in range(LISTEN_FDS)]
# s80 = socket.fromfd(fds[0], socket.AF_INET, socket.SOCK_STREAM)
# s443 = socket.fromfd(fds[1], socket.AF_INET, socket.SOCK_STREAM)
#
# SYSTEMD_SOCKET_FD_PORT_80 = 3
# SYSTEMD_SOCKET_FD_PORT_443 = 4


def main_dev():
    addr = ('', 5000)
    s = socket.create_server(addr, backlog=128, reuse_port=True)
    print('Serving HTTP on port:', 5000)

    while True:
        c, _ = s.accept()
        req = c.recv(1024)
        print(req.decode())

        resp = b'HTTP/1.1 200 OK\r\n\r\nHello, World!'
        c.sendall(resp)
        c.close()


def main_prod():
    s = socket.fromfd(3, socket.AF_INET, socket.SOCK_STREAM)
    print('Serving HTTP on port:', 80)
    print('Serving HTTPS on port:', 443)

    while True:
        c, _ = s.accept()
        req = c.recv(1024)
        print(req.decode())

        resp = b'HTTP/1.1 200 OK\r\n\r\nHello, World!'
        c.sendall(resp)
        c.close()


if __name__ == '__main__':
    if os.getenv('LISTEN_FDS'):
        main_prod()
    else:
        main_dev()
