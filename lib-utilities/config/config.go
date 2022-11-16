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

// Package config ...
package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// Data will have the configuration data from config file
var Data configModel

// configModel is for holding all the run time configurations for the services
type configModel struct {
	SouthBoundRequestTimeoutInSecs int                      `json:"SouthBoundRequestTimeoutInSecs"` // holds the value of south bound call request time out
	ServerRediscoveryBatchSize     int                      `json:"ServerRediscoveryBatchSize"`
	FirmwareVersion                string                   `json:"FirmwareVersion"`
	RootServiceUUID                string                   `json:"RootServiceUUID"` //static uuid used for root service
	SearchAndFilterSchemaPath      string                   `json:"SearchAndFilterSchemaPath"`
	RegistryStorePath              string                   `json:"RegistryStorePath"`
	LocalhostFQDN                  string                   `json:"LocalhostFQDN"`
	EnabledServices                []string                 `json:"EnabledServices"`
	MessageBusConf                 *MessageBusConf          `json:"MessageBusConf"`
	DBConf                         *DBConf                  `json:"DBConf"`
	KeyCertConf                    *KeyCertConf             `json:"KeyCertConf"`
	AuthConf                       *AuthConf                `json:"AuthConf"`
	APIGatewayConf                 *APIGatewayConf          `json:"APIGatewayConf"`
	AddComputeSkipResources        *AddComputeSkipResources `json:"AddComputeSkipResources"`
	URLTranslation                 *URLTranslation          `json:"URLTranslation"`
	PluginStatusPolling            *PluginStatusPolling     `json:"PluginStatusPolling"`
	ExecPriorityDelayConf          *ExecPriorityDelayConf   `json:"ExecPriorityDelayConf"`
	TLSConf                        *TLSConf                 `json:"TLSConf"`
	SupportedPluginTypes           []string                 `json:"SupportedPluginTypes"`
	ConnectionMethodConf           []ConnectionMethodConf   `json:"ConnectionMethodConf"`
	EventConf                      *EventConf               `json:"EventConf"`
	ResourceRateLimit              []string                 `json:"ResourceRateLimit"`
	RequestLimitCountPerSession    int                      `json:"RequestLimitCountPerSession"`
	SessionLimitCountPerUser       int                      `json:"SessionLimitCountPerUser"`
	LogLevel                       log.Level                `json:"LogLevel"`
	ImageRegistryAddress           string                   `json:"ImageRegistryAddress,omitempty"`
}

// DBConf holds all DB related configurations
type DBConf struct {
	Protocol                      string `json:"Protocol"`
	InMemoryHost                  string `json:"InMemoryHost"`
	InMemoryPort                  string `json:"InMemoryPort"`
	OnDiskHost                    string `json:"OnDiskHost"`
	OnDiskPort                    string `json:"OnDiskPort"`
	MaxIdleConns                  int    `json:"MaxIdleConns"`
	MaxActiveConns                int    `json:"MaxActiveConns"`
	RedisHAEnabled                bool   `json:"RedisHAEnabled"`
	InMemorySentinelPort          string `json:"InMemorySentinelPort"`
	OnDiskSentinelPort            string `json:"OnDiskSentinelPort"`
	InMemoryPrimarySet            string `json:"InMemoryPrimarySet"`
	OnDiskPrimarySet              string `json:"OnDiskPrimarySet"`
	RedisInMemoryPasswordFilePath string `json:"RedisInMemoryPasswordFilePath"`
	RedisOnDiskPasswordFilePath   string `json:"RedisOnDiskPasswordFilePath"`
	RedisInMemoryPassword         []byte
	RedisOnDiskPassword           []byte
}

// MessageBusConf holds all message bus configurations
type MessageBusConf struct {
	MessageBusConfigFilePath string   `json:"MessageBusConfigFilePath"`
	MessageBusType           string   `json:"MessageBusType"`
	MessageBusQueue          []string `json:"MessageBusQueue"`
}

