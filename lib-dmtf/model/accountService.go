//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// RestrictedPrivileges - The set of restricted Redfish privileges
type RestrictedPrivileges string

// SupportedAccountTypes - The account types supported by the service
type SupportedAccountTypes string

// LocalAccountAuth - An indication of how the service uses the accounts collection within
// this account service as part of authentication
type LocalAccountAuth string

// AuthenticationTypes - The type of authentication used to connect to the external account provider
type AuthenticationTypes string

// CertificateMappingAttribute - The client certificate attribute to map to a user
type CertificateMappingAttribute string

// AccountProviderType - The type of external account provider to which this service connects
type AccountProviderType string

// Mode -
type Mode string

// PrivilegeType - The set of restricted Redfish privileges
type PrivilegeType string

// PasswordExchangeProtocols - Indicates the allowed TACACS+ password exchange protocols.
type PasswordExchangeProtocols string

const (
	// RestrictedPrivilegesLogin - PrivilegeType
	RestrictedPrivilegesLogin RestrictedPrivileges = "Login"

	// RestrictedPrivilegesConfigureManager - PrivilegeType
	RestrictedPrivilegesConfigureManager RestrictedPrivileges = "ConfigureManager"

	// RestrictedPrivilegesConfigureUsers - PrivilegeType
	RestrictedPrivilegesConfigureUsers RestrictedPrivileges = "ConfigureUsers"

	//RestrictedPrivilegesConfigureSelf - PrivilegeType
	RestrictedPrivilegesConfigureSelf RestrictedPrivileges = "ConfigureSelf"

	//RestrictedPrivilegesConfigureComponents - PrivilegeType
	RestrictedPrivilegesConfigureComponents RestrictedPrivileges = "ConfigureComponents"

	//RestrictedPrivilegesNoAuth - PrivilegeType
	RestrictedPrivilegesNoAuth RestrictedPrivileges = "NoAuth"

	//RestrictedPrivilegesConfigureCompositionInfrastructure - PrivilegeType
	RestrictedPrivilegesConfigureCompositionInfrastructure RestrictedPrivileges = "ConfigureCompositionInfrastructure"

	//RestrictedPrivilegesAdministrateSystems - PrivilegeType
	RestrictedPrivilegesAdministrateSystems RestrictedPrivileges = "AdministrateSystems"

	//RestrictedPrivilegesOperateSystems - PrivilegeType
	RestrictedPrivilegesOperateSystems RestrictedPrivileges = "OperateSystems"

	//RestrictedPrivilegesAdministrateStorage - PrivilegeType
	RestrictedPrivilegesAdministrateStorage RestrictedPrivileges = "AdministrateStorage"

	//RestrictedPrivilegesOperateStorageBackup - PrivilegeType
	RestrictedPrivilegesOperateStorageBackup RestrictedPrivileges = "OperateStorageBackup"

	//SupportedAccountTypesRedfish - AccountTypes
	SupportedAccountTypesRedfish SupportedAccountTypes = "Redfish"

	//SupportedAccountTypesSNMP - AccountTypes
	SupportedAccountTypesSNMP SupportedAccountTypes = "SNMP"

	//SupportedAccountTypesOEM - AccountTypes
	SupportedAccountTypesOEM SupportedAccountTypes = "OEM"

	//SupportedAccountTypeHostConsole - AccountTypes
	SupportedAccountTypeHostConsole SupportedAccountTypes = "HostConsole"

	//SupportedAccountTypesManagerConsole - AccountTypes
	SupportedAccountTypesManagerConsole SupportedAccountTypes = "ManagerConsole"

	//SupportedAccountTypesIPMI  - AccountTypes
	SupportedAccountTypesIPMI SupportedAccountTypes = "IPMI"

	//SupportedAccountTypesKVMIP  - AccountTypes
	SupportedAccountTypesKVMIP SupportedAccountTypes = "KVMIP"

	//SupportedAccountTypesVirtualMedia - AccountTypes
	SupportedAccountTypesVirtualMedia SupportedAccountTypes = "VirtualMedia"

	//SupportedAccountTypesWebUI - AccountTypes
	SupportedAccountTypesWebUI SupportedAccountTypes = "WebUI"

	//LocalAccountAuthEnabled - The service never authenticates users based on the account service-defined accounts collection
	LocalAccountAuthEnabled LocalAccountAuth = "Enabled"

	//LocalAccountAuthDisabled -The service authenticates users based on the account service-defined accounts collection.
	LocalAccountAuthDisabled LocalAccountAuth = "Disabled"

	//LocalAccountAuthFallback - The service authenticates users based on the account service-defined accounts collection
	//only if any external account providers are currently unreachable
	LocalAccountAuthFallback LocalAccountAuth = "Fallback"

	//LocalAccountAuthLocalFirst - The service first authenticates users based on the account service-defined accounts collection.
	//  If authentication fails, the service authenticates by using external account providers
	LocalAccountAuthLocalFirst LocalAccountAuth = "LocalFirst"

	//AuthenticationTypesToken - An opaque authentication token
	AuthenticationTypesToken AuthenticationTypes = "Token"

	//AuthenticationTypesKerberosKeytab - A Kerberos keytab
	AuthenticationTypesKerberosKeytab AuthenticationTypes = "KerberosKeytab"

	//AuthenticationTypesUsernameAndPassword - A user name and password combination
	AuthenticationTypesUsernameAndPassword AuthenticationTypes = "UsernameAndPassword"

	//AuthenticationTypesOEM - An OEM-specific authentication mechanism.
	AuthenticationTypesOEM AuthenticationTypes = "OEM"

	//CertificateMappingAttributeWhole -Match the whole certificate
	CertificateMappingAttributeWhole CertificateMappingAttribute = "Whole"

	//CertificateMappingAttributeCommonName - Match the Common Name (CN) field in the provided certificate to the username
	CertificateMappingAttributeCommonName CertificateMappingAttribute = "CommonName"

	//CertificateMappingAttributeUserPrincipalName - Match the User Principal Name (UPN) field in the provided certificate to the username
	CertificateMappingAttributeUserPrincipalName CertificateMappingAttribute = "UserPrincipalName"

	//AccountProviderTypeRedfishService - An external Redfish service.
	AccountProviderTypeRedfishService AccountProviderType = "RedfishService"

	//AccountProviderTypeActiveDirectoryService - An external Active Directory service
	AccountProviderTypeActiveDirectoryService AccountProviderType = "ActiveDirectoryService"

	//AccountProviderTypeLDAPService -A generic external LDAP service
	AccountProviderTypeLDAPService AccountProviderType = "LDAPService"

	//AccountProviderTypeOEM - An OEM-specific external authentication or directory service
	AccountProviderTypeOEM AccountProviderType = "OEM"

	//AccountProviderTypeTACACSplus - "An external TACACS+ service
	AccountProviderTypeTACACSplus AccountProviderType = "TACACSplus"

	//AccountProviderTypeOAuth2 -An external OAuth 2.0 service
	AccountProviderTypeOAuth2 AccountProviderType = "OAuth2"

	//ModeDiscovery - OAuth 2.0 service information for token validation is downloaded by the service
	ModeDiscovery Mode = "Discovery"

	//ModeOffline - OAuth 2.0 service information for token validation is configured by a client
	ModeOffline Mode = "Offline"

	//PasswordExchangeProtocolsASCII - The ASCII Login method
	PasswordExchangeProtocolsASCII PasswordExchangeProtocols = "ASCII"

	//PasswordExchangeProtocolsPAP - The PAP Login method
	PasswordExchangeProtocolsPAP PasswordExchangeProtocols = "PAP"

	//PasswordExchangeProtocolsCHAP - The CHAP Login method
	PasswordExchangeProtocolsCHAP PasswordExchangeProtocols = "CHAP"

	//PasswordExchangeProtocolsMSCHAPv1 - The MS-CHAP v1 Login method
	PasswordExchangeProtocolsMSCHAPv1 PasswordExchangeProtocols = "MSCHAPv1"

	//PasswordExchangeProtocolsMSCHAPv2 - The MS-CHAP v2 Login method
	PasswordExchangeProtocolsMSCHAPv2 PasswordExchangeProtocols = "MSCHAPv2"

	//PrivilegeTypeLogin -Can log in to the service and read Resources
	PrivilegeTypeLogin PrivilegeType = "Login"

	//PrivilegeTypeConfigureManager -Can configure managers
	PrivilegeTypeConfigureManager PrivilegeType = "ConfigureManager"

	//PrivilegeTypeConfigureUsers -Can configure users and their accounts
	PrivilegeTypeConfigureUsers PrivilegeType = "ConfigureUsers"

	//PrivilegeTypeConfigureSelf -Can change the password for the current user account and log out of their own sessions
	PrivilegeTypeConfigureSelf PrivilegeType = "ConfigureSelf"

	//PrivilegeTypeConfigureComponents - Can configure components that this service manages
	PrivilegeTypeConfigureComponents PrivilegeType = "ConfigureComponents"

	//PrivilegeTypeNoAuth -Authentication is not required
	PrivilegeTypeNoAuth PrivilegeType = "NoAuth"

	//PrivilegeTypeConfigureCompositionInfrastructure - Can view and configure composition service resources
	PrivilegeTypeConfigureCompositionInfrastructure PrivilegeType = "ConfigureCompositionInfrastructure"

	//PrivilegeTypeAdministrateSystems - Administrator for systems found in the systems collection.  Able to manage boot configuration, keys, and certificates for systems
	PrivilegeTypeAdministrateSystems PrivilegeType = "AdministrateSystems"

	//PrivilegeTypeOperateSystems -Operator for systems found in the systems collection.  Able to perform resets and configure interfaces
	PrivilegeTypeOperateSystems PrivilegeType = "OperateSystems"

	//PrivilegeTypeAdministrateStorage - Administrator for storage subsystems and storage systems found in the storage collection and storage system collection respectively
	PrivilegeTypeAdministrateStorage PrivilegeType = "AdministrateStorage"

	//PrivilegeTypeOperateStorageBackup -Operator for storage backup functionality for storage subsystems and storage systems found
	//in the storage collection and storage system collection respectively
	PrivilegeTypeOperateStorageBackup PrivilegeType = "OperateStorageBackup"
)

