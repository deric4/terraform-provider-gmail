package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var oauthScopes = []string{
	gmail.GmailLabelsScope,
}

// Config is the structure used to instantiate the Gmail provider.
type Config struct {
	gmail *gmail.Service
}

func (c *Config) loadAndValidate() error {
	// read file []byte
	b, err := ioutil.ReadFile("/Users/miguel/Downloads/credentials.json")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	config, err := google.ConfigFromJSON(b, oauthScopes...)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	client := getClient(config)

	//client, err := google.DefaultClient(context.Background(), oauthScopes...)
	//if err != nil {
	//	return err
	//}

	client.Transport = logging.NewTransport("Google", client.Transport)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraform.VersionString())

	// Create gmail service.
	gmailSvc, err := gmail.New(client)
	if err != nil {
		log.Printf("%s", err)
		return nil
	}

	gmailSvc.UserAgent = userAgent
	c.gmail = gmailSvc

	return nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "/Users/miguel/.token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	time.Sleep(30 * time.Second)

	code, err := ioutil.ReadFile("/Users/miguel/.auth_code")
	if err != nil {
		return nil
	}

	tok, err := config.Exchange(context.TODO(), string(code[:]))
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
