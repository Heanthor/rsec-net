package cmd

import (
	"github.com/spf13/cobra"
)

var announceCmd = &cobra.Command{
	Use:   "announce",
	Short: "Announce the node onto the network.",
	Long:  `Announce the node onto the network.`,
	Run: func(c *cobra.Command, args []string) {

	},
}

func init() {
	announceCmd.Flags().IntP("announcePort", "a", 1145, "Port to announce on")
	announceCmd.Flags().IntP("port", "p", 1146, "Port to transmit data on")

	rootCmd.AddCommand(announceCmd)
}
