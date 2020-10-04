package user_activation

import (
	"fmt"
	"necsam/config"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var (
	Cmd  *cobra.Command
	code string
)

func init() {
	Cmd = &cobra.Command{
		Use:   "activate",
		Short: "Activates registered user",
		Long:  `Connects to necsam api and performs user activation`,
		Run: func(cmd *cobra.Command, args []string) {
			apiURL := config.Get("api_url")
			client := resty.New()
			res, err := client.R().
				SetHeader("Content-type", "application/json").
				Get(apiURL + "/v1/auth/activate/" + code)
			if err != nil {
				fmt.Println("Failed to perform request")
				fmt.Printf("Got error: %#v", err)
				os.Exit(1)
			}
			if res.StatusCode() == 204 {
				fmt.Println("Activation was successful")
				os.Exit(0)
			}
			if res.StatusCode() == 400 {
				fmt.Println("Bad Request")
				fmt.Println(res.String())
				os.Exit(2)
			}
			if res.StatusCode() == 500 {
				fmt.Println("Server Error")
				os.Exit(2)
			}
			if res.StatusCode() == 404 {
				fmt.Println("Not Found")
				os.Exit(2)
			}
		},
	}

	Cmd.Flags().StringVarP(&code, "code", "c", "", "Activation Code (required)")
	Cmd.MarkFlagRequired("code")
}
