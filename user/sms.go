package user

// The client of the sms service of Aliyun(阿理云短信服务客户端)
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Client struct {
	Request    *Request
	GatewayURL string
	Client     *http.Client
}

func New(gatewayURL string) *Client {
	client := new(Client)
	client.Request = &Request{}
	client.GatewayURL = gatewayURL
	client.Client = &http.Client{}
	return client
}

func (client *Client) Execute(accessKeyID, accessKeySecret, phoneNumbers, signName, templateCode, templateParam string) (*Response, error) {
	err := client.Request.SetParamsValue(accessKeyID, phoneNumbers, signName, templateCode, templateParam)
	if err != nil {
		return nil, err
	}
	endpoint, err := client.Request.BuildEndpoint(accessKeySecret, client.GatewayURL)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("GET", endpoint, nil)
	response, err := client.Client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result := new(Response)
	err = json.Unmarshal(body, result)

	result.RawResponse = body
	return result, err
}

// SendCode
func SendCode(phoneNumbers string, code string) error {
	gatewayURL := "http://dysmsapi.aliyuncs.com/"
	accessKeyID := "LTAIogbXuxU4hkrH"
	accessKeySecret := "jCOzTGv6vG6IR3Dm2EgKN6iWX6kXZQ"
	signName := "晁高锋"
	templateCode := "SMS_139981678"
	templateParam := fmt.Sprintf("{\"code\":\"%s\"}", code)
	client := New(gatewayURL)
	result, err := client.Execute(accessKeyID, accessKeySecret, phoneNumbers, signName, templateCode, templateParam)
	if err != nil {
		panic("Failed to send Message: " + err.Error())
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	if result.IsSuccessful() {
		fmt.Println("[SMS] A SMS is sent successfully:", phoneNumbers, string(resultJSON))
		return nil
	} else {
		fmt.Println("[SMS] Failed to send a SMS:", phoneNumbers, string(resultJSON))
		return fmt.Errorf("[SMS] Failed to send a SMS: %s, %s", phoneNumbers, string(resultJSON))
	}
}

// MakeCode 生成验证码
func MakeCode() (code string) {
	code = strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(8999) + 1000)
	return
}

// VailMobile 验证手机号
func VailMobile(mobile string) error {
	if len(mobile) < 11 {
		return errors.New("手机号码位数不正确")
	}
	reg, err := regexp.Compile("^1[3-8][0-9]{9}$")
	if err != nil {
		panic("regexp error")
	}
	if !reg.MatchString(mobile) {
		return errors.New("手机号码格式不正确")
	}
	return nil
}

// VailCode 验证验证码
func VailCode(code string) error {
	if len(code) != 4 {
		return errors.New("验证码位数不正确")
	}
	c, err := regexp.Compile("^[0-9]{4}$")
	if err != nil {
		panic("regexp error")
	}
	if !c.MatchString(code) {
		return errors.New("验证码格式不正确")
	}
	return nil
}
