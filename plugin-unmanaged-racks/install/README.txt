This document will provide steps on how to create docker image & deploy through Helm chart.

- Clone the ur repository
- Run below commands to set envoirnment variables:
   $ export ODIMRA_GROUP_ID=<group_id>
   $ export ODIMRA_USER_ID=<user_id>

- Run below command to create docker image
   $ cd plugin-ur
   $ ./build_images.sh

- Run below command to deploy ur plugin using helm chart
   $ cd plugin-ur/install/Kubernetes/helmcharts
   $ ./deploying_chart.sh <namespace>

Note: Before deploying ur plugin we need to populate values(now manually by editing the required files) to be picked by helm chart before deployment.
Folowing files are required to be edited with the user configured values:
   - plugin-ur/install/Kubernetes/helmcharts/ur-platformconfig/values.yaml
   - plugin-ur/install/Kubernetes/helmcharts/urplugin-config/values.yaml
   - plugin-ur/install/Kubernetes/helmcharts/urplugin-pv-pvc/values.yaml
   - plugin-ur/install/Kubernetes/helmcharts/urplugin/values.yaml

- Run below command to undeploy ur plugin 
   $ cd plugin-ur/install/Kubernetes/helmcharts
   $ ./undeploy_chart.sh <namespace>