// AccountService the supported properties,
// this structure should be updated once ODIMRA supports more properties
// The AccountService schema defines an account service.  The properties are common to, and enable management of,
// all user accounts.  The properties include the password requirements and control features, such as account
// lockout.  Properties and actions in this service specify general behavior that should be followed
// for typical accounts, however implementations may override these behaviors for special accounts
// or situations to avoid denial of service or other deadlock situations.
// Reference  :AccountService.v1_12_0.json
type AccountService struct {
	ODataContext                       string                              `json:"@odata.context,omitempty"`
	ODataEtag                          string                              `json:"@odata.etag,omitempty"`
	ODataID                            string                              `json:"@odata.id"`
	ODataType                          string                              `json:"@odata.type"`
	AccountLockoutCounterResetAfter    int                                 `json:"AccountLockoutCounterResetAfter,omitempty"`
	AccountLockoutCounterResetEnabled  bool                                `json:"AccountLockoutCounterResetEnabled,omitempty"`
	AccountLockoutDuration             int                                 `json:"AccountLockoutDuration,omitempty"`
	AccountLockoutThreshold            int                                 `json:"AccountLockoutThreshold,omitempty"`
	Actions                            *OemActions                         `json:"Actions,omitempty"`
	ActiveDirectory                    *ExternalAccountProvider            `json:"ActiveDirectory,omitempty"`
	AdditionalExternalAccountProviders *AdditionalExternalAccountProviders `json:"AdditionalExternalAccountProviders,omitempty"`
	AuthFailureLoggingThreshold        int                                 `json:"AuthFailureLoggingThreshold,omitempty"`
	LDAP                               *ExternalAccountProvider            `json:"LDAP,omitempty"`
	MultiFactorAuth                    *MultiFactorAuth                    `json:"MultiFactorAuth,omitempty"`
	OAuth2                             *ExternalAccountProvider            `json:"OAuth2,omitempty"`
	Oem                                *Oem                                `json:"Oem,omitempty"`
	PrivilegeMap                       *PrivilegeMap                       `json:"PrivilegeMap,omitempty"`
	RestrictedOemPrivileges            []string                            `json:"RestrictedOemPrivileges,omitempty"`
	RestrictedPrivileges               []string                            `json:"RestrictedPrivileges,omitempty"`  //enum
	SupportedAccountTypes              []string                            `json:"SupportedAccountTypes,omitempty"` //enum
	SupportedOEMAccountTypes           []string                            `json:"SupportedOEMAccountTypes,omitempty"`
	TACACSplus                         *ExternalAccountProvider            `json:"TACACSplus,omitempty"`
	ID                                 string                              `json:"Id"`
	Name                               string                              `json:"Name"`
	Description                        string                              `json:"Description,omitempty"`
	Status                             Status                              `json:"Status,omitempty"`
	Accounts                           Link                                `json:"Accounts,omitempty"`
	Roles                              Link                                `json:"Roles,omitempty"`
	MinPasswordLength                  int                                 `json:"MinPasswordLength,omitempty"`
	MaxPasswordLength                  int                                 `json:"MaxPasswordLength,omitempty"`
	PasswordExpirationDays             int                                 `json:"PasswordExpirationDays,omitempty"`
	ServiceEnabled                     bool                                `json:"ServiceEnabled,omitempty"`
	LocalAccountAuth                   string                              `json:"LocalAccountAuth,omitempty"` //enum
}

