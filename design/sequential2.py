import socket

HOST, PORT = '0.0.0.0', 8888
ADDR = (HOST, PORT)


def runserver(s):
    while True:
        c, _ = s.accept()

        req = c.recv(1024)
        if not req:
            c.close()
            continue

        req_lines = req.split(b'\r\n')
        method, path, version = req_lines[0].decode().split()
        print(method, path, version)

        resp = b'HTTP/1.1 200 OK\r\n\r\nHello, World!'
        c.send(resp)
        c.close()


if __name__ == '__main__':
    s = socket.create_server(ADDR, backlog=128, reuse_port=True)

    print('Serving HTTP on:', ADDR)
    runserver(s)
