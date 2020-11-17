export DOCKER_BUILDKIT=1

build_args="--build-arg ODIMRA_USER_ID=2021 \
	    --build-arg ODIMRA_GROUP_ID=2021 \
	    --build-arg http_proxy=http://web-proxy.corp.hpecorp.net:8080/ \
	    --build-arg https_proxy=http://web-proxy.corp.hpecorp.net:8080/"

/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.accountSession -t account-session:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.aggregation -t aggregation:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.api -t api:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.consul -t consul:1.6 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.events -t events:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.fabrics -t fabrics:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.grfplugin -t grf-plugin:1.0 $build_args .
#/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.iloplugin -t iloplugin:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.kafka -t odim_kafka:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.managers -t managers:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.redis -t redis:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.systems -t systems:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.task -t task:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.update -t update:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.zookeeper -t odim_zookeeper:1.0 $build_args .
