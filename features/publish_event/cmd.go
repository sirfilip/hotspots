package publish_event

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"necsam/config"
)

var (
	Cmd         *cobra.Command
	auth        string
	title       string
	description string
	cost        float64
	date        string
	latitude    float64
	longitude   float64
)

func init() {
	Cmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish event",
		Long:  "Publish event using server api",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := time.Parse("2006-01-02", date); err != nil {
				fmt.Println("Invalid date. Must be of format YYYY-MM-DD")
				os.Exit(1)
			}

			apiURL := config.Get("api_url")
			client := resty.New()
			res, err := client.R().
				SetHeader("Content-type", "application/json").
				SetHeader("Authorization", "Bearer "+auth).
				SetBody(map[string]interface{}{
					"title":       title,
					"description": description,
					"cost":        cost,
					"date":        date,
					"latitude":    latitude,
					"longitude":   longitude,
				}).
				Post(apiURL + "/v1/events")
			if err != nil {
				fmt.Println("Failed to perform request")
				fmt.Printf("Got error: %#v", err)
				os.Exit(1)
			}
			if res.StatusCode() == 200 {
				fmt.Println("Login was successful")
				fmt.Println(res.String())
				os.Exit(0)
			}
			if res.StatusCode() == 400 {
				fmt.Println("Bad Request")
				fmt.Println(res.String())
				os.Exit(1)
			}
			if res.StatusCode() == 500 {
				fmt.Println("Server Error")
				os.Exit(1)
			}
		},
	}

	Cmd.Flags().StringVar(&auth, "auth", "", "Auth Token (required)")
	Cmd.MarkFlagRequired("auth")

	Cmd.Flags().StringVar(&title, "title", "", "Title (required)")
	Cmd.MarkFlagRequired("title")

	Cmd.Flags().StringVar(&description, "description", "", "Description (required)")
	Cmd.MarkFlagRequired("password")

	Cmd.Flags().Float64Var(&cost, "cost", 0, "Cost (reqired)")
	Cmd.MarkFlagRequired("cost")

	Cmd.Flags().StringVar(&date, "date", "", "Date (required YYYY-MM-DD)")
	Cmd.MarkFlagRequired("date")

	Cmd.Flags().Float64Var(&latitude, "latitude", 0, "Latitude (reqired)")
	Cmd.MarkFlagRequired("latitude")

	Cmd.Flags().Float64Var(&longitude, "longitude", 0, "Longitude (reqired)")
	Cmd.MarkFlagRequired("longitude")
}
