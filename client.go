package go5paisa

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	// "log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	baseURL string = "https://Openapi.5paisa.com/VendorsAPI/Service1.svc"

	loginRoute          string = "/V3/LoginRequestMobileNewbyEmail"
	marginRoute         string = "/V3/Margin"
	orderBookRoute      string = "/V2/OrderBook"
	holdingsRoute       string = "/V2/Holding"
	positionsRoute      string = "/V1/NetPositionNetWise"
	orderPlacementRoute string = "/V1/OrderRequest"
	orderStatusRoute    string = "/OrderStatus"
	tradeInfoRoute      string = "/TradeInformation"

	// Request codes
	marginRequestCode         string = "5PMarginV3"
	orderBookRequestCode      string = "5POrdBkV2"
	holdingsRequestCode       string = "5PHoldingV2"
	positionsRequestCode      string = "5PNPNWV1"
	tradeInfoRequestCode      string = "5PTrdInfo"
	orderStatusRequestCode    string = "5POrdStatus"
	orderPlacementRequestCode string = "5POrdReq"
	loginRequestCode          string = "5PLoginV3"

	// Content Type
	contentType string = "application/json"
)

// Config is the app configuration
type Config struct {
	AppName       string
	AppSource     string
	UserID        string
	Password      string
	UserKey       string
	EncryptionKey string
}

// AppConfig is a reusable config struct
type AppConfig struct {
	config *Config
	head   *payloadHead
}

//Client is the client configuration
type Client struct {
	clientCode string
	connection *http.Client
	appConfig  *AppConfig
}

// Init initializes the AppConfig struct
func Init(c *Config) *AppConfig {
	head := &payloadHead{
		AppName:     c.AppName,
		AppVer:      "1.0",
		Key:         c.UserKey,
		OsName:      "WEB",
		RequestCode: "",
		UserID:      c.UserID,
		Password:    c.Password,
	}
	appConfig := &AppConfig{
		config: c,
		head:   head,
	}
	return appConfig
}

//Login logs in a client
func Login(conf *AppConfig, encryptedEmail string, encryptedPassword string, encryptedDOB string) (*Client, error) {
	//encryptedEmail := encrypt(conf.config.EncryptionKey, email)
	//encryptedPassword := encrypt(conf.config.EncryptionKey, password)
	//encryptedDOB := encrypt(conf.config.EncryptionKey, dob)
	var client *Client
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return client, err
	}
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
	loginRequestBody := loginBody{
		Email:          encryptedEmail,
		Password:       encryptedPassword,
		LocalIP:        "192.168.1.1",
		PublicIP:       "192.168.1.1",
		SerialNumber:   "",
		MAC:            "",
		MachineID:      "039377",
		VersionNo:      "1.7",
		RequestNo:      "1",
		My2PIN:         encryptedDOB,
		ConnectionType: "1",
	}
	conf.head.RequestCode = loginRequestCode
	loginDetails := loginPayload{
		Head: conf.head,
		Body: loginRequestBody,
	}
	jsonValue, _ := json.Marshal(loginDetails)
	res, err := httpClient.Post(baseURL+loginRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return client, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return client, err
	}
	var b body
	parseResBody(resBody, &b)
	if b.ClientCode == "" || b.ClientCode == "INVALID CODE" {
		return client, errors.New(b.Message)
	}
	client = &Client{
		clientCode: b.ClientCode,
		connection: httpClient,
		appConfig:  conf,
	}
	return client, nil
}