// Authentication redfish structure
// The information required to authenticate to the external service
// This type shall contain the information required to authenticate to the external service
type Authentication struct {
	AuthenticationType string `json:"AuthenticationType,omitempty"` //enum
	EncryptionKey      string `json:"EncryptionKey,omitempty"`
	EncryptionKeySet   bool   `json:"EncryptionKeySet,omitempty"`
	KerberosKeytab     string `json:"KerberosKeytab,omitempty"`
	Oem                *Oem   `json:"Oem,omitempty"`
	Password           string `json:"Password,omitempty"`
	Username           string `json:"Username,omitempty"`
}

// ClientCertificate redfish structure
// Various settings for client certificate authentication such as mTLS or CAC/PIV
// This type shall contain settings for client certificate authentication
type ClientCertificate struct {
	CertificateMappingAttribute     string       `json:"CertificateMappingAttribute,omitempty"` //enum
	Certificates                    Certificates `json:"Certificates,omitempty"`
	Enabled                         bool         `json:"Enabled,omitempty"`
	RespondToUnauthenticatedClients bool         `json:"RespondToUnauthenticatedClients,omitempty"`
}

// ExternalAccountProvider redfish structure
// The external account provider services that can provide accounts for this manager to use for authentication
// This type shall contain properties that represent external account provider services
// that can provide accounts for this manager to use for authentication
type ExternalAccountProvider struct {
	AccountProviderType string             `json:"AccountProviderType,omitempty"` //enum
	Authentication      *Authentication    `json:"Authentication,omitempty"`
	Certificates        *Certificates      `json:"Certificates,omitempty"`
	LDAPService         *LDAPService       `json:"LDAPService,omitempty"`
	OAuth2Service       *OAuth2Service     `json:"OAuth2Service,omitempty"`
	PasswordSet         bool               `json:"PasswordSet,omitempty"`
	Priority            int                `json:"Priority"`
	RemoteRoleMapping   *RoleMapping       `json:"RemoteRoleMapping"`
	ServiceAddresses    []string           `json:"ServiceAddresses,omitempty"`
	ServiceEnabled      bool               `json:"ServiceEnabled,omitempty"`
	TACACSplusService   *TACACSplusService `json:"TACACSplusService,omitempty"`
}