// KeyCertConf is for holding all security oriented configuration
type KeyCertConf struct {
	RootCACertificatePath string `json:"RootCACertificatePath"`
	RPCPrivateKeyPath     string `json:"RPCPrivateKeyPath"`  // location where the Private key is stored
	RPCCertificatePath    string `json:"RPCCertificatePath"` // location where the CA signed certificate is stored
	RSAPublicKeyPath      string `json:"RSAPublicKeyPath"`
	RSAPrivateKeyPath     string `json:"RSAPrivateKeyPath"`
	RootCACertificate     []byte
	RPCPrivateKey         []byte
	RPCCertificate        []byte
	RSAPublicKey          []byte
	RSAPrivateKey         []byte
}

// AuthConf holds all authentication related configurations
type AuthConf struct {
	SessionTimeOutInMins            float64        `json:"SessionTimeOutInMins"`
	ExpiredSessionCleanUpTimeInMins float64        `json:"ExpiredSessionCleanUpTimeInMins"`
	PasswordRules                   *PasswordRules `json:"PasswordRules"`
}

// PasswordRules defines rules for password complexity
type PasswordRules struct {
	MinPasswordLength       int    `json:"MinPasswordLength"`       // holds the value  of min password length
	MaxPasswordLength       int    `json:"MaxPasswordLength"`       // holds the value of max password length
	AllowedSpecialCharcters string `json:"AllowedSpecialCharcters"` // holds all value of  all sppecial charcters
}

// APIGatewayConf holds API gateway related configurations
type APIGatewayConf struct {
	Host            string `json:"Host"`
	Port            string `json:"Port"`
	PrivateKeyPath  string `json:"PrivateKeyPath"`
	CertificatePath string `json:"CertificatePath"`
	PrivateKey      []byte
	Certificate     []byte
}

// AddComputeSkipResources stores list of resources which need to ignored while inserting the contents to DB while adding Computer System
type AddComputeSkipResources struct {
	SkipResourceListUnderSystem  []string `json:"SkipResourceListUnderSystem"`  // holds the list of resources which needs to be ignored for storing in DB under system resource
	SkipResourceListUnderManager []string `json:"SkipResourceListUnderManager"` // holds the list of resources which needs to be ignored for storing in DB under manager resource
	SkipResourceListUnderChassis []string `json:"SkipResourceListUnderChassis"` // holds the list of resources which needs to be ignored for storing in DB under chassis resource
	SkipResourceListUnderOthers  []string `json:"SkipResourceListUnderOthers"`  // holds the list of resources which needs to be ignored for storing in DB under a generic resource apart from system,manager and chassis
}

// URLTranslation ...
type URLTranslation struct {
	NorthBoundURL map[string]string `json:"NorthBoundURL"` // holds value of NorthBound Translation
	SouthBoundURL map[string]string `json:"SouthBoundURL"` // holds value of SouthBound Translation
}

// PluginStatusPolling stores all inforamtion related to status polling
type PluginStatusPolling struct {
	PollingFrequencyInMins  int `json:"PollingFrequencyInMins"` // holds value of  duration in which status polling to be intiated ,value will be in minutes
	MaxRetryAttempt         int `json:"MaxRetryAttempt"`        // holds value number retry attempts
	RetryIntervalInMins     int `json:"RetryIntervalInMins"`    // holds value of  duration in which retry of status polling to be intiated,value will be in minutes
	ResponseTimeoutInSecs   int `json:"ResponseTimeoutInSecs"`  // holds value of duation in which it need wait for resposne ,value will be in seconds
	StartUpResouceBatchSize int `json:"StartUpResouceBatchSize"`
}

// ExecPriorityDelayConf holds priority and delay configurations for exec actions
type ExecPriorityDelayConf struct {
	MinResetPriority    int `json:"MinResetPriority"`
	MaxResetPriority    int `json:"MaxResetPriority"`
	MaxResetDelayInSecs int `json:"MaxResetDelayInSecs"`
}

// TLSConf holds TLS confifurations used in https queries
type TLSConf struct {
	VerifyPeer            bool     `json:"VerifyPeer"`
	MinVersion            string   `json:"MinVersion"`
	MaxVersion            string   `json:"MaxVersion"`
	PreferredCipherSuites []string `json:"PreferredCipherSuites"`
}

// ConnectionMethodConf is for connection method type and variant
type ConnectionMethodConf struct {
	ConnectionMethodType    string `json:"ConnectionMethodType"`
	ConnectionMethodVariant string `json:"ConnectionMethodVariant"`
}

