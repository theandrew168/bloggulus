from setuptools import find_packages, setup

with open('README.md') as f:
    readme = f.read()

setup(
    name='bloggulus',
    version='0.0.1',
    author='Andrew Dailey',
    author_email='info@shallowbrooksoftware.com',
    description='RSS feed aggregator powered by Django',
    long_description=readme,
    long_description_content_type='text/markdown',
    url='https://github.com/theandrew168/bloggulus',
    packages=find_packages(),
    include_package_data=True,
    python_requires='>=3.0',
    classifiers=[
        'Environment :: Web Environment',
        'Framework :: Django',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3 :: Only',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Topic :: Internet :: WWW/HTTP',
        'Topic :: Internet :: WWW/HTTP :: Dynamic Content',
    ],
)
