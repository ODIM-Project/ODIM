/********************************************************************************
* (C) Copyright [2020] Hewlett Packard Enterprise Development LP		*
*										*
* Licensed under the Apache License, Version 2.0 (the "License"); you may	*
* not use this file except in compliance with the License. You may obtain	*
* a copy of the License at							*
*										*
*    http://www.apache.org/licenses/LICENSE-2.0					*
*										*
* Unless required by applicable law or agreed to in writing, software		*
* distributed under the License is distributed on an "AS IS" BASIS, WITHOUT	*
* WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the	*
* License for the specific language governing permissions and limitations	*
* under the License.								*
*********************************************************************************/

// Package config ...
package config

import "crypto/tls"

// Host defines if the application is Server or client
type Host int8

const (
	// Server is for defining application type is  Server
	Server Host = iota
	// Client is for defining application type is Client
	Client
)

const (
	// DefaultFirmwareVersion - default FirmwareVersion value
	DefaultFirmwareVersion = "1.0"
	// DefaultSessionTimeOutInMins - default SessionTimeOutInMins value
	DefaultSessionTimeOutInMins = 30
	// DefaultExpiredSessionCleanUpTimeInMins - default ExpiredSessionCleanUpTimeInMins value
	DefaultExpiredSessionCleanUpTimeInMins = 15
	// DefaultDBProtocol - default Protocol value
	DefaultDBProtocol = "tcp"
	// DefaultDBMaxActiveConns - default MaxActiveConns value
	DefaultDBMaxActiveConns = 120
	// DefaultDBMaxIdleConns - default MaxIdleConns value
	DefaultDBMaxIdleConns = 10
	// DefaultAuthFailureLoggingThreshold - default AuthFailureLoggingThreshold value
	DefaultAuthFailureLoggingThreshold = 3
	// DefaultAccountLockoutThreshold - default AccountLockoutThreshold value
	DefaultAccountLockoutThreshold = 5
	// DefaultAccountLockoutDuration - default AccountLockoutDuration value
	DefaultAccountLockoutDuration = 30
	// DefaultAccountLockoutCounterResetAfter - default AccountLockoutCounterResetAfter value
	DefaultAccountLockoutCounterResetAfter = 30
	// DefaultMinPasswordLength - default MinPasswordLengt value
	DefaultMinPasswordLength = 12
	// DefaultMaxPasswordLength - default MaxPasswordLength value
	DefaultMaxPasswordLength = 16
	// DefaultAllowedSpecialCharcters - default AllowedSpecialCharcters value
	DefaultAllowedSpecialCharcters = "~!@#$%^&*-+_|(){}:;<>,.?/"
	// DefaultPollingFrequencyInMins - default PollingFrequencyInMins value
	DefaultPollingFrequencyInMins = 30
	// DefaultMaxRetryAttempt - default MaxRetryAttempt value
	DefaultMaxRetryAttempt = 3
	// DefaultRetryIntervalInMins - default RetryIntervalInMins value
	DefaultRetryIntervalInMins = 3
	// DefaultResponseTimeoutInSecs - default ResponseTimeoutInSecs value
	DefaultResponseTimeoutInSecs = 3
	// DefaultStartUpResouceBatchSize - default StartUpResouceBatchSize value
	DefaultStartUpResouceBatchSize = 10
	// DefaultMinResetPriority - default MinResetPriority value
	DefaultMinResetPriority = 1
	// DefaultMaxResetDelay - maximum delay in seconds a reset action can wait
	DefaultMaxResetDelay = 36000
	// DefaultHTTPConnTimeout - default HTTPConnTimeout value
	DefaultHTTPConnTimeout = 10
	// DefaultHTTPMaxIdleConns - default HTTPMaxIdleConns value
	DefaultHTTPMaxIdleConns = 100
	// DefaultHTTPIdleConnTimeout - default HTTPIdleConnTimeout value
	DefaultHTTPIdleConnTimeout = 90
	// DefaultHTTPUseKeepAlive - default HTTPUseKeepAlive value
	DefaultHTTPUseKeepAlive = false
	// DefaultHTTPMaxIdleConnPerHost - default HTTPMaxIdleConnPerHost value
	DefaultHTTPMaxIdleConnPerHost = -1
	// DefaultHTTPExpectContinueTimeout - default HTTPExpectContinueTimeout value
	DefaultHTTPExpectContinueTimeout = 1
	// DefaultTLSHandShakeTimeout - default TLSHandShakeTimeout value
	DefaultTLSHandShakeTimeout = 10
	// DefaultTLSMinVersion - default minimum TLS version supported
	DefaultTLSMinVersion = tls.VersionTLS12
	// DefaultTLSMaxVersion - default maximum TLS version supported
	DefaultTLSMaxVersion = tls.VersionTLS12
	// DefaultTLSServerVerify - indicator for performing server validation
	DefaultTLSServerVerify = true
)

var (
	// DefaultSkipListUnderSystem - holds the default list of resources which needs to be ignored for storing in DB under system resource
	DefaultSkipListUnderSystem = []string{"Chassis", "LogServices", "Managers"}
	// DefaultSkipListUnderManager - holds the default list of resources which needs to be ignored for storing in DB under manager resource
	DefaultSkipListUnderManager = []string{"Chassis", "LogServices", "Systems"}
	// DefaultSkipListUnderChassis - holds the default list of resources which needs to be ignored for storing in DB under chassis resource
	DefaultSkipListUnderChassis = []string{"Managers", "Systems", "Devices"}
	// DefaultSkipListUnderOthers - holds the default list of resources which needs to be ignored for storing in DB under any other resource
	DefaultSkipListUnderOthers = []string{"Power", "Thermal", "SmartStorage"}
	// DefaultCipherSuiteList - default cipher suite list
	DefaultCipherSuiteList = []uint16{
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}
	// SupportedCipherSuitesList - list of cipher suites supported by GO
	// TODO : list needs to be updated, everytime GO adds a new cipher suite
	SupportedCipherSuitesList = map[string]uint16{
		"TLS_RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		"TLS_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	}
	// SupportedTLSVersions - list of TLS versions supported by GO
	// TODO : list needs to be updated, everytime GO adds a new TLS version
	SupportedTLSVersions = map[string]uint16{
		"TLS_1.0": tls.VersionTLS10,
		"TLS_1.1": tls.VersionTLS11,
		"TLS_1.2": tls.VersionTLS12,
	}
)