// GoogleAuthenticator redfish structure
// Various settings for Google Authenticator multi-factor authentication
// This type shall contain settings for Google Authenticator multi-factor authentication
type GoogleAuthenticator struct {
	Enabled      bool   `json:"Enabled,omitempty"`
	SecretKey    string `json:"SecretKey,omitempty"`
	SecretKeySet bool   `json:"SecretKeySet,omitempty"`
}

// LDAPSearchSettings redfish structure
// The settings to search a generic LDAP service
// This type shall contain all required settings to search a generic LDAP service
type LDAPSearchSettings struct {
	BaseDistinguishedNames []string `json:"BaseDistinguishedNames,omitempty"`
	GroupNameAttribute     string   `json:"GroupNameAttribute,omitempty"`
	GroupsAttribute        string   `json:"GroupsAttribute,omitempty"`
	SSHKeyAttribute        string   `json:"SSHKeyAttribute,omitempty"`
	UsernameAttribute      string   `json:"UsernameAttribute,omitempty"`
}

// LDAPService redfish structure
// The settings required to parse a generic LDAP service
// This type shall contain all required settings to parse a generic LDAP service
type LDAPService struct {
	Oem            *Oem                `json:"Oem,omitempty"`
	SearchSettings *LDAPSearchSettings `json:"SearchSettings,omitempty"`
}

