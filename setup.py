from setuptools import setup

setup(
    name='mysocketctl',
    version='0.1',
    py_modules=['mysocketctl'],
    install_requires=[
        'Click','requests','pyjwt','prettytable'
    ],
    entry_points='''
        [console_scripts]
        mysocketctl=mysocketctl:cli
    ''',
)
