from datetime import datetime, timezone
import io
import os
import socket
import sys

from bloggulus.app import app

# References:
# https://www.python.org/dev/peps/pep-3333/

# os.environ['LISTEN_FDS'] will hold number of FDs
#
# fds = [fd + 3 for fd in range(LISTEN_FDS)]
# s80 = socket.fromfd(fds[0], socket.AF_INET, socket.SOCK_STREAM)
# s443 = socket.fromfd(fds[1], socket.AF_INET, socket.SOCK_STREAM)
#
# SYSTEMD_SOCKET_FD_PORT_80 = 3
# SYSTEMD_SOCKET_FD_PORT_443 = 4


def serve_wsgi_app(s, app):
    host, port = s.getsockname()[:2]
    server_name = socket.getfqdn(host)
    server_port = str(port)

    while True:
        c, _ = s.accept()

        req = c.recv(1024)
        print(req.decode())

        req_lines = req.split(b'\r\n')
        method, path, version = req_lines[0].decode().split()

        env = {}

        # Required WSGI vars
        env['wsgi.version']      = (1, 0)
        env['wsgi.url_scheme']   = 'http'
        env['wsgi.input']        = io.BytesIO(req)
        env['wsgi.errors']       = sys.stderr
        env['wsgi.multithread']  = False
        env['wsgi.multiprocess'] = True
        env['wsgi.run_once']     = False

        # Required CGI vars
        env['REQUEST_METHOD']    = method
        env['PATH_INFO']         = path
        env['SERVER_NAME']       = server_name
        env['SERVER_PORT']       = server_port

        resp_status = ''
        resp_headers = [
            ('Date', datetime.now(timezone.utc).strftime('%a, %d %b %Y %H:%M:%S GMT')),
            ('Server', 'BloggulusWSGI 0.0.1'),
        ]

        def start_response(status, response_headers, exc_info=None):
            nonlocal resp_status
            resp_status = status
            resp_headers.extend(response_headers)

        result = app(env, start_response)

        resp_lines = []
        resp_lines.append('HTTP/1.1 {}'.format(resp_status).encode())
        for header in resp_headers:
            resp_lines.append('{}: {}'.format(*header).encode())

        resp = b'\r\n'.join(resp_lines)
        resp += b'\r\n\r\n'
        for data in result:
            resp += data

        print(resp.decode())

        c.sendall(resp)
        c.close()


def main_dev(app):
    addr = ('', 5000)
    s = socket.create_server(addr, backlog=128, reuse_port=True)
    print('Serving HTTP on port:', 5000)

    serve_wsgi_app(s, app)


def main_prod(app):
    s = socket.fromfd(3, socket.AF_INET, socket.SOCK_STREAM)
    print('Serving HTTP on port:', 80)
#    print('Serving HTTPS on port:', 443)

    serve_wsgi_app(s, app)


if __name__ == '__main__':
    if os.getenv('LISTEN_FDS'):
        main_prod(app.wsgi_app)
    else:
        main_dev(app.wsgi_app)
