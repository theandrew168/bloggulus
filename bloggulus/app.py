from http import HTTPStatus
from pprint import pprint

from jinja2 import Environment, PackageLoader, select_autoescape


def status_string(status):
    return '{} {}'.format(status.value, status.phrase)


class Application:

    def __init__(self, web_root, redirect_http=True):
        self.web_root = web_root
        self.redirect_http = redirect_http

        self.templates = Environment(
            loader=PackageLoader('bloggulus', 'templates'),
            autoescape=select_autoescape(['html', 'jinja2']))

    # satisfies the WSGI protocol (PEP 3333)
    def __call__(self, environ, start_response):
        headers = []
        resp = []

        # redirect http to https (80 to 443)
        if self.redirect_http and environ['wsgi.url_scheme'] == 'http':
            target = 'https://' + environ['HTTP_HOST'] + environ['SCRIPT_NAME'] + environ['PATH_INFO']
            headers.append(('Location', target))
            start_response(status_string(HTTPStatus.MOVED_PERMANENTLY), headers)
            return resp

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
