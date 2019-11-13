package reddit

import (
	"github.com/pkg/errors"
	"github.com/turnage/graw/reddit"
	"strings"
	"sync"
)

var (
	//ErrMalformedParam is returned when a parameter is not reasonable.
	ErrMalformedParam = errors.New("err_malformed_param")
)

// Reddtior is an interface for interacting with the reddit api.
type Reddtior interface {
	Listing(path, after string) (reddit.Harvest, error)
}

// Interactor is used for interfacing with reddit.
type Interactor struct {
	reddit Reddtior
}

// Represents the types of Reddit entities supported by the bot.
const (
	Comment = "COMMENT"
	POST    = "POST"
)

// Response Represents a response a user can get from reddit.
type Response struct {
	Author   string
	Body     string
	PostedAt uint64
	Type     string
}

// SafeResponse should be used if you are going to use this package concurrently to scan multiple subreddits
// at once.
type SafeResponse struct {
	Res []*Response
	Mux *sync.Mutex
}

// Config represents all the things needed to create a reddit bot.
// All should be supplied.
type Config struct {
	userAgent    string
	clientID     string
	clientSecret string
	username     string
	password     string
}

// NewSafeResponse returns a SafeResponse with some defaults.
func NewSafeResponse() *SafeResponse {
	return &SafeResponse{
		Res: []*Response{},
		Mux: &sync.Mutex{},
	}
}

// NewConfig does some basic validation on parameters before returning a new Config used to
// create a Interactor. Can return an error if they are not correct.
func NewConfig(userAgent string,
	clientID string,
	clientSecret string,
	username string,
	password string,
) (*Config, error) {

	switch {
	case clientID == "":
		return nil, errors.Wrap(ErrMalformedParam, "clientID")
	case clientSecret == "":
		return nil, errors.Wrap(ErrMalformedParam, "clientSecret")
	case username == "":
		return nil, errors.Wrap(ErrMalformedParam, "username")
	case password == "":
		return nil, errors.Wrap(ErrMalformedParam, "password")
	}

	return &Config{userAgent: userAgent,
		clientID:     clientID,
		clientSecret: clientSecret,
		username:     username,
		password:     password,
	}, nil
}

// NewInteractor is used to create a new reddit Interactor.
// It returns an error if a new reddit bot cannot be created.
func NewInteractor(c Config) (*Interactor, error) {
	b := reddit.BotConfig{
		Agent: c.userAgent,
		App: reddit.App{
			ID:       c.clientID,
			Secret:   c.clientSecret,
			Username: c.username,
			Password: c.password,
		},
		Rate: 0,
	}
	r, err := reddit.NewBot(b)
	if err != nil {
		return nil, err
	}

	return &Interactor{reddit: r}, nil
}

//Search should take a subreddit of interest and a channel(?) and return into that channel any matches.
func (i Interactor) Search(subredditToMonitor string, keyPhrasesToMonitor ...string) ([]*Response, error) {
	var res []*Response

	h, err := i.reddit.Listing(subredditToMonitor, "")
	if err != nil {
		return nil, err
	}

	for _, p := range h.Posts {
		for _, k := range keyPhrasesToMonitor {
			if strings.Contains(strings.ToLower(p.SelfText), strings.ToLower(k)) ||
				strings.Contains(strings.ToLower(p.Title), strings.ToLower(k)) {
				res = append(res, &Response{
					Author:   p.Author,
					Body:     p.SelfText,
					PostedAt: p.CreatedUTC,
					Type:     POST,
				})
			}
		}
	}

	for _, c := range h.Comments {
		for _, k := range keyPhrasesToMonitor {
			if strings.Contains(strings.ToLower(c.Body), strings.ToLower(k)) {
				res = append(res, &Response{
					Author:   c.Author,
					Body:     c.Body,
					PostedAt: c.CreatedUTC,
					Type:     Comment,
				})
			}
		}
	}
	return res, err
}

// ValidateSubreddits is used to take phrases from the command line
// and ensure they are in the format expected by the Reddit Interactor.
func ValidateSubreddits(subs string) ([]string, error) {
	sp := strings.Split(subs, ",")
	if len(sp) == 0 {
		return nil, errors.New("err_no_subs")
	}
	for _, v := range sp {
		if !strings.HasPrefix(v, "/r/") {
			return nil, errors.New("err_not_a_valid_sub")
		}
	}
	return sp, nil
}

// ValidatePhrases is used to take phrases from the command line
// and ensure they are in the format expected by the Reddit Interactor.
func ValidatePhrases(phrases string) ([]string, error) {
	sp := strings.Split(phrases, ",")
	if len(sp) == 0 {
		return nil, errors.New("err_no_subs")
	}
	return sp, nil
}
