/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// PluginConfig struct holds configuration of URP plugin
type PluginConfig struct {
	Host               string   `yaml:"Host" envconfig:"HOST"`
	Port               string   `yaml:"Port" envconfig:"PORT"`
	UserName           string   `yaml:"UserName" envconfig:"BASIC_AUTH_USERNAME"`
	Password           string   `yaml:"Password" envconfig:"BASIC_AUTH_PASSWORD"`
	RootServiceUUID    string   `yaml:"RootServiceUUID" envconfig:"SERVICE_ROOT_UUID"`
	OdimURL            string   `yaml:"OdimURL" envconfig:"ODIM_URL"`
	OdimUserName       string   `yaml:"OdimUserName" envconfig:"ODIM_USERNAME"`
	OdimPassword       string   `yaml:"OdimPassword" envconfig:"ODIM_PASSWORD"`
	FirmwareVersion    string   `yaml:"FirmwareVersion" envconfig:"FIRMWARE_VERSION"`
	TLSConf            *TLSConf `yaml:"TLSConf"`
	RSAPrivateKeyPath  string   `yaml:"RSAPrivateKeyPath" envconfig:"RSA_PRIVATE_KEY_PATH"`
	RSAPublicKeyPath   string   `yaml:"RSAPublicKeyPath" envconfig:"RSA_PUBLIC_KEY_PATH"`
	PKIRootCAPath      string   `yaml:"PKIRootCACertificatePath" envconfig:"PKI_ROOT_CA_PATH"`
	PKIPrivateKeyPath  string   `yaml:"PKIPrivateKeyPath" envconfig:"PKI_PRIVATE_KEY_PATH"`
	PKICertificatePath string   `yaml:"PKICertificatePath" envconfig:"PKI_CERTIFICATE_PATH_PATH"`
	LogLevel           string   `yaml:"LogLevel" envconfig:"LOG_LEVEL"`
	RedisAddress       string   `yaml:"RedisAddress" envconfig:"REDIS_ADDRESS"`
	SentinelMasterName string   `yaml:"SentinelMasterName" envconfig:"SENTINEL_MASTER_NAME"`
}

// TLSConf holds details related with URP's NB interface TLS configuration
type TLSConf struct {
	MinVersion uint16 `yaml:"MinVersion"`
	MaxVersion uint16 `yaml:"MaxVersion"`
}

// ReadPluginConfiguration loads URP's configuration from path defined behind PLUGIN_CONFIG_FILE_PATH env variable
func ReadPluginConfiguration() (*PluginConfig, error) {
	pc := new(PluginConfig)

	if cp := os.Getenv("PLUGIN_CONFIG_FILE_PATH"); cp != "" {
		if configData, err := ioutil.ReadFile(cp); err == nil {
			_ = yaml.Unmarshal(configData, pc)
		} else {
			logging.Warnf("Cannot load configuration file: %s", err)
		}
	}

	if err := envconfig.Process("PLUGIN", pc); err != nil {
		logging.Warnf("Cannot load ENV configuration: %s", err)
	}

	return pc, validate(pc)
}

func validate(pc *PluginConfig) error {
	if pc.LogLevel == "" {
		pc.LogLevel = "debug"
	}
	if pc.OdimURL == "" {
		return fmt.Errorf("OdimURL has to be specified")
	}

	if _, e := url.Parse(pc.OdimURL); e != nil {
		return fmt.Errorf("given OdimURL is not correct URL")
	}

	if _, err := uuid.Parse(pc.RootServiceUUID); err != nil {
		return err
	}

	if pc.FirmwareVersion == "" {
		return fmt.Errorf("no value set for FirmwareVersion, setting default: %s", pc.FirmwareVersion)
	}

	if pc.RootServiceUUID == "" {
		return fmt.Errorf("error: no value set for rootServiceUUID")
	}

	if pc.Host == "" {
		return fmt.Errorf("no value set for Plugin Host")
	}

	if pc.Port == "" {
		return fmt.Errorf("error: no value set for Plugin Port")
	}

	if pc.UserName == "" {
		return fmt.Errorf("error: no value set for Plugin Username")
	}

	if pc.Password == "" {
		return fmt.Errorf("error: no value set for Plugin Password")
	}

	if _, err := ioutil.ReadFile(pc.PKICertificatePath); err != nil {
		return fmt.Errorf("error: value check failed for CertificatePath:%s with %v", pc.PKICertificatePath, err)
	}
	if _, err := ioutil.ReadFile(pc.PKIPrivateKeyPath); err != nil {
		return fmt.Errorf("error: value check failed for PrivateKeyPath:%s with %v", pc.PKIPrivateKeyPath, err)
	}
	if _, err := ioutil.ReadFile(pc.PKIRootCAPath); err != nil {
		return fmt.Errorf("error: value check failed for RootCACertificatePath:%s with %v", pc.PKIRootCAPath, err)
	}

	if pc.TLSConf == nil {
		return fmt.Errorf("TLSConf not provided, setting default value")
	}
	if pc.TLSConf.MinVersion == 0 || pc.TLSConf.MinVersion == 0x0301 || pc.TLSConf.MinVersion == 0x0302 {
		return fmt.Errorf("configured TLSConf.{MinVersion} is wrong")
	}
	if pc.TLSConf.MaxVersion == 0 || pc.TLSConf.MaxVersion == 0x0301 || pc.TLSConf.MaxVersion == 0x0302 {
		return fmt.Errorf("configured TLSConf.{MaxVersion} is wrong")
	}
	return nil
}
