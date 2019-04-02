package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	kubeAPIToken string
	svcURL       string
)

type Service struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Labels struct {
			App        string `json:"app"`
			LastUpdate string `json:"lastUpdate"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		ClusterIP string `json:"clusterIP"`
		Ports     []struct {
			Name       string `json:"name"`
			Port       int    `json:"port"`
			Protocol   string `json:"protocol"`
			TargetPort int    `json:"targetPort"`
		} `json:"ports"`
		Selector struct {
			App string `json:"app"`
		} `json:"selector"`
		SessionAffinity string `json:"sessionAffinity"`
		Type            string `json:"type"`
	} `json:"spec"`
}

func main() {
	if err := appInit(); err != nil {
		log.Printf("[error] %v\n", err.Error())
		os.Exit(1)
	}

	log.Println("[info] Starting cron job...")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log.Printf("[info] GET %s\n", svcURL)

	req, _ := http.NewRequest("GET", svcURL, nil)
	req.Header.Add(`Authorization`, fmt.Sprintf("Bearer %s", kubeAPIToken))
	req.Header.Add(`Accept`, `application/json`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		os.Exit(1)
	}

	svc := Service{}

	if err = json.NewDecoder(resp.Body).Decode(&svc); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	svc.Metadata.Labels.LastUpdate = strconv.FormatInt(time.Now().Unix(), 10)

	payload, err := json.Marshal(&svc)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	req, _ = http.NewRequest("PUT", svcURL, bytes.NewBuffer(payload))
	req.Header.Add(`Authorization`, fmt.Sprintf("Bearer %s", kubeAPIToken))
	req.Header.Add(`Content-Type`, `application/json`)

	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[error] %v\n", err.Error())
			log.Printf("[error] could not read response, status was: %v\n", resp.StatusCode)
			os.Exit(1)
		}
		log.Printf("[error] HTTP PUT code %v\n", string(bodyBytes))
		log.Printf("[error] HTTP PUT body '%v'\n", resp.StatusCode)
		os.Exit(1)
	}

	log.Printf("[info] service updated at %v\n", time.Now().Format(time.RFC3339))
}

func appInit() error {
	missing := []string{}

	kubeAPIToken = os.Getenv("KUBE_API_TOKEN")
	if len(kubeAPIToken) == 0 {
		missing = append(missing, "KUBE_API_TOKEN")
	}

	kubeAPIUrl := os.Getenv("KUBE_API_URL")
	if len(kubeAPIUrl) == 0 {
		missing = append(missing, "KUBE_API_URL")
	}

	svcURLPath := os.Getenv("SVC_URL_PATH")
	if len(svcURLPath) == 0 {
		missing = append(missing, "SVC_URL_PATH")
	}

	svcURL = fmt.Sprintf("%s%s", kubeAPIUrl, svcURLPath)

	if len(missing) != 0 {
		return fmt.Errorf("Missing environment variables: %v", missing)
	}

	return nil
}
