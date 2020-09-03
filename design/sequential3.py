import socketserver

HOST, PORT = '0.0.0.0', 8888
ADDR = (HOST, PORT)


class HTTPHandler(socketserver.BaseRequestHandler):
    allow_reuse_address = True
    request_queue_size = 128

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
    print('Serving HTTP on:', ADDR)
    with socketserver.TCPServer(ADDR, HTTPHandler) as server:
        server.serve_forever()
