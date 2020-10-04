package login

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"necsam/config"
)

var (
	Cmd      *cobra.Command
	email    string
	password string
)

func init() {
	Cmd = &cobra.Command{
		Use:   "login",
		Short: "Login User",
		Long:  "Login user using server api",
		Run: func(cmd *cobra.Command, args []string) {
			apiURL := config.Get("api_url")
			client := resty.New()
			res, err := client.R().
				SetHeader("Content-type", "application/json").
				SetBody(map[string]interface{}{
					"email":    email,
					"password": password,
				}).
				Post(apiURL + "/v1/auth/login")
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
				os.Exit(2)
			}
			if res.StatusCode() == 500 {
				fmt.Println("Server Error")
				os.Exit(2)
			}
		},
	}

	Cmd.Flags().StringVarP(&email, "email", "e", "", "Email (required)")
	Cmd.MarkFlagRequired("email")

	Cmd.Flags().StringVarP(&password, "password", "p", "", "Password (required)")
	Cmd.MarkFlagRequired("password")
}
