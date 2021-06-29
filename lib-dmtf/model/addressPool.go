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

package model

// AddressPool is the redfish AddressPool model according to the 2020.3 release
type AddressPool struct {
	ODataContext string               `json:"@odata.context,omitempty"`
	ODataEtag    string               `json:"@odata.etag,omitempty"`
	ODataID      string               `json:"@odata.id"`
	ODataType    string               `json:"@odata.type"`
	Actions      *OemActions          `json:"Actions,omitempty"`
	Description  string               `json:"Description,omitempty"`
	Ethernet     *AddressPoolEthernet `json:"Ethernet,omitempty"`
	GenZ         *GenZ                `json:"GenZ,omitempty"`
	ID           string               `json:"id"`
	Links        *AddressPoolLinks    `json:"Links,omitempty"`
	Name         string               `json:"Name"`
	Oem          interface{}          `json:"Oem,omitempty"`
	Status       *Status              `json:"Status,omitempty"`
}

// OemActions redfish model
type OemActions struct {
	Oem *Oem `json:"Oem,omitempty"`
}

// AddressPoolEthernet redfish model
type AddressPoolEthernet struct {
	BFDSingleHopOnly  *BFDSingleHopOnly    `json:"BFDSingleHopOnly,omitempty"`
	BGPEvpn           *BGPEvpn             `json:"BGPEvpn,omitempty"`
	EBGP              *EBGP                `json:"EBGP,omitempty"`
	IPv4              *IPv4                `json:"IPv4,omitempty"`
	MultiProtocolEBGP *EBGP                `json:"MultiProtocolEBGP,omitempty"`
	MultiProtocolIBGP *CommonBGPProperties `json:"MultiProtocolIBGP,omitempty"`
	SystemMACRange    *AddressRange        `json:"SystemMACRange,omitempty"`
}

// GenZ redfish model
type GenZ struct {
	AccessKey string `json:"AccessKey,omitempty"`
	MaxCID    int    `json:"MaxCID,omitempty"`
	MaxSID    int    `json:"MaxSID,omitempty"`
	MinCID    int    `json:"MinCID,omitempty"`
	MinSID    int    `json:"MinSID,omitempty"`
}

// BFDSingleHopOnly redfish structure
type BFDSingleHopOnly struct {
	DemandModeEnabled                 bool   `json:"DemandModeEnabled,omitempty"`
	DesiredMinTxIntervalMilliseconds  int    `json:"DesiredMinTxIntervalMilliseconds,omitempty"`
	KeyChain                          string `json:"KeyChain,omitempty"`
	LocalMultiplier                   int    `json:"LocalMultiplier,omitempty"`
	MeticulousModeEnabled             bool   `json:"MeticulousModeEnabled,omitempty"`
	RequiredMinRxIntervalMilliseconds int    `json:"RequiredMinRxIntervalMilliseconds,omitempty"`
	SourcePort                        int    `json:"SourcePort,omitempty"`
}

// BGPEvpn redfish structure
type BGPEvpn struct {
	ARPProxyEnabled                  bool          `json:"ARPProxyEnabled,omitempty"`
	ARPSupressionEnabled             bool          `json:"ARPSupressionEnabled,omitempty"`
	AnycastGatewayIPAddress          string        `json:"AnycastGatewayIPAddress,omitempty"`
	AnycastGatewayMACAddress         string        `json:"AnycastGatewayMACAddress,omitempty"`
	ESINumberRange                   *NumberRange  `json:"ESINumberRange,omitempty"`
	EVINumberRange                   *NumberRange  `json:"EVINumberRange,omitempty"`
	GatewayIPAddress                 string        `json:"GatewayIPAddress,omitempty"`
	GatewayIPAddressRange            *AddressRange `json:"GatewayIPAddressRange,omitempty"`
	NDPProxyEnabled                  bool          `json:"NDPProxyEnabled,omitempty"`
	NDPSupressionEnabled             bool          `json:"NDPSupressionEnabled,omitempty"`
	RouteDistinguisherRange          *AddressRange `json:"RouteDistinguisherRange,omitempty"`
	RouteTargetRange                 *AddressRange `json:"RouteTargetRange,omitempty"`
	UnderlayMulticastEnabled         bool          `json:"UnderlayMulticastEnabled,omitempty"`
	UnknownUnicastSuppressionEnabled bool          `json:"UnknownUnicastSuppressionEnabled,omitempty"`
	VLANIdentifierAddressRange       *NumberRange  `json:"VLANIdentifierAddressRange,omitempty"`
}

