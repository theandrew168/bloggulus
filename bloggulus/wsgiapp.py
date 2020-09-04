from http import HTTPStatus
from pprint import pformat

from jinja2 import Environment, PackageLoader, select_autoescape


def status_string(status):
    return '{} {}'.format(status.value, status.phrase)


class Application:

    def __init__(self, web_root):
        self.web_root = web_root

        self.templates = Environment(
            loader=PackageLoader('bloggulus', 'templates'),
            autoescape=select_autoescape(['html', 'jinja2']))

    def __call__(self, environ, start_response):
        headers = []
        resp = []

        # TODO: can make this a dict-based class with a master regex pat
        # that'd even work for args! /feed/[0-9]+
        if environ['PATH_INFO'] == '/':
            headers.append(('Content-Type', 'text/plain'))
            resp.append(b'index\n')
        elif environ['PATH_INFO'] == '/foo':
            headers.append(('Content-Type', 'text/plain'))
            resp.append(b'foo\n')
        elif environ['PATH_INFO'] == '/bar':
            headers.append(('Content-Type', 'text/plain'))
            resp.append(b'bar\n')
        elif environ['PATH_INFO'] == '/template':
            headers.append(('Content-Type', 'text/html'))
            t = self.templates.get_template('page.jinja2')
            r = t.render(name='mad jinja2 action!')
            resp.append(r.encode())

        # TODO: check static files and serve (adjust content-type)
        # TODO: else 404

#        env = pformat(environ)
#        resp.append(env.encode())

        start_response(status_string(HTTPStatus.OK), headers)
        return resp
