package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Heanthor/rsec-net/internal/udp"

	"github.com/rs/zerolog"

	"github.com/pkg/profile"
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
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if viper.GetBool("verbose") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		n := viper.GetString("nodeName")
		if n == "" {
			log.Panic().Msg("node name is required")
		}
		a := viper.GetString("announceAddr")
		d := viper.GetString("dataAddr")
		i := viper.GetInt("announceInterval")

		initalizeAnnounce(n, d, a, i)
	},
}

func init() {
	announceCmd.Flags().StringP("announceAddr", "a", "239.0.0.0:1145", "Address to announce on")
	announceCmd.Flags().BoolP("announceMulticast", "m", false, "true if announcing using multicast")
	announceCmd.Flags().StringP("dataAddr", "d", ":1146", "Address to transmit data on")
	announceCmd.Flags().StringP("nodeName", "n", "", "Node name")
	announceCmd.Flags().IntP("announceInterval", "i", 5, "interval (in seconds) to announce presence to the network")

	viper.BindPFlag("announceAddr", announceCmd.Flags().Lookup("announceAddr"))
	viper.BindPFlag("announceMulticast", announceCmd.Flags().Lookup("announceMulticast"))
	viper.BindPFlag("dataAddr", announceCmd.Flags().Lookup("dataAddr"))
	viper.BindPFlag("nodeName", announceCmd.Flags().Lookup("nodeName"))
	viper.BindPFlag("announceInterval", announceCmd.Flags().Lookup("announceInterval"))

	rootCmd.AddCommand(announceCmd)
}

func initalizeAnnounce(nodeName, addr, announceAddr string, interval int) {
	if viper.GetBool("profile") {
		// start cpu profiling
		defer profile.Start().Stop()
	}

	settings := net.InterfaceSettings{
		AnnounceInterval: time.Second * time.Duration(interval),
	}

	// create data connection
	u, err := udp.NewUniNet(addr)
	if err != nil {
		log.Panic().Err(err).Msg("unable to create udp data connection")
	}

	// create announce connection
	var announceConn udp.NetCommunicator

	if viper.GetBool("announceMulticast") {
		announceConn, err = udp.NewMulticastNet(addr)
	} else {
		announceConn, err = udp.NewUniNet(addr)
	}
	if err != nil {
		log.Panic().Err(err).Msg("unable to create udp announce connection")
	}

	i, err := net.NewInterface(nodeName, u, announceConn, settings)
	if err != nil {
		log.Panic().Err(err).Msg("unable to start net interface")
	}
	i.StartAnnounce()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info().Msg("CTRL-C pressed, stopping...")
		i.Close()
		os.Exit(0)
	}()

	time.Sleep(time.Hour * time.Duration(1))
}
