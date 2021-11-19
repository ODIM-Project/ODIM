FROM python:3.7

WORKDIR /ODIM/

COPY svc-composition-service /ODIM/svc-composition-service

COPY install/Docker/dockerfiles/scripts/build_cs.sh .
RUN chmod 755 build_cs.sh

RUN /bin/bash build_cs.sh