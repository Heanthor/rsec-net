package cmd

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	// pprof network server
	_ "net/http/pprof"

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

		initalizeAnnounce(n)
	},
}

func init() {
	announceCmd.Flags().String("announceAddr", "239.0.0.0:1145", "Address to announce on (host:port)")
	announceCmd.Flags().String("announceListenPort", "1145", "Port to listen for announce packets on")
	announceCmd.Flags().BoolP("announceMulticast", "m", false, "true if announcing using multicast")
	announceCmd.Flags().String("dataListenPort", "1146", "Port to listen for data packets on")
	announceCmd.Flags().StringP("nodeName", "n", "", "Node name")
	announceCmd.Flags().IntP("announceInterval", "i", 5, "interval (in seconds) to announce presence to the network")

	viper.BindPFlag("announceAddr", announceCmd.Flags().Lookup("announceAddr"))
	viper.BindPFlag("announceListenPort", announceCmd.Flags().Lookup("announceListenPort"))
	viper.BindPFlag("announceMulticast", announceCmd.Flags().Lookup("announceMulticast"))
	viper.BindPFlag("dataAddr", announceCmd.Flags().Lookup("dataAddr"))
	viper.BindPFlag("dataListenPort", announceCmd.Flags().Lookup("dataListenPort"))
	viper.BindPFlag("nodeName", announceCmd.Flags().Lookup("nodeName"))
	viper.BindPFlag("announceInterval", announceCmd.Flags().Lookup("announceInterval"))

	rootCmd.AddCommand(announceCmd)
}

func initalizeAnnounce(nodeName string) {
	var profiler interface{ Stop() }

	if viper.GetBool("profile") {
		pp := viper.GetString("profilePath")
		mode := profile.CPUProfile
		doProfile := true

		switch viper.GetString("profileMode") {
		case "cpu":
			log.Info().Msg("Profiling in cpu mode")
		case "mem":
			log.Info().Msg("Profiling in memory mode")
			mode = profile.MemProfile
		case "mutex":
			log.Info().Msg("Profiling in mutex mode")
			mode = profile.MutexProfile
		case "block":
			runtime.SetBlockProfileRate(1)
			log.Info().Msg("Profiling in block mode")
			mode = profile.BlockProfile
		case "goroutine":
			log.Info().Msg("Profiling in goroutine mode")
			mode = profile.GoroutineProfile
		default:
			doProfile = false
			// do nothing
		}

		if doProfile {
			profiler = profile.Start(mode, profile.ProfilePath(pp), profile.NoShutdownHook)
		}
	}

	if viper.GetBool("netProfile") {
		go func() {
			runtime.SetBlockProfileRate(1)
			port := viper.GetInt("netProfilePort")
			log.Info().Msgf("Started pprof http server on %d", port)
			http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
		}()
	}

	interval := viper.GetInt("announceInterval")
	settings := net.InterfaceSettings{
		AnnounceInterval: time.Second * time.Duration(interval),
	}

	// create data connections
	dataReceive := viper.GetString("dataListenPort")
	listenAddr := ":" + dataReceive
	dr, err := udp.NewUniReader(listenAddr)
	if err != nil {
		log.Panic().Err(err).Str("listenAddr", listenAddr).Msg("unable to create udp data UniReader")
	}

	// create announce connection
	var ar udp.NetReader
	announceSend := viper.GetString("announceAddr")
	announceReceive := viper.GetString("announceListenPort")
	aListenAddr := ":" + announceReceive
	if viper.GetBool("announceMulticast") {
		ar, err = udp.NewMulticastReader(aListenAddr)
	} else {
		ar, err = udp.NewUniReader(aListenAddr)
	}
	if err != nil {
		log.Panic().Err(err).Str("aListenAddr", aListenAddr).Msg("unable to create udp announce NetReader")
	}

	as, err := udp.NewUDPWriter(announceSend)
	if err != nil {
		log.Panic().Err(err).Str("dataAddr", announceSend).Msg("unable to create announce udp data UDPWriter")
	}

	i, err := net.NewInterface(nodeName, dr, as, ar, settings)
	if err != nil {
		log.Panic().Err(err).Msg("unable to start net interface")
	}
	i.StartAnnounce()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-c
	log.Info().Msg("CTRL-C pressed, stopping...")
	profiler.Stop()
	i.Close()
	os.Exit(0)
}