// MFABypass redfish structure
// Multi-factor authentication bypass settings
// This type shall contain multi-factor authentication bypass settings
type MFABypass struct {
	BypassTypes []string `json:"BypassTypes,omitempty"`
}

// MicrosoftAuthenticator redfish structure
// Various settings for Microsoft Authenticator multi-factor authentication
// This type shall contain settings for Microsoft Authenticator multi-factor authentication
type MicrosoftAuthenticator struct {
	Enabled      bool   `json:"Enabled,omitempty"`
	SecretKey    string `json:"SecretKey,omitempty"`
	SecretKeySet bool   `json:"SecretKeySet,omitempty"`
}

// MultiFactorAuth redfish structure
// Multi-factor authentication settings
// This type shall contain multi-factor authentication settings
type MultiFactorAuth struct {
	ClientCertificate      *ClientCertificate      `json:"ClientCertificate,omitempty"`
	GoogleAuthenticator    *GoogleAuthenticator    `json:"GoogleAuthenticator,omitempty"`
	MicrosoftAuthenticator *MicrosoftAuthenticator `json:"MicrosoftAuthenticator,omitempty"`
	SecurID                *SecurID                `json:"SecurID,omitempty"`
}

// OAuth2Service redfish structure
// Various settings to parse an OAuth 2.0 service
// This type shall contain settings for parsing an OAuth 2.0 service
type OAuth2Service struct {
	Audience                []string `json:"Audience,omitempty"`
	Issuer                  string   `json:"Issuer,omitempty"`
	Mode                    string   `json:"Mode,omitempty"` //enum
	OAuthServiceSigningKeys string   `json:"OAuthServiceSigningKeys,omitempty"`
}

// SecurID redfish structure
// Various settings for RSA SecurID multi-factor authentication
// This type shall contain settings for RSA SecurID multi-factor authentication
type SecurID struct {
	Certificates    *Certificates `json:"Certificates,omitempty"`
	ClientID        string        `json:"ClientID,omitempty"`
	ClientSecret    string        `json:"ClientSecret,omitempty"`
	ClientSecretSet bool          `json:"ClientSecretSet,omitempty"`
	Enabled         bool          `json:"Enabled,omitempty"`
	ServerURI       string        `json:"ServerURI,omitempty"`
}

// RoleMapping redfish structure
// The mapping rules that are used to convert the external account providers account
// information to the local Redfish role
// This type shall contain mapping rules that are used to convert the external account
// providers account information to the local Redfish role
type RoleMapping struct {
	LocalRole   string     `json:"LocalRole,omitempty"`
	MFABypass   *MFABypass `json:"MFABypass,omitempty"`
	Oem         *Oem       `json:"Oem,omitempty"`
	RemoteGroup string     `json:"RemoteGroup,omitempty"`
	RemoteUser  string     `json:"RemoteUser,omitempty"`
}

// TACACSplusService redfish structure
// Various settings to parse a TACACS+ service
// This type shall contain settings for parsing a TACACS+ service
type TACACSplusService struct {
	PasswordExchangeProtocols string `json:"PasswordExchangeProtocols,omitempty"` //enum
	PrivilegeLevelArgument    string `json:"PrivilegeLevelArgument,omitempty"`
}

// AdditionalExternalAccountProviders redfish structure
type AdditionalExternalAccountProviders struct {
	ODataContext         string   `json:"@odata.context,omitempty"`
	ODataEtag            string   `json:"@odata.etag,omitempty"`
	ODataID              string   `json:"@odata.id"`
	ODataType            string   `json:"@odata.type"`
	Description          string   `json:"Description,omitempty"`
	Members              []string `json:"Members"`
	MembersODataCount    int      `json:"Members@odata.count"`
	MembersODataNextLink string   `json:"Members@odata.nextLink,omitempty"`
	Name                 string   `json:"Name"`
	Oem                  *Oem     `json:"Oem,omitempty"`
}

