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

package datacommunicator

// -----------------------------------------------------------------------------
// IMPORT Section
// -----------------------------------------------------------------------------
import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
)

// MQCONFIGFILE define the file name and location of the Client config file for
// MQ Platform (contains KAFKA related client configurations) in
// TOML format. This configuration file should be placed under common odimra config
// files location to make it easier to deploy

// -----------------------------------------------------------------------------
// CLIENT CONFIGURATION FILE HANDLING
// -----------------------------------------------------------------------------

// Sample Configuration File below
// [KAFKA]
// # Kafka Server List.
// KServers   = "yourLocalhostFQDN"
// # Listening port. DEFAULT = 9092
// KLPort     = 9092
// # Timeout of KAFKA Server connection drop / Keepalive.
// KTimeout   = 10
// # TLS Configuration Data
// KAFKACertFile       = "path/to/kafka/server.crt"
// KAFKAKeyFile        = "path/to/kafka/kafka.key"
// KAFKACAFile         = "path/to/kafka/CA.crt"

// MQF define the configuration File content for KAFKA in Golang
// structure format. These configurations are embedded into MQF structure for direct
// access to the data.
type MQF struct {
	KafkaF       *KafkaF       `toml:"KAFKA"`
	RedisStreams *RedisStreams `toml:"RedisStreams"`
}

// KafkaF defines the KAFKA Server connection configurations. This structure
// will be extended once we are adding the TLS Authentication and Message
// encoding capability.
type KafkaF struct {

	// KServersInfo defines the list of Kafka Server URI/Nodename:port. DEFAULT = [localhost:9092]
	KServersInfo []string `toml:"KServersInfo"`
	// KTimeout defines the timeout for Kafka Server connection.
	// DEFAULT = 10 (in seconds)
	KTimeout int `toml:"KTimeout"`
	// KAFKACertFile defines the TLS Certificate File for KAFKA. No DEFAULT
	KAFKACertFile string `toml:"KAFKACertFile"`
	// KAFKAKeyFile defines the TLS Key File for KAFKA. No DEFAULT
	KAFKAKeyFile string `toml:"KAFKAKeyFile"`
	// KAFKACAFile defines the KAFKA Certification Authority. No DEFAULT
	KAFKACAFile string `toml:"KAFKACAFile"`
}

// RedisStreams  defines the Redis  connection configurations.
type RedisStreams struct {
	RedisServerAddress string `toml:"RedisServerAddress"`
	RedisServerPort    string `toml:"RedisServerPort"`
	SentinalAddress    string `toml:"SentinalAddress"`
}

// MQ Create both MQF and KafkaPacket Objects. MQF will be used to store
// all config information including Server URL, Port, User credentials
// and other configuration information, which is for Future Expansion.
var MQ MQF

// SetConfiguration defines the function to read the client side configuration file any
// configuration data, which need / should be provided by MQ user would be taken
// directly from the user by asking to fill a structure.  THIS DATA DETAILS
// SHOULD BE DEFINED AS PART OF INTERFACE DEFINITION.
func SetConfiguration(filePath string) error {
	if _, err := toml.DecodeFile(filePath, &MQ); err != nil {
		return fmt.Errorf("Configuration File - %v Read Error: %v", filePath, err)
	}
	if MQ.KafkaF != nil {
		if len(MQ.KafkaF.KServersInfo) <= 0 {
			return fmt.Errorf("no value found for KServersInfo in messagebus config file")
		}
		if MQ.KafkaF.KTimeout == 0 {
			log.Warn("no value found for KTimeout in messagebus config file, using default time 10 seconds")
			MQ.KafkaF.KTimeout = 10
		}
		if MQ.KafkaF.KAFKACertFile == "" {
			return fmt.Errorf("no value found for KAFKACertFile in messagebus config file")
		}
		if MQ.KafkaF.KAFKAKeyFile == "" {
			return fmt.Errorf("no value found for KAFKAKeyFile in messagebus config file")
		}
		if MQ.KafkaF.KAFKACAFile == "" {
			return fmt.Errorf("no value found for KAFKACAFile in messagebus config file")
		}
	}
	return nil
}
