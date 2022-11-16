#(C) Copyright [2022] American Megatrends International LLC
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.

# EVENT_DESTINATION_URL = "/ODIM/v1/EventService/CompositionServiceEvent"
"""
EVENT_TYPES = [
    "ResourceRemoved", "ResourceAdded", "ResourceUpdated", "StatusChange",
    "Alert", "MetricReport", "Other"
]
"""

# EVENT_SUBSCRIPTION_URL = "/redfish/v1/EventService/Subscriptions"


RESOURCE_BLOCK_TEMP = {
    "@odata.context": "/redfish/v1/$metadata#ResourceBlock.ResourceBlock.",
    "@odata.type": "#ResourceBlock.v1_4_0.ResourceBlock",
    "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks",
    "Id": "ComputerSystemBlock",
    "Name": "ComputerSystem Block",
    "Description": "ComputerSystem Block",
    "ResourceBlockType": [],
    "Status": {
        "State": "Enabled",
        "Health": "OK"
    },
    "CompositionStatus": {
        "Reserved": False,
        "CompositionState": "Unused",
        "SharingCapable": False,
        "MaxCompositions": 1,
        "NumberOfCompositions": 0
    },
    "Client": None,
    "Pool": "Free"
}

RESOURCE_ZONE_TEMP = {
    "@odata.type": "#Zone.v1_6_0.Zone",
    "Id": "1",
    "Name": "Resource Zone 1",
    "Status": {
        "State": "Enabled",
        "Health": "OK"
    },
    "ZoneType": "ZoneOfResourceBlocks",
    "Links": {
        "ResourceBlocks": []
    },
    "@odata.id": "/redfish/v1/CompositionService/ResourceZones"
}

COMPOSITION_SERVICE_NAME = "svc.composition.service"

ACCOUNT_SESSION_NAME = "svc.account.session"

# PrivilegeLogin defines the privilege for login
PrivilegeLogin = "Login"
# PrivilegeConfigureManager defines the privilege for configuraton manager
PrivilegeConfigureManager = "ConfigureManager"
# PrivilegeConfigureUsers defines the privilege for user configuratons
PrivilegeConfigureUsers = "ConfigureUsers"
# PrivilegeConfigureSelf defines the privilege for self configuratons
PrivilegeConfigureSelf = "ConfigureSelf"
# PrivilegeConfigureComponents defines the privilege for component configuratons
PrivilegeConfigureComponents = "ConfigureComponents"