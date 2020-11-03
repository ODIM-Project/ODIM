package config

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

type PluginConfig struct {
	ID                      string          `json:"ID"`
	Host                    string          `json:"Host"`
	Port                    string          `json:"Port"`
	UserName                string          `json:"UserName"`
	Password                string          `json:"Password"`
	RootServiceUUID         string          `json:"RootServiceUUID"`
	OdimraNBUrl             string          `json:"OdimraNBUrl"`
	FirmwareVersion         string          `json:"FirmwareVersion"`
	SessionTimeoutInMinutes float64         `json:"SessionTimeoutInMinutes"`
	EventConf               *EventConf      `json:"EventConf"`
	KeyCertConf             *KeyCertConf    `json:"KeyCertConf"`
	URLTranslation          *URLTranslation `json:"URLTranslation"`
	TLSConf                 *TLSConf        `json:"TLSConf"`
}

//EventConf is for holding all events related configuration
type EventConf struct {
	DestURI      string `json:"DestinationURI"`
	ListenerHost string `json:"ListenerHost"`
	ListenerPort string `json:"ListenerPort"`
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

// ReadPluginConfiguration will extract the config data from file
func ReadPluginConfiguration() (*PluginConfig, error) {
	configFilePath := os.Getenv("PLUGIN_CONFIG_FILE_PATH")
	if configFilePath == "" {
		return nil, fmt.Errorf("error: no value set to environment variable PLUGIN_CONFIG_FILE_PATH")
	}
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error: failed to read the config file: %v", err)
	}

	pc := new(PluginConfig)
	err = json.Unmarshal(configData, pc)
	if err != nil {
		return nil, fmt.Errorf("error: failed to unmarshal config data: %v", err)
	}

	err = validate(pc)
	if err != nil {
		return nil, err
	}
	return pc, nil

}

// validate will validate configurations read and assign default values, where required
func validate(pc *PluginConfig) error {

	if pc.OdimraNBUrl == "" {
		return fmt.Errorf("OdimraNBUrl has to be specified")
	}

	if _, e := url.Parse(pc.OdimraNBUrl); e != nil {
		return fmt.Errorf("given OdimraNBUrl is not correct URL")
	}

	if _, err := uuid.Parse(pc.RootServiceUUID); err != nil {
		return err
	}

	if pc.FirmwareVersion == "" {
		log.Println("warn: no value set for FirmwareVersion, setting default value")
		pc.FirmwareVersion = "1.0"
	}

	if pc.RootServiceUUID == "" {
		return fmt.Errorf("error: no value set for rootServiceUUID")
	}

	if pc.SessionTimeoutInMinutes == 0 {
		log.Println("warn: no value set for SessionTimeoutInMinutes, setting default value")
		pc.SessionTimeoutInMinutes = 30
	}

	if pc.ID == "" {
		return fmt.Errorf("no value set for Plugin ID, setting default value")
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

	if pc.EventConf == nil {
		return fmt.Errorf("no value found for EventConf")
	}

	if pc.EventConf.DestURI == "" {
		return fmt.Errorf("o value set for EventURI")
	}

	if pc.EventConf.ListenerHost == "" {
		return fmt.Errorf("no value set for ListenerHost")
	}

	if pc.EventConf.ListenerPort == "" {
		return fmt.Errorf("no value set for ListenerPort")
	}

	if pc.KeyCertConf == nil {
		return fmt.Errorf("error: no value found for KeyCertConf")
	}
	if cert, err := ioutil.ReadFile(pc.KeyCertConf.CertificatePath); err != nil {
		return fmt.Errorf("error: value check failed for CertificatePath:%s with %v", pc.KeyCertConf.CertificatePath, err)
	} else {
		pc.KeyCertConf.Certificate = cert
	}
	if pk, err := ioutil.ReadFile(pc.KeyCertConf.PrivateKeyPath); err != nil {
		return fmt.Errorf("error: value check failed for PrivateKeyPath:%s with %v", pc.KeyCertConf.PrivateKeyPath, err)
	} else {
		pc.KeyCertConf.PrivateKey = pk
	}
	if ca, err := ioutil.ReadFile(pc.KeyCertConf.RootCACertificatePath); err != nil {
		return fmt.Errorf("error: value check failed for RootCACertificatePath:%s with %v", pc.KeyCertConf.RootCACertificatePath, err)
	} else {
		pc.KeyCertConf.RootCACertificate = ca
	}
	if pc.TLSConf == nil {
		return fmt.Errorf("TLSConf not provided, setting default value")
	}

	if pc.URLTranslation == nil {
		log.Println("warn: URL translation not provided, setting default value")
		pc.URLTranslation = &URLTranslation{
			NorthBoundURL: map[string]string{
				"ODIM": "redfish",
			},
			SouthBoundURL: map[string]string{
				"redfish": "ODIM",
			},
		}
	}
	if len(pc.URLTranslation.NorthBoundURL) <= 0 {
		log.Println("warn: NorthBoundURL is empty, setting default value")
		pc.URLTranslation.NorthBoundURL = map[string]string{
			"ODIM": "redfish",
		}
	}
	if len(pc.URLTranslation.SouthBoundURL) <= 0 {
		log.Println("warn: SouthBoundURL is empty, setting default value")
		pc.URLTranslation.SouthBoundURL = map[string]string{
			"redfish": "ODIM",
		}
	}
	return nil
}