// EventConf stores all inforamtion related to event delivery configurations
type EventConf struct {
	DeliveryRetryAttempts        int `json:"DeliveryRetryAttempts"`        // holds value of retrying event posting to destination
	DeliveryRetryIntervalSeconds int `json:"DeliveryRetryIntervalSeconds"` // holds value of retrying events posting in interval
}

// SetConfiguration will extract the config data from file
func SetConfiguration() error {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		return fmt.Errorf("No value set to environment variable CONFIG_FILE_PATH")
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
	var err error
	if err = CheckRootServiceuuid(Data.RootServiceUUID); err != nil {
		return err
	}
	if err = checkMiscellaneousConf(); err != nil {
		return err
	}
	if err = checkDBConf(); err != nil {
		return err
	}
	if err = checkMessageBusConf(); err != nil {
		return err
	}
	if err = checkKeyCertConf(); err != nil {
		return err
	}
	if err = checkAPIGatewayConf(); err != nil {
		return err
	}
	if err = checkTLSConf(); err != nil {
		return err
	}
	if err = checkConnectionMethodConf(); err != nil {
		return err
	}
	if err = checkEventConf(); err != nil {
		return err
	}
	if err = checkResourceRateLimit(); err != nil {
		return err
	}
	checkAuthConf()
	checkAddComputeSkipResources()
	checkURLTranslation()
	checkPluginStatusPolling()
	checkExecPriorityDelayConf()

	return nil
}

func checkMiscellaneousConf() error {
	if Data.FirmwareVersion == "" {
		log.Warn("No value set for FirmwareVersion, setting default value")
		Data.FirmwareVersion = DefaultFirmwareVersion
	}
	if Data.RootServiceUUID == "" {
		return fmt.Errorf("error: no value set for rootServiceUUID")
	}
	if Data.SouthBoundRequestTimeoutInSecs > 0 {
		DefaultHTTPClient.Timeout = time.Duration(Data.SouthBoundRequestTimeoutInSecs) * time.Second
	}
	if Data.LocalhostFQDN == "" {
		return fmt.Errorf("error: no value set for localhostFQDN")
	}
	if _, err := os.Stat(Data.SearchAndFilterSchemaPath); err != nil {
		return fmt.Errorf("error: value check failed for SearchAndFilterSchemaPath:%s with %v", Data.SearchAndFilterSchemaPath, err)
	}
	if _, err := os.Stat(Data.RegistryStorePath); err != nil {
		return fmt.Errorf("error: value check failed for RegistryStorePath:%s with %v", Data.RegistryStorePath, err)
	}
	if len(Data.EnabledServices) == 0 {
		return fmt.Errorf("error: no value set for EnabledServices")
	}
	if len(Data.SupportedPluginTypes) == 0 {
		return fmt.Errorf("error: no value set for SupportedPluginTypes")
	}
	return nil
}

func checkDBConf() error {
	if Data.DBConf == nil {
		return fmt.Errorf("error: DBConf is not provided")
	}
	if Data.DBConf.Protocol != DefaultDBProtocol {
		log.Warn("Incorrect value configured for DB Protocol, setting default value")
		Data.DBConf.Protocol = DefaultDBProtocol
	}
	if Data.DBConf.InMemoryHost == "" {
		return fmt.Errorf("error: no value configured for DB InMemoryHost")
	}
	if Data.DBConf.InMemoryPort == "" {
		return fmt.Errorf("error: no value configured for DB InMemoryPort")
	}
	if Data.DBConf.OnDiskHost == "" {
		return fmt.Errorf("error: no value configured for DB OnDiskHost")
	}
	if Data.DBConf.OnDiskPort == "" {
		return fmt.Errorf("error: no value configured for DB OnDiskPort")
	}
	if Data.DBConf.MaxActiveConns == 0 {
		log.Warn("No value configured for MaxActiveConns, setting default value")
		Data.DBConf.MaxActiveConns = DefaultDBMaxActiveConns
	}
	if Data.DBConf.MaxIdleConns == 0 {
		log.Warn("No value configured for MaxIdleConns, setting default value")
		Data.DBConf.MaxIdleConns = DefaultDBMaxIdleConns
	}
	if Data.DBConf.RedisHAEnabled {
		if err := checkDBHAConf(); err != nil {
			return err
		}
	}
	var err error
	if Data.DBConf.RedisInMemoryPasswordFilePath != "" && Data.KeyCertConf.RSAPrivateKeyPath != "" {
		if Data.DBConf.RedisInMemoryPassword, err = decryptRSAOAEPEncryptedPasswords(Data.DBConf.RedisInMemoryPasswordFilePath); err != nil {
			return fmt.Errorf("error: while decrypting password from the passwordFilePath:%s with %v", Data.DBConf.RedisInMemoryPasswordFilePath, err)
		}
	}
	if Data.DBConf.RedisOnDiskPasswordFilePath != "" && Data.KeyCertConf.RSAPrivateKeyPath != "" {
		if Data.DBConf.RedisOnDiskPassword, err = decryptRSAOAEPEncryptedPasswords(Data.DBConf.RedisOnDiskPasswordFilePath); err != nil {
			return fmt.Errorf("error: while decrypting password from the passwordFilePath:%s with %v", Data.DBConf.RedisOnDiskPasswordFilePath, err)
		}
	}
	return nil
}

