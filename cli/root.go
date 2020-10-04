package cli

import (
	"github.com/spf13/cobra"

	"necsam/features/display_user_profile"
	"necsam/features/login"
	"necsam/features/publish_event"
	"necsam/features/register"
	"necsam/features/user_activation"
)

var RootCmd = &cobra.Command{Use: "necsam"}

func init() {
	RootCmd.AddCommand(register.Cmd)
	RootCmd.AddCommand(user_activation.Cmd)
	RootCmd.AddCommand(login.Cmd)
	RootCmd.AddCommand(publish_event.Cmd)
	RootCmd.AddCommand(display_user_profile.Cmd)
}
