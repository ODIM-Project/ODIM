//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package common ...
package common

import (
	"context"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	// K8sODIMNamespace has the namepace where service is deployed
	K8sODIMNamespace string
	// isK8sDeployment is for indicating if it is kubernetes deployment
	isK8sDeployment bool
	// updateConf is for updating all required configuration once
	// process lifetime
	updateConf sync.Once
)

// updateConfig is for updating all the required configurations
func updateConfig() {
	var exists bool
	K8sODIMNamespace, exists = os.LookupEnv("ODIM_NAMESPACE")
	if !exists {
		log.Info("ODIM_NAMESPACE environment variable not found, not a kubernetes deployment")
		return
	}
	if K8sODIMNamespace != "" {
		isK8sDeployment = true
	} else {
		log.Fatalf("value not set for ODIM_NAMESPACE environment variable")
	}
	return
}

// IsK8sDeployment is for finding out if it is kubernetes deployment
func IsK8sDeployment() bool {
	updateConf.Do(updateConfig)
	return isK8sDeployment
}

// GetServiceEndpointAddresses is for getting the list of IP addresses of all the pods
// belonging to the passed kubernetes service.
// Expects service name as defined in kubernetes deployment and returns the string slice
// with IP addresses of the pods and error if operation failed
// IsK8sDeployment should be used before calling GetServiceEndpointAddresses to check
// if it is kubernetes deployment
func GetServiceEndpointAddresses(srvName string) ([]string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to set k8s config: %s", err.Error())
	}

	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to set k8s clientset config: %s", err.Error())
	}

	endpoints, err := clientset.CoreV1().Endpoints(K8sODIMNamespace).Get(context.TODO(), srvName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s service endpoint info: %s", srvName, err.Error())
	}

	addrList := make([]string, 0)
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			addrList = append(addrList, address.IP)
		}
	}

	return addrList, nil
}