func decryptRSAOAEPEncryptedPasswords(passwordFilePath string) ([]byte, error) {
	privateKeyStr, err := ioutil.ReadFile(Data.KeyCertConf.RSAPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("value check failed for RSAPrivateKeyPath:%s with %v", Data.KeyCertConf.RSAPrivateKeyPath, err)
	}

	block, _ := pem.Decode(privateKeyStr)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key for the RSAPrivateKeyPath:%s",
			Data.KeyCertConf.RSAPrivateKeyPath)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key for the RSAPrivateKeyPath:%s with %v",
			Data.KeyCertConf.RSAPrivateKeyPath, err)
	}

	cipherText, err := ioutil.ReadFile(passwordFilePath)
	if err != nil {
		return nil, fmt.Errorf("value check failed for passwordFilePath:%s with %v", passwordFilePath, err)
	}

	ct, err := base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return nil, fmt.Errorf("value check failed for passwordFilePath:%s with %v", passwordFilePath, err)
	}

	rng := rand.Reader
	password, err := rsa.DecryptOAEP(sha512.New(), rng, privateKey, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("password decryption failed for the passwordFilePath:%s with %v", passwordFilePath, err)
	}

	return password, nil
}

func checkMessageBusConf() error {
	if Data.MessageBusConf == nil {
		return fmt.Errorf("error: MessageBusConf is not provided")
	}
	if Data.MessageBusConf.MessageBusType == "" {
		log.Warn("No value set for MessageBusType, setting default value")
		Data.MessageBusConf.MessageBusType = "Kafka"
	}
	if Data.MessageBusConf.MessageBusType == "Kafka" {
		if _, err := os.Stat(Data.MessageBusConf.MessageBusConfigFilePath); err != nil {
			return fmt.Errorf("Value check failed for MessageBusConfigFilePath:%s with %v", Data.MessageBusConf.MessageBusConfigFilePath, err)
		}
		if len(Data.MessageBusConf.MessageBusQueue) <= 0 {
			log.Warn("No value set for MessageBusQueue, setting default value")
			Data.MessageBusConf.MessageBusQueue = []string{"REDFISH-EVENTS-TOPIC"}
		}
	}
	if !AllowedMessageBusTypes[Data.MessageBusConf.MessageBusType] {
		return fmt.Errorf("error: invalid value configured for MessageBusType")
	}
	return nil
}

func checkDBHAConf() error {
	if Data.DBConf.InMemorySentinelPort == "" {
		return fmt.Errorf("error: no value configured for DB InMemorySentinelPort")
	}
	if Data.DBConf.OnDiskSentinelPort == "" {
		return fmt.Errorf("error: no value configured for DB OnDiskSentinelPort")
	}
	if Data.DBConf.InMemoryPrimarySet == "" {
		return fmt.Errorf("error: no value configured for DB InMemoryPrimarySet")
	}
	if Data.DBConf.OnDiskPrimarySet == "" {
		return fmt.Errorf("error: no value configured for DB OnDiskPrimarySet")
	}
	return nil
}

