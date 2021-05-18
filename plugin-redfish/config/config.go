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

//Package config ...
package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"

	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
)

// Data will have the configuration data from config file
var Data configModel

// configModel is for holding all the run time configurations for the svc-redfish-plugin
type configModel struct {
	FirmwareVersion         string            `json:"FirmwareVersion"` //FirmwareVersion of plugin of the plugin
	RootServiceUUID         string            `json:"RootServiceUUID"`
	SessionTimeoutInMinutes float64           `json:"SessionTimeoutInMinutes"` //plugin token time out in minutes
	PluginConf              *PluginConf       `json:"PluginConf"`
	LoadBalancerConf        *LoadBalancerConf `json:"LoadBalancerConf"`
	EventConf               *EventConf        `json:"EventConf"`
	MessageBusConf          *MessageBusConf   `json:"MessageBusConf"`
	KeyCertConf             *KeyCertConf      `json:"KeyCertConf"`
	URLTranslation          *URLTranslation   `json:"URLTranslation"`
	TLSConf                 *TLSConf          `json:"TLSConf"`
}

//PluginConf is for holding all the plugin related configurations
type PluginConf struct {
	ID       string `json:"ID"` // PluginID hold the id of the plugin
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

//LoadBalancerConf is for holding all load balancer related configurations
type LoadBalancerConf struct {
	Host string `json:"LBHost"`
	Port string `json:"LBPort"`
}

//EventConf is for holding all events related configuration
type EventConf struct {
	DestURI      string `json:"DestinationURI"`
	ListenerHost string `json:"ListenerHost"`
	ListenerPort string `json:"ListenerPort"`
}

// MessageBusConf will have configuration data of MessageBusConf
type MessageBusConf struct {
	MessageQueueConfigFilePath string   `json:"MessageQueueConfigFilePath"` // Message Queue Config File Path
	EmbType                    string   `json:"MessageBusType"`
	EmbQueue                   []string `json:"MessageBusQueue"`
}

//KeyCertConf is for holding all security oriented configuration
type KeyCertConf struct {
	RootCACertificatePath string `json:"RootCACertificatePath"` // RootCACertificate will be added to truststore
	PrivateKeyPath        string `json:"PrivateKeyPath"`        // plugin private key
	CertificatePath       string `json:"CertificatePath"`       // plugin certificate
	RootCACertificate     []byte
	PrivateKey            []byte
	Certificate           []byte
}

// URLTranslation ...
type URLTranslation struct {
	NorthBoundURL map[string]string `json:"NorthBoundURL"` // holds value of NorthBound Translation
	SouthBoundURL map[string]string `json:"SouthBoundURL"` // holds value of SouthBound Translation
}

// TLSConf holds TLS confifurations used in https queries
type TLSConf struct {
	MinVersion            string   `json:"MinVersion"`
	MaxVersion            string   `json:"MaxVersion"`
	VerifyPeer            bool     `json:"VerifyPeer"`
	PreferredCipherSuites []string `json:"PreferredCipherSuites"`
}

// SetConfiguration will extract the config data from file
func SetConfiguration() error {
	configFilePath := os.Getenv("PLUGIN_CONFIG_FILE_PATH")
	if configFilePath == "" {
		return fmt.Errorf("No value set to environment variable PLUGIN_CONFIG_FILE_PATH")
	}
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read the config file: %v", err)
	}
	err = json.Unmarshal(configData, &Data)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config data: %v", err)
	}

	return ValidateConfiguration()
}

// ValidateConfiguration will validate configurations read and assign default values, where required
func ValidateConfiguration() error {
	if err := lutilconf.CheckRootServiceuuid(Data.RootServiceUUID); err != nil {
		return err
	}
	if Data.FirmwareVersion == "" {
		log.Warn("No value set for FirmwareVersion, setting default value")
		Data.FirmwareVersion = "1.0"
	}
	if Data.RootServiceUUID == "" {
		return fmt.Errorf("No value set for rootServiceUUID")
	}
	if Data.SessionTimeoutInMinutes == 0 {
		log.Warn("No value set for SessionTimeoutInMinutes, setting default value")
		Data.SessionTimeoutInMinutes = 30
	}
	if err := checkPluginConf(); err != nil {
		return err
	}
	if err := checkEventConf(); err != nil {
		return err
	}
	if err := checkMessageBusConf(); err != nil {
		return err
	}
	if err := checkCertsAndKeysConf(); err != nil {
		return err
	}
	if err := checkTLSConf(); err != nil {
		return err
	}
	checkLBConf()
	checkURLTranslationConf()
	return nil
}

func checkPluginConf() error {
	if Data.PluginConf == nil {
		return fmt.Errorf("No value found for PluginConf")
	}
	if Data.PluginConf.ID == "" {
		log.Warn("No value set for Plugin ID, setting default value")
		Data.PluginConf.ID = "GRF"
	}
	if Data.PluginConf.Host == "" {
		return fmt.Errorf("No value set for plugin Host")
	}
	if Data.PluginConf.Port == "" {
		return fmt.Errorf("No value set for plugin Port")
	}
	if Data.PluginConf.UserName == "" {
		return fmt.Errorf("No value set for plugin Username")
	}
	if Data.PluginConf.Password == "" {
		return fmt.Errorf("No value set for plugin Password")
	}
	return nil
}

