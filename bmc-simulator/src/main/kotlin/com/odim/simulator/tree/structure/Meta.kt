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

package com.odim.simulator.tree.structure

import com.odim.simulator.tree.RedfishVersion
import com.odim.simulator.tree.ResourceFactory

data class Meta(val type: Type,
                var parent: TreeElement? = null,
                var createdVersion: RedfishVersion? = null,
                var resourceFactory: ResourceFactory? = null)
