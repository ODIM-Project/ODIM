#!/usr/bin/env bash
#(C) Copyright [2022] American Megatrends International LLC
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.

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
