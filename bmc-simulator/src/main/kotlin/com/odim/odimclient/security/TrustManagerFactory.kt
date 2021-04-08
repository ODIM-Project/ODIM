/*
 * Copyright (c) Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.odim.odimclient.security

import com.odim.simulator.CoreConfig.EXTERNAL_TRUSTSTORE_LOCATION
import com.odim.simulator.CoreConfig.EXTERNAL_TRUSTSTORE_PASSWORD
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import org.slf4j.LoggerFactory.getLogger
import java.io.File
import java.security.KeyStore
import java.security.KeyStore.getDefaultType
import javax.net.ssl.TrustManager
import javax.net.ssl.TrustManagerFactory
import javax.net.ssl.TrustManagerFactory.getDefaultAlgorithm

object TrustManagerFactory {
    private val logger = getLogger(this.javaClass)

    fun createTrustAllSecurityManager(): Array<TrustManager> = arrayOf(TrustAllSecurityManager())

    fun createTrustManagerWithServerCertificateVerification(): Array<TrustManager> {
        val trustStorePassword = getConfigProperty<String>(EXTERNAL_TRUSTSTORE_PASSWORD)

        val trustStore = KeyStore.getInstance(getDefaultType()).also {
                val trustStoreLocation = getConfigProperty<String>(EXTERNAL_TRUSTSTORE_LOCATION)
                logger.debug("Using external truststore: $trustStoreLocation")

                File(trustStoreLocation).inputStream().use { inputStream ->
                    it.load(inputStream, trustStorePassword.toCharArray())
                }
        }

        val trustManagerFactory = TrustManagerFactory.getInstance(getDefaultAlgorithm()).also { it.init(trustStore) }
        return trustManagerFactory.trustManagers
    }
}
