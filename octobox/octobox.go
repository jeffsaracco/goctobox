package octobox

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	//Client holds the information of the octobox installation
	Client struct {
		InstanceAddr string
		APIToken     string
	}

	//NotificationSubject defines the subject structure in the response data
	NotificationSubject struct {
		Title string `json:"title"`
		Type  string `json:"type"`
	}

	//NotificationRepo defines the repo structure in the response data
	NotificationRepo struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Owner string `json:"owner"`
		URL   string `json:"repo_url"`
	}

	//Notification defines the return from the notifications api endpoint
	Notification struct {
		ID       int                 `json:"id"`
		GitHubID int                 `json:"github_id"`
		Reason   string              `json:"reason"`
		Unread   bool                `json:"unread"`
		WebURL   string              `json:"web_url"`
		Subject  NotificationSubject `json:"subject"`
		Repo     NotificationRepo    `json:"repo"`
	}

	//APIResponse holds the root level response from the API
	APIResponse struct {
		Notifications []Notification `json:"notifications"`
	}
)

//New returns a new, configured Client
func New(instanceAddr, token string) *Client {
	return &Client{
		InstanceAddr: instanceAddr,
		APIToken:     token,
	}
}

//GetNotifications returns a list of the notifications
func (c *Client) GetNotifications() *[]Notification {
	req, _ := http.NewRequest("GET", c.InstanceAddr+"/notifications.json", nil)
	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	rs, _ := http.DefaultClient.Do(req)
	// if err != nil {
	// panic(err) // More idiomatic way would be to print the error and die unless it's a serious error
	// }
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println("DATA ", c.APIToken, c.InstanceAddr, string(bodyBytes))

	data := APIResponse{}

	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		panic(err)
	}

	return &data.Notifications
}