//check load balancer configuration
func checkLBConf() {
	if Data.LoadBalancerConf == nil {
		log.Warn("No value set for LoadBalancerConf, setting default value")
		Data.LoadBalancerConf = &LoadBalancerConf{
			Host: Data.EventConf.ListenerHost,
			Port: Data.EventConf.ListenerPort,
		}
		return
	}
	if Data.LoadBalancerConf.Host == "" || Data.LoadBalancerConf.Port == "" {
		log.Warn("No value set for LBHost/LBPort, setting ListenerHost/ListenerPort value")
		Data.LoadBalancerConf.Host = Data.EventConf.ListenerHost
		Data.LoadBalancerConf.Port = Data.EventConf.ListenerPort
	}
}

func checkEventConf() error {
	if Data.EventConf == nil {
		return fmt.Errorf("No value found for EventConf")
	}
	if Data.EventConf.DestURI == "" {
		return fmt.Errorf("No value set for EventURI")
	}
	if Data.EventConf.ListenerHost == "" {
		return fmt.Errorf("No value set for ListenerHost")
	}
	if Data.EventConf.ListenerPort == "" {
		return fmt.Errorf("No value set for ListenerPort")
	}
	return nil
}

//Check or apply default values for message bus to be used by this plugin
func checkMessageBusConf() error {
	if Data.MessageBusConf == nil {
		return fmt.Errorf("No value found for MessageBusConf")
	}
	if _, err := os.Stat(Data.MessageBusConf.MessageQueueConfigFilePath); err != nil {
		return fmt.Errorf("Value check failed for MessageQueueConfigFilePath:%s with %v", Data.MessageBusConf.MessageQueueConfigFilePath, err)
	}
	if Data.MessageBusConf.EmbType == "" {
		log.Warn("No value set for MessageBusType, setting default value")
		Data.MessageBusConf.EmbType = "Kafka"
	}
	if len(Data.MessageBusConf.EmbQueue) <= 0 {
		log.Warn("No value set for MessageBusQueue, setting default value")
		Data.MessageBusConf.EmbQueue = []string{"REDFISH-EVENTS-TOPIC"}
	}
	return nil
}

//Check or apply default values for certs/keys used by this plugin
func checkCertsAndKeysConf() error {
	var err error
	if Data.KeyCertConf == nil {
		return fmt.Errorf("No value found for KeyCertConf")
	}
	if Data.KeyCertConf.Certificate, err = ioutil.ReadFile(Data.KeyCertConf.CertificatePath); err != nil {
		return fmt.Errorf("Value check failed for CertificatePath:%s with %v", Data.KeyCertConf.CertificatePath, err)
	}
	if Data.KeyCertConf.PrivateKey, err = ioutil.ReadFile(Data.KeyCertConf.PrivateKeyPath); err != nil {
		return fmt.Errorf("Value check failed for PrivateKeyPath:%s with %v", Data.KeyCertConf.PrivateKeyPath, err)
	}
	if Data.KeyCertConf.RootCACertificate, err = ioutil.ReadFile(Data.KeyCertConf.RootCACertificatePath); err != nil {
		return fmt.Errorf("Value check failed for RootCACertificatePath:%s with %v", Data.KeyCertConf.RootCACertificatePath, err)
	}
	return nil
}

//Check or apply default values for URL translation from ODIM <=> redfish
func checkURLTranslationConf() {
	if Data.URLTranslation == nil {
		log.Warn("URL translation not provided, setting default value")
		Data.URLTranslation = &URLTranslation{
			NorthBoundURL: map[string]string{
				"ODIM": "redfish",
			},
			SouthBoundURL: map[string]string{
				"redfish": "ODIM",
			},
		}
		return
	}
	if len(Data.URLTranslation.NorthBoundURL) <= 0 {
		log.Warn("NorthBoundURL is empty, setting default value")
		Data.URLTranslation.NorthBoundURL = map[string]string{
			"ODIM": "redfish",
		}
	}
	if len(Data.URLTranslation.SouthBoundURL) <= 0 {
		log.Warn("SouthBoundURL is empty, setting default value")
		Data.URLTranslation.SouthBoundURL = map[string]string{
			"redfish": "ODIM",
		}
	}
}

func checkTLSConf() error {
	if Data.TLSConf == nil {
		log.Warn("TLSConf not provided, setting default value")
		Data.TLSConf = &TLSConf{}
		lutilconf.SetDefaultTLSConf()
		return nil
	}

	var err error
	lutilconf.SetVerifyPeer(Data.TLSConf.VerifyPeer)
	if err = lutilconf.SetTLSMinVersion(Data.TLSConf.MinVersion); err != nil {
		return err
	}
	if err = lutilconf.SetTLSMaxVersion(Data.TLSConf.MaxVersion); err != nil {
		return err
	}
	if err = lutilconf.ValidateConfiguredTLSVersions(); err != nil {
		return err
	}
	if err = lutilconf.SetPreferredCipherSuites(Data.TLSConf.PreferredCipherSuites); err != nil {
		return err
	}
	return nil
}
