from datetime import datetime, timezone
import io
import multiprocessing as mp
import os
import selectors
import socket
import sys
import traceback

from .app import app

# References:
# https://www.python.org/dev/peps/pep-3333/
# https://ruslanspivak.com/lsbaws-part2/

# os.environ['LISTEN_FDS'] will hold number of FDs
#
# fds = [fd + 3 for fd in range(LISTEN_FDS)]
# s80 = socket.fromfd(fds[0], socket.AF_INET, socket.SOCK_STREAM)
# s443 = socket.fromfd(fds[1], socket.AF_INET, socket.SOCK_STREAM)
#
# SYSTEMD_SOCKET_FD_PORT_80 = 3
# SYSTEMD_SOCKET_FD_PORT_443 = 4

# Benchmarks (hey http://bloggulus.com):
# single process, serial - 4.66 RPS
# pre-forked (1 process), serial - 4.70 RPS
# pre-forked (2 process), serial -
# pre-forked (4 process), serial -
# pre-forked (1 process), async -
# pre-forked (2 process), async -
# pre-forked (4 process), async -


def handle_client(mux, sock, mask):
    req = sock.recv(1024)

    # TODO: why does this ever occur?
    if len(req) == 0:
        mux.unregister(sock)
        sock.close()
        return

#    print(req.decode())

    try:
        req_lines = req.split(b'\r\n')
        method, path, version = req_lines[0].decode().split()
    except:
        print(req)
        raise

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
    env['SERVER_NAME']       = 'bloggulus.com'
    env['SERVER_PORT']       = '80'

    resp_status = ''
    resp_headers = [
        ('Date', datetime.now(timezone.utc).strftime('%a, %d %b %Y %H:%M:%S GMT')),
        ('Server', 'BloggulusWSGI 0.0.1'),
    ]

    def start_response(status, response_headers, exc_info=None):
        nonlocal resp_status
        resp_status = status
        resp_headers.extend(response_headers)

    result = app.wsgi_app(env, start_response)

    resp_lines = []
    resp_lines.append('HTTP/1.1 {}'.format(resp_status).encode())
    for header in resp_headers:
        resp_lines.append('{}: {}'.format(*header).encode())

    resp = b'\r\n'.join(resp_lines)
    resp += b'\r\n\r\n'
    for data in result:
        resp += data

#    print(resp.decode())

    sock.sendall(resp)
    mux.unregister(sock)
    sock.close()


def accept_client(mux, sock, mask):
    try:
        c, _ = sock.accept()
    except BlockingIOError:
        # http://lkml.iu.edu/hypermail/linux/kernel/1508.0/01413.html
        #
        # New client connections will wake up all workers but
        # only one will actually get the accept(). The others will
        # all raise a blocking IO error (EGAIN under the hood?).
        # An Linux-specific optimization would be to set
        # EPOLLEXCLUSIVE on the underlying epoll object. This way,
        # only a single process would get woken up to accept()
        # the new client.
        return

    # the worker that _does_ accept() the client can move forward
    c.setblocking(False)
    mux.register(c, selectors.EVENT_READ, handle_client)


def wsgi_worker(s):
    mux = selectors.DefaultSelector()
    mux.register(s, selectors.EVENT_READ, accept_client)

    try:
        while True:
            events = mux.select()
            for key, mask in events:
                handler = key.data

                try:
                    handler(mux, key.fileobj, mask)
                except:
                    # don't crash the worker if something goes wrong
                    traceback.print_exc()
                    mux.unregister(key.fileobj)
                    key.fileobj.close()
    except KeyboardInterrupt:
        pass


def runserver(s):
    mp.set_start_method('spawn')

    workers = []
    for _ in range(os.cpu_count()):
        worker = mp.Process(target=wsgi_worker, args=(s,))
        worker.start()
        workers.append(worker)
        print('started worker {}'.format(worker.pid))

    for worker in workers:
        try:
            worker.join()
        except KeyboardInterrupt:
            pass
        print('stopped worker {}'.format(worker.pid))
