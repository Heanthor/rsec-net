package cmd

import (
	"time"

	"github.com/spf13/viper"

	"github.com/Heanthor/rsec-net/pkg/net"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var announceCmd = &cobra.Command{
	Use:   "start-node",
	Short: "Add the node onto the network.",
	Long:  `Add the node onto the network.`,
	Run: func(c *cobra.Command, args []string) {
		flags := c.Flags()
		n, err := flags.GetString("nodeName")
		if err != nil {
			log.Panic().Err(err).Msg("unable to find node name")
		}
		a, err := flags.GetString("announceAddr")
		if err != nil {
			log.Panic().Err(err).Msg("unable to find announceAddr")
		}
		d, err := flags.GetString("addr")
		if err != nil {
			log.Panic().Err(err).Msg("unable to find addr")
		}
		i, err := flags.GetInt("announceInterval")
		if err != nil {
			log.Panic().Err(err).Msg("unable to find announceInterval")
		}

		initalizeAnnounce(n, d, a, i)
	},
}

func init() {
	announceCmd.Flags().StringP("announceAddr", "a", "224.0.0.1:1145", "Address to announce on")
	announceCmd.Flags().StringP("addr", "d", ":1146", "Address to transmit data on")
	announceCmd.Flags().StringP("nodeName", "n", "", "Node name")
	announceCmd.Flags().IntP("announceInterval", "i", 5, "interval (in seconds) to announce presence to the network")

	viper.BindPFlag("announceAddr", announceCmd.Flags().Lookup("announceAddr"))
	viper.BindPFlag("addr", announceCmd.Flags().Lookup("addr"))
	viper.BindPFlag("nodeName", announceCmd.Flags().Lookup("nodeName"))
	viper.BindPFlag("announceInterval", announceCmd.Flags().Lookup("announceInterval"))

	announceCmd.MarkFlagRequired("nodeName")

	rootCmd.AddCommand(announceCmd)
}

func initalizeAnnounce(nodeName, addr, announceAddr string, interval int) {
	settings := net.InterfaceSettings{
		AnnounceInterval: time.Second * time.Duration(interval),
	}
	i, err := net.NewInterface(nodeName, addr, announceAddr, settings)
	if err != nil {
		log.Panic().Err(err).Msg("unable to start net interface")
	}
	i.StartAnnounce()
}
