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

package com.odim.simulator.repo.bmc.utils

import io.javalin.http.UploadedFile
import java.io.ByteArrayInputStream
import java.io.File
import javax.servlet.MultipartConfigElement
import javax.servlet.http.HttpServletRequest

fun getFileMetadata(servletRequest: HttpServletRequest): UploadedFile {
    servletRequest.setAttribute("org.eclipse.jetty.multipartConfig", MultipartConfigElement(System.getProperty("java.io.tmpdir")))
    return servletRequest.parts.iterator().next().run {
        UploadedFile(
                contentType = contentType,
                content = ByteArrayInputStream(inputStream.readBytes()),
                filename = submittedFileName,
                size = size,
                contentLength = size.toInt(),
                extension = submittedFileName.replaceBeforeLast(".", "")
        )
    }
}

fun uploadFile(uploadedName: String, uploadedFileBytes: ByteArray?) {
    val byteArrayInputStream = ByteArrayInputStream(uploadedFileBytes)
    byteArrayInputStream.use { inputStream ->
        // Output file will be saved into Working Directory location.
        val nameOnFilesystem = uploadedName.plus("_out")
        File(nameOnFilesystem).outputStream().use { outputStream ->
            inputStream.copyTo(outputStream)
        }
    }
}

fun filenameWithoutExtension(uploadedFile: UploadedFile) = uploadedFile.filename.substringBeforeLast(uploadedFile.extension)
