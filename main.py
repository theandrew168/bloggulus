import http.server
import os
import socket
import socketserver
import ssl
import subprocess
import threading

from waitress import serve

from bloggulus.app import Application

# os.environ['LISTEN_FDS'] will hold number of FDs
#
# fds = [fd + 3 for fd in range(LISTEN_FDS)]
# s80 = socket.fromfd(fds[0], socket.AF_INET, socket.SOCK_STREAM)
# s443 = socket.fromfd(fds[1], socket.AF_INET, socket.SOCK_STREAM)
#
# SYSTEMD_SOCKET_FD_PORT_80 = 3
# SYSTEMD_SOCKET_FD_PORT_443 = 4


class FileServerHandler(http.server.SimpleHTTPRequestHandler):
    pass


class FileServer:

    def __init__(self, sock, directory):
        self.sock = sock
        self.directory = directory
        self.thread = None
        self.httpd = None

    def serve(self):
        handler = FileServerHandler
        handler.directory = self.directory
        self.httpd = socketserver.TCPServer(None, handler, bind_and_activate=False)
        self.httpd.socket = self.sock
        self.httpd.serve_forever()

    def __enter__(self):
        self.thread = threading.Thread(target=self.serve)
        self.thread.start()

    def __exit__(self, *args):
        self.httpd.shutdown()
        self.thread.join()


# TODO: this guy is rough
def require_https(sock):
    while True:
        try:
            client, _ = sock.accept()
            req = client.recv(1024)
            if not req:
                client.close()
                continue

            req_lines = req.split(b'\r\n')
            method, path, version = req_lines[0].decode().split()

            resp = '\r\n'.join([
                'HTTP/1.1 301 Moved Permanently',
                'Location: https://bloggulus.com{}'.format(path),
            ])
            resp += '\r\n\r\n'

            client.sendall(resp.encode())
            client.close()
        except:
            pass


if __name__ == '__main__':
    threads = os.cpu_count() * 4

    if os.getenv('LISTEN_FDS'):
        s80 = socket.fromfd(3, socket.AF_INET, socket.SOCK_STREAM)
        s80.setblocking(False)

        s443 = socket.fromfd(4, socket.AF_INET, socket.SOCK_STREAM)
        s443.setblocking(False)

        home = os.environ['HOME']
        ssl_dir = os.path.join(home, 'live', 'www.bloggulus.com')
        web_root = os.environ['BLOGGULUS_WEB_ROOT']

        # TODO handle renewals either here (thread) or elsewhere
        # expose web root on port 80 and get TLS files from Let's Encrypt
        with FileServer(s80, web_root):
            certbot_cmd = [
                'certbot',
                'certonly',
                '--agree-tos',
                '--keep-until-expiring',
                '--register-unsafely-without-email',
                '--config-dir', home,
                '--work-dir', home,
                '--logs-dir', home,
                '-d', 'www.bloggulus.com',
                '-d', 'bloggulus.com',
                '--webroot',
                '-w', web_root,
            ]
            rc = subprocess.run(certbot_cmd)
            print(rc)

        # start HTTP to HTTPS redirect server
        thread = threading.Thread(target=require_https, args=(s80,), daemon=True)
        thread.start()

        # call ssl.wrap_socket() on s443
        context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        context.load_cert_chain(
            os.path.join(ssl_dir, 'fullchain.pem'),
            os.path.join(ssl_dir, 'privkey.pem'))
        s443 = context.wrap_socket(s443, server_side=True)

        app = Application(web_root)
        serve(app, sockets=[s443], threads=threads, url_scheme='https')
    else:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(('0.0.0.0', 8080))
        s.listen(128)
        s.setblocking(False)

        thread = threading.Thread(target=require_https, args=(s,), daemon=True)
        thread.start()

        app = Application('./web')
        serve(app, sockets=[s], threads=threads, expose_tracebacks=True)
