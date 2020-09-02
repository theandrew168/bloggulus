import socketserver

HOST, PORT = '0.0.0.0', 8888
ADDR = (HOST, PORT)


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
    print('Serving HTTP on:', ADDR)
    with ThreadingTCPServer(ADDR, HTTPHandler) as server:
        server.serve_forever()