// EBGP redfish structure
type EBGP struct {
	AllowDuplicateASEnabled bool             `json:"AllowDuplicateASEnabled,omitempty"`
	AllowOverrideASEnabled  bool             `json:"AllowOverrideASEnabled,omitempty"`
	AlwaysCompareMEDEnabled bool             `json:"AlwaysCompareMEDEnabled,omitempty"`
	ASNumberRange           *NumberRange     `json:"ASNumberRange,omitempty"`
	BGPLocalPreference      int              `json:"BGPLocalPreference,omitempty"`
	BGPNeighbor             *BGPNeighbor     `json:"BGPNeighbor,omitempty"`
	BGPRoute                *BGPRoute        `json:"BGPRoute,omitempty"`
	BGPWeight               int              `json:"BGPWeight,omitempty"`
	GracefulRestart         *GracefulRestart `json:"GracefulRestart,omitempty"`
	MED                     int              `json:"MED,omitempty"`
	MultihopEnabled         bool             `json:"MultihopEnabled,omitempty"`
	MultihopTTL             int              `json:"MultihopTTL,omitempty"`
	MultiplePaths           *MultiplePaths   `json:"MultiplePaths,omitempty"`
	SendCommunityEnabled    bool             `json:"SendCommunityEnabled,omitempty"`
}

// IPv4 redfish model
type IPv4 struct {
	AnycastGatewayIPAddress       string        `json:"AnycastGatewayIPAddress,omitempty"`
	AnycastGatewayMACAddress      string        `json:"AnycastGatewayMACAddress,omitempty"`
	DHCP                          *DHCP         `json:"DHCP,omitempty"`
	DNSDomainName                 string        `json:"DNSDomainName,omitempty"`
	DNSServer                     string        `json:"DNSServer,omitempty"`
	DistributeIntoUnderlayEnabled bool          `json:"DistributeIntoUnderlayEnabled,omitempty"`
	EBGPAddressRange              *AddressRange `json:"EBGPAddressRange,omitempty"`
	FabricLinkAddressRange        *AddressRange `json:"FabricLinkAddressRange,omitempty"`
	GatewayIPAddress              string        `json:"GatewayIPAddress,omitempty"`
	HostAddressRange              *AddressRange `json:"HostAddressRange,omitempty"`
	IBGPAddressRange              *AddressRange `json:"IBGPAddressRange,omitempty"`
	LoopbackAddressRange          *AddressRange `json:"LoopbackAddressRange,omitempty"`
	ManagementAddressRange        *AddressRange `json:"ManagementAddressRange,omitempty"`
	NTPOffsetHoursMinutes         int           `json:"NTPOffsetHoursMinutes,omitempty"`
	NTPServer                     string        `json:"NTPServer,omitempty"`
	NTPTimezone                   string        `json:"NTPTimezone,omitempty"`
	NativeVLAN                    int           `json:"NativeVLAN,omitempty"`
	VLANIdentifierAddressRange    *NumberRange  `json:"VLANIdentifierAddressRange,omitempty"`
}

// CommonBGPProperties redfish model
type CommonBGPProperties struct {
	ASNumberRange        *NumberRange     `json:"ASNumberRange,omitempty"`
	BGPNeighbor          *BGPNeighbor     `json:"BGPNeighbor,omitempty"`
	BGPRoute             *BGPRoute        `json:"BGPRoute,omitempty"`
	GracefulRestart      *GracefulRestart `json:"GracefulRestart,omitempty"`
	MultiplePaths        *MultiplePaths   `json:"MultiplePaths,omitempty"`
	SendCommunityEnabled bool             `json:"SendCommunityEnabled,omitempty"`
}

// NumberRange is a common structure ASNumberRange, ESINumberRange, EVINumberRange, etc
type NumberRange struct {
	Lower int `json:"Lower,omitempty"`
	Upper int `json:"Upper,omitempty"`
}

// AddressRange is a common structure EBGPAddressRange, FabricLinkAddressRange, HostAddressRange, etc
type AddressRange struct {
	Lower string `json:"Lower,omitempty"`
	Upper string `json:"Upper,omitempty"`
}

