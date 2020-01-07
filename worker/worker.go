package worker

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:                        "worker",
		Short:                      "w",
		Long:                       "worker",
		Run: func(cmd *cobra.Command, args []string) {
			Main()
		},
	}
}
