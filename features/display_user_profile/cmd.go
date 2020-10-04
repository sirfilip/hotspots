package display_user_profile

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"necsam/config"
)

var (
	Cmd  *cobra.Command
	auth string
	id   string
)

func init() {
	Cmd = &cobra.Command{
		Use:   "fetch_user",
		Short: "Fetches user info",
		Long:  "Fetches user info using server api",
		Run: func(cmd *cobra.Command, args []string) {
			apiURL := config.Get("api_url")
			client := resty.New()
			res, err := client.R().
				SetHeader("Content-type", "application/json").
				SetHeader("Authorization", "Bearer "+auth).
				Get(apiURL + "/v1/users/" + id)
			if err != nil {
				fmt.Println("Failed to perform request")
				fmt.Printf("Got error: %#v", err)
				os.Exit(1)
			}
			if res.StatusCode() == 200 {
				fmt.Println(res.String())
				os.Exit(0)
			}
			if res.StatusCode() == 500 {
				fmt.Println("Server Error")
				os.Exit(1)
			}
		},
	}

	Cmd.Flags().StringVar(&auth, "auth", "", "Auth Token (required)")
	Cmd.MarkFlagRequired("auth")

	Cmd.Flags().StringVar(&id, "id", "", "UserID (required)")
	Cmd.MarkFlagRequired("id")
}