func checkKeyCertConf() error {
	var err error
	if Data.KeyCertConf == nil {
		return fmt.Errorf("error: KeyCertConf is not provided")
	}
	if Data.KeyCertConf.RootCACertificate, err = ioutil.ReadFile(Data.KeyCertConf.RootCACertificatePath); err != nil {
		return fmt.Errorf("error: value check failed for RootCACertificatePath:%s with %v", Data.KeyCertConf.RootCACertificatePath, err)
	}
	if Data.KeyCertConf.RPCPrivateKey, err = ioutil.ReadFile(Data.KeyCertConf.RPCPrivateKeyPath); err != nil {
		return fmt.Errorf("error: value check failed for RPCPrivateKeyPath:%s with %v", Data.KeyCertConf.RPCPrivateKeyPath, err)
	}
	if Data.KeyCertConf.RPCCertificate, err = ioutil.ReadFile(Data.KeyCertConf.RPCCertificatePath); err != nil {
		return fmt.Errorf("error: value check failed for RPCCertificatePath:%s with %v", Data.KeyCertConf.RPCCertificatePath, err)
	}
	if Data.KeyCertConf.RSAPublicKey, err = ioutil.ReadFile(Data.KeyCertConf.RSAPublicKeyPath); err != nil {
		return fmt.Errorf("error: value check failed for RSAPublicKeyPath:%s with %v", Data.KeyCertConf.RSAPublicKeyPath, err)
	}
	if Data.KeyCertConf.RSAPrivateKey, err = ioutil.ReadFile(Data.KeyCertConf.RSAPrivateKeyPath); err != nil {
		return fmt.Errorf("error: value check failed for RSAPrivateKeyPath:%s with %v", Data.KeyCertConf.RSAPrivateKeyPath, err)
	}
	return nil
}

func checkAuthConf() {
	if Data.AuthConf == nil {
		log.Warn("No value found for AuthConf, setting default value")
		Data.AuthConf = &AuthConf{
			SessionTimeOutInMins:            DefaultSessionTimeOutInMins,
			ExpiredSessionCleanUpTimeInMins: DefaultExpiredSessionCleanUpTimeInMins,
			PasswordRules: &PasswordRules{
				MinPasswordLength:       DefaultMinPasswordLength,
				MaxPasswordLength:       DefaultMaxPasswordLength,
				AllowedSpecialCharcters: DefaultAllowedSpecialCharcters,
			},
		}
		return
	}
	if Data.AuthConf.SessionTimeOutInMins == 0 {
		log.Warn("No value set for SessionTimeOutInMin, setting default value")
		Data.AuthConf.SessionTimeOutInMins = DefaultSessionTimeOutInMins
	}
	if Data.AuthConf.ExpiredSessionCleanUpTimeInMins == 0 {
		log.Warn("No value set for ExpiredSessionCleanUpTimeInMins, setting default value")
		Data.AuthConf.ExpiredSessionCleanUpTimeInMins = DefaultExpiredSessionCleanUpTimeInMins
	}
	checkPasswordRulesConf()
}

func checkPasswordRulesConf() {
	if Data.AuthConf.PasswordRules == nil {
		log.Warn("PasswordRules configuration is found empty, setting default value")
		Data.AuthConf.PasswordRules = &PasswordRules{
			MinPasswordLength:       DefaultMinPasswordLength,
			MaxPasswordLength:       DefaultMaxPasswordLength,
			AllowedSpecialCharcters: DefaultAllowedSpecialCharcters,
		}
		return
	}
	if Data.AuthConf.PasswordRules.MinPasswordLength <= 0 {
		log.Warn("No value set for MinPasswordLength, setting default value")
		Data.AuthConf.PasswordRules.MinPasswordLength = DefaultMinPasswordLength
	}
	if Data.AuthConf.PasswordRules.MaxPasswordLength <= 0 {
		log.Warn("No value set for MaxPasswordLength, setting default value")
		Data.AuthConf.PasswordRules.MaxPasswordLength = DefaultMaxPasswordLength
	}
	if Data.AuthConf.PasswordRules.AllowedSpecialCharcters == "" {
		log.Warn("No value set for AllowedSpecialCharcters, setting default value")
		Data.AuthConf.PasswordRules.AllowedSpecialCharcters = DefaultAllowedSpecialCharcters
	}
}

