package client

import (
	"bytes"
	"doccer/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)


// Test client
type Client struct {
	url string
	client http.Client
}

func NewClient(url string) Client {
	return Client {
		url: url,
		client: http.Client {},
	}
}

// Returns user id
func (c* Client) Register(login string, password string) (string, error) {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"login":"%s", "password":"%s"}`, login, password)))
	reqUrl := c.url + "/register"
	resp, err := c.client.Post(reqUrl, "application/json", reqJson)
	if err != nil {
		return "", err
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return c.parseField(bodyBytes, "id")
}

// Returns jwt
func (c* Client) Login(login string, password string) (string, error) {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"login":"%s", "password":"%s"}`, login, password)))
	resp, err := c.client.Post(c.url + "/login", "application/json", reqJson)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return c.parseField(bodyBytes, "Token")
	} else {
		return "", errors.New("Error: " + strconv.Itoa(resp.StatusCode))

	}
}

// Returns docId
func (c* Client) CreateDoc(text string, defaultAccess string, token string) (string, error) {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"id":"-1", "authorId":"-1", "text":"%s", "access":"%s"}`, text, defaultAccess)))

	req, _ := http.NewRequest("POST", c.url + "/docs", reqJson)

	resp, _ := c.makeAuthRequest(*req, token)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return c.parseField(bodyBytes, "id")
}


// Returns groupId
func (c* Client) CreateGroup(name string, token string) (string, error) {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"id":"-1", "creatorId":"-1", "name":"%s"}`, name)))

	req, _ := http.NewRequest("POST", c.url + "/users/groups", reqJson)

	resp, _ := c.makeAuthRequest(*req, token)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return c.parseField(bodyBytes, "id")
}

func (c* Client) AddMember(groupId string, memberId string, token string) error {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"groupId":"%s", "memberId":"%s"}`, groupId, memberId)))

	req, _ := http.NewRequest("PUT", c.url + fmt.Sprintf("/users/groups/%s/members", groupId), reqJson)

	_, _ = c.makeAuthRequest(*req, token)
	return nil
}

func (c* Client) ChangeMemberAccess(docId string, memberId string, newAccess string, token string) error {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"id":"%s", "type":0, "itemId":"%s", "access":"%s"}`, docId, memberId, newAccess)))

	req, _ := http.NewRequest("POST", c.url + fmt.Sprintf("/docs/%s/access", docId), reqJson)

	_, _ = c.makeAuthRequest(*req, token)
	return nil
}

func (c* Client) ChangeGroupAccess(docId string, groupId string, newAccess string, token string) error {
	reqJson := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"id":"%s", "type":1, "itemId":"%s", "access":"%s"}`, docId, groupId, newAccess)))

	req, _ := http.NewRequest("POST", c.url + fmt.Sprintf("/docs/%s/access", docId), reqJson)

	_, _ = c.makeAuthRequest(*req, token)
	return nil
}

func (c* Client) GetDoc(docId string, token string) (*model.Doc, error) {

	req, _ := http.NewRequest("GET", c.url + fmt.Sprintf("/docs/%s", docId), nil)

	resp, _ := c.makeAuthRequest(*req, token)

	var res *model.Doc
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c* Client) makeAuthRequest(request http.Request, token string) (*http.Response, error)  {
	request.Header.Add("AuthToken", token)
	resp, err := c.client.Do(&request)
	return resp, err
}

func (c* Client) parseField(responseBody []byte, fieldName string) (string, error) {
	bodyStr := string(responseBody)
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(bodyStr), &result)
	return result[fieldName].(string), nil
}