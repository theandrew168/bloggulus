import os
import sys

import django
from gunicorn.app import wsgiapp


def main():
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", "settings.development")
    django.setup()

    # setup args to gunicorn (since it doesn't have a python API)
    workers = os.cpu_count() * 2 + 1
    sys.argv = [
        '.',
        'bloggulus.wsgi',
        '--bind=127.0.0.1:8000',
        '--workers={}'.format(workers),
    ]

    wsgiapp.run()
