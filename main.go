// +build go1.5

package main

import (
	"fmt"
	"os"
	// "path"

	gottyclient "github.com/moul/gotty-client"
	"github.com/sirupsen/logrus"
	// "github.com/urfave/cli"
	"net/http"
	"io/ioutil"
	"strings"
    "encoding/json"
)

func main() {

	argsWithoutProg := os.Args[1:]
	
	fmt.Printf("Launching Binder: %s\n", argsWithoutProg[0])

	binder_launch := fmt.Sprintf("https://mybinder.org/build/gh/%s/master",argsWithoutProg[0])
	resp, err := http.Get(binder_launch)
	if err != nil {
		logrus.Fatalf("Communication error")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("Communication error")
	}

	messages := strings.Split(string(body),"\n")

	for i := 0; i < len(messages);  i++ {
		x := strings.SplitN(messages[i],"data: ",2)
		if len(x) == 2 {

			var result map[string]interface{}
			json.Unmarshal([]byte(x[1]), &result)
		    value, ok := result["phase"].(string)
			if !ok {
				logrus.Fatalf("Communication error")
			}

			if value == "ready" {
				url, ok := result["url"].(string)
				if !ok {
					logrus.Fatalf("Communication error")
				}

				token, ok := result["token"].(string)
				if !ok {
					logrus.Fatalf("Communication error")
				}

				fullurl := fmt.Sprintf("%sproxy/8080/?token=%s", url, token)

				client, err := gottyclient.NewClient(fullurl)
				if err != nil {
					logrus.Fatalf("Cannot create client: %v", err)
				}
				client.V2 = true

				if err = client.Loop(); err != nil {
					logrus.Fatalf("Communication error: %v", err)
				}
							
			}

		}
	}
}