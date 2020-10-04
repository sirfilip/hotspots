package register

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"necsam/config"
)

var (
	Cmd *cobra.Command

	username string
	email    string
	password string
)

func init() {
	apiURL := config.Get("api_url")

	Cmd = &cobra.Command{
		Use:   "register",
		Short: "Register new user",
		Long:  `Connects to necsam api and performs user registration`,
		Run: func(cmd *cobra.Command, args []string) {
			client := resty.New()
			res, err := client.R().
				SetHeader("Content-type", "application/json").
				SetBody(map[string]interface{}{
					"username": username,
					"email":    email,
					"password": password,
				}).
				Post(apiURL + "/v1/auth/register")
			if err != nil {
				fmt.Println("Failed to perform request")
				fmt.Printf("Got error: %#v", err)
				os.Exit(1)
			}
			if res.StatusCode() == 201 {
				fmt.Println("Registration was successful")
				os.Exit(0)
			}
			if res.StatusCode() == 400 {
				fmt.Println(res.String())
				os.Exit(2)
			}
		},
	}

	Cmd.Flags().StringVarP(&username, "username", "u", "", "Username (required)")
	Cmd.MarkFlagRequired("username")

	Cmd.Flags().StringVarP(&email, "email", "e", "", "Email (required)")
	Cmd.MarkFlagRequired("email")

	Cmd.Flags().StringVarP(&password, "password", "p", "", "Password (required)")
	Cmd.MarkFlagRequired("password")
}
