package master

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use: "master",
		Short: "",
		Long:  "master",
		Run: func(cmd *cobra.Command, args []string) {
			Main()
		},
	}
}