// PrivilegeMap redfish structure
type PrivilegeMap struct {
	ODataType         string      `json:"@odata.type"`
	Actions           *OemActions `json:"Actions,omitempty"`
	Description       string      `json:"Description,omitempty"`
	ID                string      `json:"Id"`
	Mapping           *Mapping    `json:"Mapping,omitempty"`
	Name              string      `json:"Name"`
	OEMPrivilegesUsed []string    `json:"OEMPrivilegesUsed,omitempty"`
	Oem               Oem         `json:"Oem,omitempty"`
	PrivilegesUsed    []string    `json:"PrivilegesUsed,omitempty"` //enum
}

// Mapping redfish structure
type Mapping struct {
	Entity               string             `json:"Entity,omitempty"`
	OperationMap         OperationMap       `json:"OperationMap,omitempty"`
	PropertyOverrides    TargetPrivilegeMap `json:"PropertyOverrides,omitempty"`
	ResourceURIOverrides TargetPrivilegeMap `json:"ResourceURIOverrides,omitempty"`
	SubordinateOverrides TargetPrivilegeMap `json:"SubordinateOverrides,omitempty"`
}

// TargetPrivilegeMap redfish structure
type TargetPrivilegeMap struct {
	OperationMap OperationMap `json:"OperationMap,omitempty"`
	Targets      []string     `json:"Targets,omitempty"`
}

// OperationMap redfish structure
type OperationMap struct {
	DELETE OperationPrivilege `json:"DELETE,omitempty"`
	GET    OperationPrivilege `json:"GET,omitempty"`
	HEAD   OperationPrivilege `json:"HEAD,omitempty"`
	POST   OperationPrivilege `json:"POST,omitempty"`
	PUT    OperationPrivilege `json:"PUT,omitempty"`
	PATCH  OperationPrivilege `json:"PATCH,omitempty"`
}

// OperationPrivilege redfish structure
type OperationPrivilege struct {
	Privilege []string `json:"Privilege,omitempty"`
}

// ManagerAccount the supported properties of manager account schema,
// this structure should be updated once ODIMRA supports more properties
type ManagerAccount struct {
	ODataContext           string       `json:"@odata.context,omitempty"`
	ODataEtag              string       `json:"@odata.etag,omitempty"`
	ODataID                string       `json:"@odata.id"`
	ODataType              string       `json:"@odata.type"`
	ID                     string       `json:"Id"`
	Name                   string       `json:"Name"`
	Description            string       `json:"Description,omitempty"`
	UserName               string       `json:"UserName,omitempty"`
	Password               string       `json:"Password,omitempty"`
	RoleID                 string       `json:"RoleId,omitempty"`
	Enabled                bool         `json:"Enabled,omitempty"`
	Locked                 bool         `json:"Locked,omitempty"`
	PasswordChangeRequired bool         `json:"PasswordChangeRequired,omitempty"`
	PasswordExpiration     string       `json:"PasswordExpiration,omitempty"`
	AccountExpiration      string       `json:"AccountExpiration,omitempty"`
	Links                  AccountLinks `json:"Links,omitempty"`
	AccountTypes           string       `json:"AccountTypes,omitempty"`
	Keys                   *Collection  `json:"Keys,omitempty"`
}

// AccountLinks struct definition
type AccountLinks struct {
	Role Link `json:"Role"`
}

// Role the supported properties of role schema,
// this structure should be updated once ODIMRA supports more properties
type Role struct {
	ODataContext       string   `json:"@odata.context,omitempty"`
	ODataEtag          string   `json:"@odata.etag,omitempty"`
	ODataID            string   `json:"@odata.id"`
	ODataType          string   `json:"@odata.type"`
	ID                 string   `json:"Id"`
	Name               string   `json:"Name"`
	Description        string   `json:"Description,omitempty"`
	AlternateRoleID    string   `json:"AlternateRoleId,omitempty"`
	AssignedPrivileges []string `json:"AssignedPrivileges,omitempty"`
	IsPredefined       bool     `json:"IsPredefined,omitempty"`
	Restricted         bool     `json:"Restricted,omitempty"`
	RoleID             string   `json:"RoleId,omitempty"`
}
