package tweet

import (
	"github.com/ChimeraCoder/anaconda"
)

type Client struct {
	service *anaconda.TwitterApi
}

func New(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Client {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	return &Client{api}
}

func (c *Client) Tweet(message string) error {
	_, err := c.service.PostTweet(message, nil)
	return err
}
