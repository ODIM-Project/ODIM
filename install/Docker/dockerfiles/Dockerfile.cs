FROM ubuntu:18.04

ARG ODIMRA_USER_ID
ARG ODIMRA_GROUP_ID

WORKDIR /ODIM/

RUN if [ -z "$ODIMRA_USER_ID" ] || [ -z "$ODIMRA_GROUP_ID" ]; then echo "\n[$(date)] -- ERROR -- ODIMRA_USER_ID or ODIMRA_GROUP_ID is not set\n"; exit 1; fi \
    && groupadd -r -g $ODIMRA_GROUP_ID odimra \
    && useradd -s /bin/bash -u $ODIMRA_USER_ID -m -d /home/odimra -r -g odimra odimra \
    && mkdir -p /etc/csplugin_config && chown odimra:odimra /etc/csplugin_config

RUN apt update -y && apt-get install redis-tools -y && apt-get install python3 -y && apt-get install python3-venv -y && apt-get install python3-dev -y
RUN apt-get install -y binutils libc6

COPY svc-composition-service /ODIM/svc-composition-service

COPY install/Docker/dockerfiles/scripts/build_cs.sh .
RUN chmod 755 build_cs.sh

RUN /bin/bash build_cs.sh