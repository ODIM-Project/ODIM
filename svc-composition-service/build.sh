#!/usr/bin/env bash

# Install pre requisite tools

rm -rf ./env
python3 -m venv env
source env/bin/activate

cd env/bin
./python3 -m pip install --upgrade pip

./python3 pip3 install --no-cache-dir -r ../../requirements.txt

cd ../../

cd app
# Pre-build configuration

rm -rf ./dist
rm -rf ./build
rm -f ./svc-composition-service.spec

../env/bin/python3 ../env/bin/pyinstaller --onefile --name svc-composition-service main.py
cd ..