func checkAPIGatewayConf() error {
	var err error
	if Data.APIGatewayConf == nil {
		return fmt.Errorf("error: APIGatewayConf is not provided")
	}
	if Data.APIGatewayConf.Port == "" {
		return fmt.Errorf("error: no value set for GatewayPort")
	}
	if Data.APIGatewayConf.PrivateKey, err = ioutil.ReadFile(Data.APIGatewayConf.PrivateKeyPath); err != nil {
		return fmt.Errorf("error: value check failed for PrivateKeyPath:%s with %v", Data.APIGatewayConf.PrivateKeyPath, err)
	}
	if Data.APIGatewayConf.Certificate, err = ioutil.ReadFile(Data.APIGatewayConf.CertificatePath); err != nil {
		return fmt.Errorf("error: value check failed for CertificatePath:%s with %v", Data.APIGatewayConf.CertificatePath, err)
	}
	return nil
}

func checkAddComputeSkipResources() {
	if Data.AddComputeSkipResources == nil {
		log.Warn("No value found for AddComputeRetrival, setting default value")
		Data.AddComputeSkipResources = &AddComputeSkipResources{
			SkipResourceListUnderSystem:  DefaultSkipListUnderSystem,
			SkipResourceListUnderManager: DefaultSkipListUnderManager,
			SkipResourceListUnderChassis: DefaultSkipListUnderChassis,
			SkipResourceListUnderOthers:  DefaultSkipListUnderOthers,
		}
		return
	}
	if len(Data.AddComputeSkipResources.SkipResourceListUnderSystem) == 0 {
		log.Warn("No value found for SkipResourceListUnderSystem, setting default value")
		Data.AddComputeSkipResources.SkipResourceListUnderSystem = DefaultSkipListUnderSystem
	}
	if len(Data.AddComputeSkipResources.SkipResourceListUnderManager) == 0 {
		log.Warn("No value found for SkipResourceListUnderManager, setting default value")
		Data.AddComputeSkipResources.SkipResourceListUnderManager = DefaultSkipListUnderManager
	}
	if len(Data.AddComputeSkipResources.SkipResourceListUnderChassis) == 0 {
		log.Warn("No value found for SkipResourceListUnderChassis, setting default value")
		Data.AddComputeSkipResources.SkipResourceListUnderChassis = DefaultSkipListUnderChassis
	}
	if len(Data.AddComputeSkipResources.SkipResourceListUnderOthers) == 0 {
		log.Warn("No value found for SkipResourceListUnderOthers, setting default value")
		Data.AddComputeSkipResources.SkipResourceListUnderOthers = DefaultSkipListUnderOthers
	}
}

