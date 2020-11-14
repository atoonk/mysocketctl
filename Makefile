SHELL=/bin/bash

install: 
	-rm -r build dist *.egg-info
	python3 setup.py bdist_wheel sdist
	rm -rf */__pycache_
	twine upload --repository testpypi dist/*
