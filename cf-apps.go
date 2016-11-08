package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"strings"
)

var (
	CF_URL          string
	OAUTH_TOKEN     string
	SPACES_ENDPOINT string = "/v2/spaces"
	APPS_WG         sync.WaitGroup
)

type Resource struct {
	Metadata struct {
		Guid string `json:"guid"`
	} `json:"metadata"`
	Entity struct {
		Name string `json:"name"`
	} `json:"entity"`
}

type App struct {
	Guid  string `json:"guid"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type AppDetails struct {
	Stats struct {
		Usage struct {
			Disk int64   `json:"disk"`
			Mem  int64   `json:"mem"`
			Cpu  float64 `json:"cpu"`
			Time string  `json:"time"`
		} `json:"usage"`
		Host string `json:"host"`
		Port int    `json:"port"`
		Uris []string	`json:"uris"`
	} `json:"stats"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: cf-apps https://cloudfoundry.url \"$(cf oauth-token)\"")
	}
	CF_URL = os.Args[1]
	OAUTH_TOKEN = os.Args[2]

	body, err := connectCF(SPACES_ENDPOINT)
	if err != nil {
		log.Fatal(err)
	}
	spaces, err := fetchSpaces(body)
	if err != nil {
		log.Fatal(err)
	}

	for spaceName, guid := range spaces {
		body, err := connectCF(SPACES_ENDPOINT + "/" + guid + "/summary")
		if err != nil {
			log.Fatal(err)
		}
		apps, err := fetchApps(body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(spaceName)
		displayAppDetails(apps)
	}
}

func connectCF(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", CF_URL+url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", OAUTH_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return ioutil.ReadAll(resp.Body)
}

func displayAppDetails(apps []App) {
	APPS_WG.Add(len(apps))
	for _, application := range apps {
		go displayApplicationDetails(application)
	}
	APPS_WG.Wait()
}

func displayApplicationDetails(application App) {
	defer APPS_WG.Done()
	if application.State == "STARTED" {
		endpoint := "/v2/apps/" + application.Guid + "/stats"
		body, err := connectCF(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		details, err := fetchAppDetails(body)
		if err != nil {
			log.Fatal(err)
		}
		for _, detailedApp := range details {
			fmt.Println("\t" + application.Name + "-> " + detailedApp.Stats.Host + ":" + strconv.Itoa(detailedApp.Stats.Port) + ", Uris: " + strings.Join(detailedApp.Stats.Uris, ", "))
		}
	}
}

func fetchSpaces(body []byte) (map[string]string, error) {
	if strings.Contains(string(body), "Invalid Auth Token") {
		log.Fatal("Invalid Auth Token")
	}

	var resources struct {
		Resources []Resource `json:"resources"`
	}
	err := json.Unmarshal(body, &resources)

	spaces := make(map[string]string)
	for _, res := range resources.Resources {
		spaces[res.Entity.Name] = res.Metadata.Guid
	}
	if err != nil {
		log.Println(string(body))
		log.Fatal(err)
	}
	return spaces, err
}

func fetchApps(body []byte) ([]App, error) {
	var response struct {
		Apps []App `json:"apps"`
	}
	err := json.Unmarshal(body, &response)

	if err != nil {
		log.Println(string(body))
		log.Fatal(err)
	}
	return response.Apps, err
}

func fetchAppDetails(body []byte) (map[string]AppDetails, error) {
	results := map[string]AppDetails{}
	err := json.Unmarshal(body, &results)
	if err != nil {
		log.Println(string(body))
		log.Fatal(err)
	}
	return results, err
}
