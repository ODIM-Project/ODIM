
EVENT_DESTINATION_URL = "/redfish/v1/EventService/CompositionServiceEvent"

EVENT_TYPES = ["ResourceRemoved", "ResourceAdded", "ResourceUpdated",
               "StatusChange", "Alert", "MetricReport", "Other"]

EVENT_SUBSCRIPTION_URL = "/redfish/v1/EventService/Subscriptions"

SYSTEMS_URL = "/redfish/v1/Systems"

RESOURCE_BLOCK_TEMP = {
    "@odata.context": "/redfish/v1/$metadata#ResourceBlock.ResourceBlock.",
    "@odata.type": "#ResourceBlock.v1_4_0.ResourceBlock",
    "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks",
    "Id": "ComputerSystemBlock",
    "Name": "ComputerSystem Block",
    "Description": "ComputerSystem Block",
    "ResourceBlockType": [
    ],
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
    "PoolType": "Free"
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
        "ResourceBlocks": [
        ]
    },
    "@odata.id": "/redfish/v1/CompositionService/ResourceZones"
}

COMPOSITION_SERVICE_NAME = "svc.composition.service"
