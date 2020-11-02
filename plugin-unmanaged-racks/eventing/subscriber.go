package eventing

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"net/http"
	"time"

	"net/url"
)

type Subscriber interface {
	Run()
}

func NewSubscriber(config *config.PluginConfig) Subscriber {

	baseURL, err := url.Parse("https://" + config.EventConf.ListenerHost + ":" + config.EventConf.ListenerPort + "/")
	if err != nil {
		panic(err)
	}

	destURI, err := url.Parse(config.EventConf.DestURI)
	if err != nil {
		panic(err)
	}

	destURL := baseURL.ResolveReference(destURI)
	if err != nil {
		panic(err)
	}
	return &subscriber{
		destinationURL: *destURL,
		odimRAClient:   redfish.NewClient(config.OdimraNBUrl),
	}
}

type subscriber struct {
	odimRAClient   redfish.Client
	destinationURL url.URL
}

func (s *subscriber) Run() {
	go func() {
		for true {
			s.subscribe()
			time.Sleep(time.Second * 15)
		}
	}()
}

func (s *subscriber) subscribe() (*int, error) {
	sr := createSubscriptionRequest(s.destinationURL.String())
	bodyBytes, err := json.Marshal(&sr)
	if err != nil {
		return nil, err
	}

	rsp, err := s.odimRAClient.Post("/redfish/v1/EventService/Subscriptions", bodyBytes)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode == http.StatusAccepted {
		monitor := func() (*http.Response, error) {
			return s.odimRAClient.Get(rsp.Header.Get("Location"))
		}
		for {
			r, e := monitor()
			if e != nil {
				return nil, e
			}
			if r.StatusCode != http.StatusAccepted {
				return &r.StatusCode, nil
			}
		}
	}
	return &rsp.StatusCode, nil

}

func createSubscriptionRequest(destination string) redfish.EvtSubPost {
	return redfish.EvtSubPost{
		Name:                 "URP",
		Destination:          destination,
		EventTypes:           []string{"ResourceRemoved"},
		MessageIds:           nil,
		ResourceTypes:        []string{"ComputerSystem"},
		Context:              "ODIMRA_Event",
		Protocol:             "Redfish",
		SubscriptionType:     "RedfishEvent",
		EventFormatType:      "Event",
		SubordinateResources: true,
		OriginResources: []redfish.Link{
			{
				Oid: "/redfish/v1/Systems",
			},
		},
	}
}
