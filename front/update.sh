#!/bin/bash

git pull

npm run build:prod

cd dist

ln -s ../../doc/_book ./doc
