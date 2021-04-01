package main

import (
	"fmt"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"time"

	cmd "ptcg_trader/cmd"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "server"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// go func() {
	// 	log.Info().Msgf("listen http port 6060: %+v", http.ListenAndServe("localhost:6060", nil))
	// }()

	rootCmd.AddCommand(cmd.ServerCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
