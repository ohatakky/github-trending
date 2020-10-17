package tweet

import (
	"github.com/ChimeraCoder/anaconda"
)

type Client struct {
	service *anaconda.TwitterApi
}

type Tweet struct {
	Text      string
	CreatedAt string
}

func New(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Client {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	return &Client{api}
}

func (c *Client) Tweet(message string) (*Tweet, error) {
	tweet, err := c.service.PostTweet(message, nil)
	return &Tweet{Text: tweet.Text, CreatedAt: tweet.CreatedAt}, err
}
