import socketserver

ADDR = ('0.0.0.0', 8888)


class ThreadingTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    allow_reuse_address = True
    request_queue_size = 128


class HTTPHandler(socketserver.BaseRequestHandler):

    def handle(self):
        req = self.request.recv(1024)
        if not req:
            self.request.close()
            return

        req_lines = req.split(b'\r\n')
        method, path, version = req_lines[0].decode().split()
        print(method, path, version)

        resp = b'HTTP/1.1 200 OK\r\n\r\nHello, World!'
        self.request.send(resp)
        self.request.close()


if __name__ == '__main__':
    with ThreadingTCPServer(ADDR, HTTPHandler) as server:
        print('Serving HTTP on:', ADDR)
        server.serve_forever()
