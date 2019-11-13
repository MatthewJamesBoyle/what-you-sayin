package main

import (
	"flag"
	"fmt"
	"github.com/matthewjamesboyle/whattheysayingbot/internal/reddit"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	//autoload is used for loading .env files locally.
	_ "github.com/joho/godotenv/autoload"
)

//TODO: add tests
//TODO: add documentation
//TODO: push it to github
// TODO: Add CI
// TODO: tag for release
func main() {

	var (
		useragent      string
		clientid       string
		clientSecret   string
		redditUsername string
		redditPassword string
	)

	for k, v := range map[string]*string{
		"USER_AGENT":      &useragent,
		"CLIENT_ID":       &clientid,
		"CLIENT_SECRET":   &clientSecret,
		"REDDIT_USERNAME": &redditUsername,
		"REDDIT_PASSWORD": &redditPassword,
	} {
		var ok bool
		if *v, ok = os.LookupEnv(k); !ok {
			log.Fatalf("Mising environment variable %s", k)
		}

	}

	subsCli := flag.String("subreddits", "", "A comma separated list of all the subreddits you are interested in.")

	phrasesCli := flag.String("phrases", "", "A comma separated list of all the phrases you are interested in.")

	flag.Parse()

	if *subsCli == "" {
		log.Fatalf("you need to provide atleast one subreddit.")
	}

	if *phrasesCli == "" {
		log.Fatalf("you need to provide atleast one subreddit.")
	}

	phrases, err := reddit.ValidatePhrases(*phrasesCli)
	if err != nil {
		log.Fatal(err)
	}

	subs, err := reddit.ValidateSubreddits(*subsCli)
	if err != nil {
		log.Fatal(err)
	}

	//look up how we do this.
	c, err := reddit.NewConfig(
		useragent,
		clientid,
		clientSecret,
		redditUsername,
		redditPassword,
	)
	if err != nil {
		log.Fatal(err)
	}

	r, err := reddit.NewInteractor(*c)
	if err != nil {
		log.Fatal(err)
	}

	var g errgroup.Group
	rsp := reddit.NewSafeResponse()

	for _, v := range subs {
		g.Go(func() error {
			res, err := r.Search(v, phrases...)
			if err != nil {
				return err
			}
			rsp.Mux.Lock()
			rsp.Res = append(rsp.Res, res...)
			rsp.Mux.Unlock()
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		log.Fatal("failed to get some data")
	}

	for _, v := range rsp.Res {
		fmt.Println(v.Body)
	}

}
