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

package com.odim.utils

fun ipAddressOf(value: String): IpV4Address {
    val ipAddressParts = value.split(".")
    require(ipAddressParts.size == 4) { "$value does not represent IPv4 Address" }

    val bytes = ipAddressParts.map {
        it.toIntInRange(0, 255) ?: throw IllegalArgumentException("$it is not a number in range 0..255")
    }.toIntArray()

    return IpV4Address(bytes[0], bytes[1], bytes[2], bytes[3])
}

fun String.toIpV4Address() = ipAddressOf(this)

fun ClosedRange<String>.toIpV4AddressRange() = IpV4AddressRange(
    this.start.toIpV4Address(),
    this.endInclusive.toIpV4Address()
)

data class IpV4Address(val b3: Int, val b2: Int, val b1: Int, val b0: Int) : Comparable<IpV4Address> {
    override fun toString() = "$b3.$b2.$b1.$b0"

    override fun compareTo(other: IpV4Address) = when {
        b3 != other.b3 -> b3 - other.b3
        b2 != other.b2 -> b2 - other.b2
        b1 != other.b1 -> b1 - other.b1
        else -> b0 - other.b0
    }

    operator fun rangeTo(that: IpV4Address) = IpV4AddressRange(this, that)

    operator fun inc(): IpV4Address {
        fun Int.incIf(condition: Boolean) = if (condition) (this + 1) % 256 else this

        val r0 = b0.incIf(true)
        val r1 = b1.incIf(r0 == 0)
        val r2 = b2.incIf(r0 == 0 && r1 == 0)
        val r3 = b3.incIf(r0 == 0 && r1 == 0 && r2 == 0)

        return IpV4Address(r3, r2, r1, r0)
    }
}

data class IpV4AddressRange(
    override val start: IpV4Address,
    override val endInclusive: IpV4Address
) : ClosedRange<IpV4Address>, Iterable<IpV4Address>, Sequence<IpV4Address> {

    override fun iterator(): Iterator<IpV4Address> = IpV4AddressIterator(start)

    inner class IpV4AddressIterator(private var current: IpV4Address) : Iterator<IpV4Address> {
        override fun hasNext() = current <= endInclusive
        override fun next() = if (hasNext()) current++ else throw NoSuchElementException()
    }
}

private fun String.toIntInRange(min: Int, max: Int) = this.toInt().takeIf { it in min..max }
