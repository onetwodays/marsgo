package httpclient

import (
    resty  "github.com/go-resty/resty/v2"
    "time"
)

type HttpClient struct {
    urlBase string
    Client *resty.Client
}

// Note: 工厂函数生成一个对象
func NewHttpClient(urlbase string) *HttpClient  {
    this:=&HttpClient{
        urlBase:urlbase,
        Client:resty.New(),
    }
    if len(urlbase)>0{
        this.Client.SetHostURL(this.urlBase)
    }
    return this
}
func (this *HttpClient)  SetHeaders(headers map[string]string)  {
    // SetHeader("Accept", "application/json").
    this.Client.SetHeaders(headers)
}
func (this *HttpClient)   SetAuthToken(authToken string)  {
    this.Client.SetAuthToken(authToken)
}
func (this *HttpClient)   SetTimeout(timeout time.Duration)  {
    this.Client.SetTimeout(timeout)
}

func (this *HttpClient) Get(url string) (*resty.Response,error) {

    return this.Client.R().Get(url)
}

func (this *HttpClient) GetByQueryParams(url string, queryParams map[string] string) (*resty.Response,error) {

    return this.Client.R().SetQueryParams(queryParams).Get(url)
}

//  SetQueryParams

func (this *HttpClient) Post(url string, httpBody *string) (*resty.Response,error){
    return this.Client.R().SetBody(httpBody).Post(url)
}
func (this *HttpClient) PostMap(url string, mapBody map[string] interface{}) (*resty.Response,error) {
    return this.Client.R().SetBody(mapBody).Post(url)
}






