import selectors
import socket

HOST, PORT = '0.0.0.0', 8888
ADDR = (HOST, PORT)


def handle_client(mux, sock, mask):
    req = sock.recv(1024)
    if not req:
        mux.unregister(sock)
        sock.close()
        return

    req_lines = req.split(b'\r\n')
    method, path, version = req_lines[0].decode().split()
    print(method, path, version)

    resp = b'HTTP/1.1 200 OK\r\n\r\nHello, World!'
    sock.send(resp)
    mux.unregister(sock)
    sock.close()


def accept_client(mux, sock, mask):
    c, _ = s.accept()
    c.setblocking(False)
    mux.register(c, selectors.EVENT_READ, handle_client)


def runserver(s):
    mux = selectors.DefaultSelector()
    mux.register(s, selectors.EVENT_READ, accept_client)

    while True:
        events = mux.select()
        for key, mask in events:
            handler = key.data
            handler(mux, key.fileobj, mask)


if __name__ == '__main__':
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(ADDR)
    s.listen(128)
    s.setblocking(False)

    print('Serving HTTP on:', ADDR)
    runserver(s)
