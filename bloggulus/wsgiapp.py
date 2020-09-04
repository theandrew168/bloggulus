from http import HTTPStatus
from pprint import pformat


def status_string(status):
    return '{} {}'.format(status.value, status.phrase)


class Application:

    def __init__(self, web_root, templates_root):
        self.web_root = web_root
        self.templates_root = templates_root

    def __call__(self, environ, start_response):
        response_headers = [('Content-Type', 'text/plain')]
        start_response(status_string(HTTPStatus.OK), response_headers)

        resp = []

        # TODO: can make this a dict-based class with a master regex pat
        # that'd even work for args! /feed/[0-9]+
        if environ['PATH_INFO'] == '/':
            resp.append(b'index\n')
        elif environ['PATH_INFO'] == '/foo':
            resp.append(b'foo\n')
        elif environ['PATH_INFO'] == '/bar':
            resp.append(b'bar\n')

        # TODO: check static files and serve (adjust content-type)
        # TODO: else 404

#        env = pformat(environ)
#        resp.append(env.encode())

        return resp
