package anedya

import (
	"net/http"
	"time"

	accesstokens "github.com/anedyaio/anedya-go-sdk/accessTokens"
	"github.com/anedyaio/anedya-go-sdk/dataAccess"
	"github.com/anedyaio/anedya-go-sdk/deviceLogs"
	"github.com/anedyaio/anedya-go-sdk/nodes"
	valuestore "github.com/anedyaio/anedya-go-sdk/valueStore"
	"github.com/anedyaio/anedya-go-sdk/variable"
)

type AnedyaRegion string

const (
	AP_IN_1 AnedyaRegion = "ap-in-1"
)

type Client struct {
	NodeManagement        *nodes.NodeManagement
	VariableManagement    *variable.VariableManagement
	DataManagement        *dataAccess.DataManagement
	AccessTokenManagement *accesstokens.AccessTokenManagement
	DeviceLogManagement   *deviceLogs.DeviceLogManagement
	ValueStoreManagement  *valuestore.ValueStoreManagement
}

func NewClient(baseURL, apiKey string) *Client {

	auth := &authTransport{
		apiKey: apiKey,
		next:   http.DefaultTransport,
	}

	hc := &http.Client{
		Timeout:   30 * time.Second,
		Transport: auth,
	}

	return &Client{
		NodeManagement:        nodes.NewNodeManagement(hc, baseURL),
		VariableManagement:    variable.NewVariableManagement(hc, baseURL),
		DataManagement:        dataAccess.NewDataManagement(hc, baseURL),
		AccessTokenManagement: accesstokens.NewAccessTokenManagement(hc, baseURL),
		DeviceLogManagement:   deviceLogs.NewDeviceLogManagement(hc, baseURL),
		ValueStoreManagement:  valuestore.NewValueStoreManagement(hc, baseURL),
	}
}

func DefaultURL(region AnedyaRegion) string {
	return "https://api." + string(region) + ".anedya.io"
}
