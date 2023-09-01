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

FROM ubuntu:22.04

ARG ODIMRA_USER_ID
ARG ODIMRA_GROUP_ID

WORKDIR /ODIM/

RUN if [ -z "$ODIMRA_USER_ID" ] || [ -z "$ODIMRA_GROUP_ID" ]; then echo "\n[$(date)] -- ERROR -- ODIMRA_USER_ID or ODIMRA_GROUP_ID is not set\n"; exit 1; fi \
    && groupadd -r -g $ODIMRA_GROUP_ID odimra \
    && useradd -s /bin/bash -u $ODIMRA_USER_ID -m -d /home/odimra -r -g odimra odimra
    

RUN apt-get update -y && \
    apt-get -y install redis-tools python3 python3-venv python3-dev python3-pip binutils libc6 && \
    python3 -m pip install --upgrade pip && \
    python3 -m pip install grpcio grpcio-tools


COPY svc-composition-service /ODIM/svc-composition-service
COPY lib-utilities /ODIM/lib-utilities

COPY install/Docker/dockerfiles/scripts/build_cs.sh .
RUN chmod 755 build_cs.sh

RUN /bin/bash build_cs.sh
