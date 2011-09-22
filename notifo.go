package notifo

import (
    "http"
    "url"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "json"
)

type response struct {
    Status string `json:"status"`
    Code int `json:"response_code"`
    Message string `json:"response_message"`
}

type Client struct {
    Username string
    ApiSecret string
}

type M map[string]string 

func (nc *Client) request(path string, data url.Values) (*response, os.Error) {
    // endpoint
    url := "https://api.notifo.com/v1" + path
    // encode request params
    reqBody := strings.NewReader(data.Encode())
    // build request
    req, err := http.NewRequest("POST", url, reqBody)
    if err != nil {
        return nil, err
    }
    req.Method = "POST"
    req.SetBasicAuth(nc.Username, nc.ApiSecret)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    // make request
    client := new(http.Client)
    r,err := client.Do(req)
    // check connection
    if err != nil {
        return nil,err
    }
    defer r.Body.Close()
    // read response
    body,_ := ioutil.ReadAll(r.Body)
    // decode json
    response := new(response)
    err = json.Unmarshal(body, &response)
    if err != nil {
        return nil,err
    }
    // check for success code
    if response.Code != 2201 {
        return nil,fmt.Errorf("notifo: %s", response.Message)
    }
    // no error
    return response,nil
}

// Register a client
func (nc *Client) Subscribe(username string) os.Error {
    data := url.Values{}
    data.Add("username", username)
    _,err := nc.request("/subscribe_user", data)
    return err
}

// optional params for "options" map are:
//
// label: label describing the "application" (used only if being sent from a User account; the Service label is automatically applied if being sent from a Service account)
// title: name of "notification event
// uri: the uri that will be loaded when the notification is opened; if specified, must be urlencoded; if a web address, must start with http:// or https://
//
func (nc *Client) Send(to string, msg string, options... M) os.Error {
    data := url.Values{}
    data.Add("to", to)
    data.Add("msg", msg)
    // grab optional parameters
    for _,option := range options {
        for k,v := range option {
            data.Add(k,v)
        }
    }
    _,err := nc.request("/send_notification", data)
    return err
}