// BGPNeighbor redfish structure
type BGPNeighbor struct {
	Address                             string     `json:"Address,omitempty"`
	AllowOwnASEnabled                   bool       `json:"AllowOwnASEnabled,omitempty"`
	ConnectRetrySeconds                 int        `json:"ConnectRetrySeconds,omitempty"`
	HoldTimeSeconds                     int        `json:"HoldTimeSeconds,omitempty"`
	KeepaliveIntervalSeconds            int        `json:"KeepaliveIntervalSeconds,omitempty"`
	LocalAS                             int        `json:"LocalAS,omitempty"`
	LogStateChangesEnabled              bool       `json:"LogStateChangesEnabled,omitempty"`
	MaxPrefix                           *MaxPrefix `json:"MaxPrefix,omitempty"`
	MinimumAdvertisementIntervalSeconds int        `json:"MinimumAdvertisementIntervalSeconds,omitempty"`
	PassiveModeEnabled                  bool       `json:"PassiveModeEnabled,omitempty"`
	PathMTUDiscoveryEnabled             bool       `json:"PathMTUDiscoveryEnabled,omitempty"`
	PeerAS                              int        `json:"PeerAS,omitempty"`
	ReplacePeerASEnabled                bool       `json:"ReplacePeerASEnabled,omitempty"`
	TCPMaxSegmentSizeBytes              int        `json:"TCPMaxSegmentSizeBytes,omitempty"`
	TreatAsWithdrawEnabled              bool       `json:"TreatAsWithdrawEnabled,omitempty"`
}

// BGPRoute redfish model
type BGPRoute struct {
	AdvertiseInactiveRoutesEnabled bool `json:"AdvertiseInactiveRoutesEnabled,omitempty"`
	DistanceExternal               int  `json:"DistanceExternal,omitempty"`
	DistanceInternal               int  `json:"DistanceInternal,omitempty"`
	DistanceLocal                  int  `json:"DistanceLocal,omitempty"`
	ExternalCompareRouterIDEnabled bool `json:"ExternalCompareRouterIdEnabled,omitempty"`
	FlapDampingEnabled             bool `json:"FlapDampingEnabled,omitempty"`
	SendDefaultRouteEnabled        bool `json:"SendDefaultRouteEnabled,omitempty"`
}

// GracefulRestart redfish model
type GracefulRestart struct {
	GracefulRestartEnabled bool `json:"GracefulRestartEnabled,omitempty"`
	HelperModeEnabled      bool `json:"HelperModeEnabled,omitempty"`
	StaleRoutesTimeSeconds int  `json:"StaleRoutesTimeSeconds,omitempty"`
	TimeSeconds            int  `json:"TimeSeconds,omitempty"`
}

// MultiplePaths redfish model
type MultiplePaths struct {
	MaximumPaths            int  `json:"MaximumPaths,omitempty"`
	UseMultiplePathsEnabled bool `json:"UseMultiplePathsEnabled,omitempty"`
}

// MaxPrefix redfish model
type MaxPrefix struct {
	MaxPrefixNumber             int     `json:"MaxPrefixNumber,omitempty"`
	RestartTimerSeconds         int     `json:"RestartTimerSeconds,omitempty"`
	ShutdownThresholdPercentage float64 `json:"ShutdownThresholdPercentage,omitempty"`
	ThresholdWarningOnlyEnabled bool    `json:"ThresholdWarningOnlyEnabled,omitempty"`
}

// DHCP redfish model
type DHCP struct {
	DHCPInterfaceMTUBytes int    `json:"DHCPInterfaceMTUBytes,omitempty"`
	DHCPRelayEnabled      bool   `json:"DHCPRelayEnabled,omitempty"`
	DHCPServer            string `json:"DHCPServer,omitempty"`
}

// AddressPoolLinks is the struct to links under a AddressPool
type AddressPoolLinks struct {
	Endpoints      []Link      `json:"Endpoints,omitempty"`
	EndpointsCount int         `json:"Endpoints@odata.count,omitempty"`
	Zones          []Link      `json:"Zones,omitempty"`
	ZonesCount     int         `json:"Zones@odata.count,omitempty"`
	Oem            interface{} `json:"Oem,omitempty"`
}
