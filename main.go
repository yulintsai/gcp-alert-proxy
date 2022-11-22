package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

var (
	StartMsg = "GCP Alert Proxy Serving😉"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"w"},
			Usage:   "gcp",
			Action: func(c *cli.Context) error {
				fmt.Println("GCP Alert Proxy is Starting ...")

				authUser := c.Args().First()
				if authUser == "" {
					fmt.Println("[warn] AuthUser is empty")
				}

				authPassword := c.Args().Get(1)
				if authPassword == "" {
					fmt.Println("[warn] AuthPassword is empty")
				}

				listenPort := c.Args().Get(2)
				if listenPort == "" {
					fmt.Println("[warn] listenPort is empty")
				}

				runGin(authUser, authPassword, listenPort)

				return nil
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func runGin(authUser, authPassword, listenPort string) {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "GCP Alert Proxy Pong",
		})
	})

	authorized := r.Group("/gcp", gin.BasicAuth(gin.Accounts{
		authUser: authPassword,
	}))
	authorized.POST("/tg/:bot_token/:chat_id", ToTelegram)

	r.Run(":" + listenPort)
}

type GCPWebHook struct {
	Incident struct {
		IncidentID              string `json:"incident_id"`
		ScopingProjectID        string `json:"scoping_project_id"`
		ScopingProjectNumber    int    `json:"scoping_project_number"`
		URL                     string `json:"url"`
		StartedAt               int    `json:"started_at"`
		EndedAt                 int    `json:"ended_at"`
		State                   string `json:"state"`
		ResourceID              string `json:"resource_id"`
		ResourceName            string `json:"resource_name"`
		ResourceDisplayName     string `json:"resource_display_name"`
		ResourceTypeDisplayName string `json:"resource_type_display_name"`
		Resource                struct {
			Type   string `json:"type"`
			Labels struct {
				InstanceID string `json:"instance_id"`
				ProjectID  string `json:"project_id"`
				Zone       string `json:"zone"`
			} `json:"labels"`
		} `json:"resource"`
		Metric struct {
			Type        string `json:"type"`
			DisplayName string `json:"displayName"`
			Labels      struct {
				InstanceName string `json:"instance_name"`
			} `json:"labels"`
		} `json:"metric"`
		Metadata struct {
			SystemLabels struct {
				Labelkey string `json:"labelkey"`
			} `json:"system_labels"`
			UserLabels struct {
				Labelkey string `json:"labelkey"`
			} `json:"user_labels"`
		} `json:"metadata"`
		PolicyName       string `json:"policy_name"`
		PolicyUserLabels struct {
			UserLabel1 string `json:"user-label-1"`
			UserLabel2 string `json:"user-label-2"`
		} `json:"policy_user_labels"`
		ConditionName  string `json:"condition_name"`
		ThresholdValue string `json:"threshold_value"`
		ObservedValue  string `json:"observed_value"`
		Condition      struct {
			Name               string `json:"name"`
			DisplayName        string `json:"displayName"`
			ConditionThreshold struct {
				Filter       string `json:"filter"`
				Aggregations []struct {
					AlignmentPeriod  string `json:"alignmentPeriod"`
					PerSeriesAligner string `json:"perSeriesAligner"`
				} `json:"aggregations"`
				Comparison     string  `json:"comparison"`
				ThresholdValue float64 `json:"thresholdValue"`
				Duration       string  `json:"duration"`
				Trigger        struct {
					Count int `json:"count"`
				} `json:"trigger"`
			} `json:"conditionThreshold"`
		} `json:"condition"`
		Documentation struct {
			Content  string `json:"content"`
			MimeType string `json:"mime_type"`
		} `json:"documentation"`
		Summary string `json:"summary"`
	} `json:"incident"`
	Version string `json:"version"`
}

func ToTelegram(c *gin.Context) {
	body, ioReadAllErr := ioutil.ReadAll(c.Request.Body)
	if ioReadAllErr != nil {
		fmt.Println("[warn] ioReadAllErr:", ioReadAllErr)
		c.Data(400, "", []byte(ioReadAllErr.Error()))
	}

	ret := &GCPWebHook{}

	jerr := json.Unmarshal(body, &ret)
	if jerr != nil {
		c.Data(400, "", []byte(jerr.Error()))
	}

	token := c.Param("bot_token")
	id := c.Param("chat_id")
	telUrl := "https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + id + "&text="
	msg := "[warn] \n" + ret.Incident.PolicyName + "\n" + ret.Incident.Documentation.Content + "\n" + ret.Incident.URL
	resp, _ := http.Get(telUrl + url.QueryEscape(msg))
	resp.Body.Close()

	c.Data(204, "", nil)
}