func checkURLTranslation() {
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

func checkPluginStatusPolling() {
	if Data.PluginStatusPolling == nil {
		log.Warn("PluginStatusPolling not provided, setting default value")
		Data.PluginStatusPolling = &PluginStatusPolling{
			PollingFrequencyInMins:  DefaultPollingFrequencyInMins,
			MaxRetryAttempt:         DefaultMaxRetryAttempt,
			RetryIntervalInMins:     DefaultRetryIntervalInMins,
			ResponseTimeoutInSecs:   DefaultResponseTimeoutInSecs,
			StartUpResouceBatchSize: DefaultStartUpResouceBatchSize,
		}
		return
	}
	if Data.PluginStatusPolling.PollingFrequencyInMins <= 0 {
		log.Warn("No value found for PollingFrequencyInMins, setting default value")
		Data.PluginStatusPolling.PollingFrequencyInMins = DefaultPollingFrequencyInMins
	}
	if Data.PluginStatusPolling.MaxRetryAttempt <= 0 {
		log.Warn("No value found for MaxRetryAttempt, setting default value")
		Data.PluginStatusPolling.MaxRetryAttempt = DefaultMaxRetryAttempt
	}
	if Data.PluginStatusPolling.RetryIntervalInMins <= 0 {
		log.Warn("No value found for RetryIntervalInMins, setting default value")
		Data.PluginStatusPolling.RetryIntervalInMins = DefaultRetryIntervalInMins
	}
	if Data.PluginStatusPolling.ResponseTimeoutInSecs <= 0 {
		log.Warn("No value found for ResponseTimeoutInSecs, setting default value")
		Data.PluginStatusPolling.ResponseTimeoutInSecs = DefaultResponseTimeoutInSecs
	}
	if Data.PluginStatusPolling.StartUpResouceBatchSize <= 0 {
		log.Warn("No value found for StartUpResouceBatchSize, setting default value")
		Data.PluginStatusPolling.StartUpResouceBatchSize = DefaultStartUpResouceBatchSize
	}
}

func checkExecPriorityDelayConf() {
	if Data.ExecPriorityDelayConf == nil {
		log.Warn("ExecPriorityDelayConf not provided, setting default value")
		Data.ExecPriorityDelayConf = &ExecPriorityDelayConf{
			MinResetPriority:    DefaultMinResetPriority,
			MaxResetPriority:    DefaultMinResetPriority + 1,
			MaxResetDelayInSecs: DefaultMaxResetDelay,
		}
		return
	}
	if Data.ExecPriorityDelayConf.MinResetPriority <= 0 {
		log.Warn("No value found for MinResetPriority, setting default value")
		Data.ExecPriorityDelayConf.MinResetPriority = DefaultMinResetPriority
	}
	if Data.ExecPriorityDelayConf.MaxResetPriority <= Data.ExecPriorityDelayConf.MinResetPriority {
		log.Warn("no value found for MaxResetPriority, setting default value")
		Data.ExecPriorityDelayConf.MaxResetPriority = Data.ExecPriorityDelayConf.MinResetPriority + 1
	}
	if Data.ExecPriorityDelayConf.MaxResetDelayInSecs <= 0 ||
		Data.ExecPriorityDelayConf.MaxResetDelayInSecs > DefaultMaxResetDelay {
		log.Warn("No value found for MaxResetDelayInSecs, setting default value")
		Data.ExecPriorityDelayConf.MaxResetDelayInSecs = DefaultMaxResetDelay
	}
}

func checkTLSConf() error {
	if Data.TLSConf == nil {
		log.Warn("TLSConf not provided, setting default value")
		Data.TLSConf = &TLSConf{}
		SetDefaultTLSConf()
		return nil
	}

	var err error
	SetVerifyPeer(Data.TLSConf.VerifyPeer)
	if err = SetTLSMinVersion(Data.TLSConf.MinVersion); err != nil {
		return err
	}
	if err = SetTLSMaxVersion(Data.TLSConf.MaxVersion); err != nil {
		return err
	}
	if err = ValidateConfiguredTLSVersions(); err != nil {
		return err
	}
	if err = SetPreferredCipherSuites(Data.TLSConf.PreferredCipherSuites); err != nil {
		return err
	}
	return nil
}

//CheckRootServiceuuid function is used to validate format of Root Service UUID. The same function is used in plugin-redfish config.go
func CheckRootServiceuuid(uid string) error {
	_, err := uuid.Parse(uid)
	return err
}

func checkConnectionMethodConf() error {
	var err error
	if len(Data.ConnectionMethodConf) == 0 {
		return fmt.Errorf("error: ConnectionMethodConf is not provided")
	}
	return err
}

func checkEventConf() error {
	if Data.EventConf == nil {
		log.Warn("EventConf not provided, setting default value")
		Data.EventConf = &EventConf{
			DeliveryRetryAttempts:        DefaultDeliveryRetryAttempts,
			DeliveryRetryIntervalSeconds: DefaultDeliveryRetryIntervalSeconds,
		}
		return nil
	}
	if Data.EventConf.DeliveryRetryAttempts <= 0 {
		log.Warn("No value found for DeliveryRetryAttempts, setting default value")
		Data.EventConf.DeliveryRetryAttempts = DefaultDeliveryRetryAttempts
	}
	if Data.EventConf.DeliveryRetryIntervalSeconds <= 0 {
		log.Warn("No value found for DeliveryRetryIntervalSeconds, setting default value")
		Data.EventConf.DeliveryRetryIntervalSeconds = DefaultDeliveryRetryIntervalSeconds
	}
	return nil
}

func checkResourceRateLimit() error {
	for _, val := range Data.ResourceRateLimit {
		resourceLimit := strings.Split(val, ":")
		if len(resourceLimit) > 1 && resourceLimit[1] != "" {
			_, err := strconv.Atoi(resourceLimit[1])
			if err != nil {
				return fmt.Errorf("time should be in integer format: %v", err.Error())
			}
		}
	}
	return nil
}
