from setuptools import setup

setup(
    name='mysocketctl',
    version='0.1',
    description="CLI tool for mysocket.io",
    long_description=open("README.rst").read(),
    py_modules=['mysocketctl'],
    install_requires=[
        'Click','requests','pyjwt','prettytable'
    ],
    classifiers=[
        "Intended Audience :: Developers",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: MacOS :: MacOS X",
        "Operating System :: POSIX",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3 :: Only",
        "Programming Language :: Python :: 3.5",
        "Programming Language :: Python :: 3.6",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: Implementation :: CPython",
        "Programming Language :: Python :: Implementation :: PyPy",
        "Topic :: Software Development :: Libraries :: Python Modules"
    ]
    entry_points='''
        [console_scripts]
        mysocketctl=mysocketctl:cli
    ''',
)
