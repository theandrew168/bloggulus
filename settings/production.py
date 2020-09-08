import os

from .default import *


DEBUG = False
ALLOWED_HOSTS = ['bloggulus.com']
SECRET_KEY = os.environ['BLOGGULUS_SECRET_KEY']
STATIC_ROOT = BASE_DIR / 'bloggulus' / 'static'
