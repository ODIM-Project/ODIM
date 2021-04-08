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

import com.odim.odimclient.security.TrustManagerFactory.createTrustAllSecurityManager
import com.odim.odimclient.security.TrustManagerFactory.createTrustManagerWithServerCertificateVerification
import org.apache.http.conn.ssl.NoopHostnameVerifier
import org.apache.http.impl.client.HttpClients
import java.security.SecureRandom
import javax.net.ssl.SSLContext

data class ClientSecurityConfigurator(private val serverCertVerificationEnabled: Boolean) {
    private val trustManagers =
            if (serverCertVerificationEnabled) {
                createTrustManagerWithServerCertificateVerification()
            } else {
                createTrustAllSecurityManager()
            }

    private val sslContext = SSLContext.getInstance("TLSv1.2").also {
        it.init(null, trustManagers, SecureRandom())
    }

    val httpClient = HttpClients.custom()
            .setSSLHostnameVerifier(NoopHostnameVerifier()).setSSLContext(sslContext)
            .build()!!
}
