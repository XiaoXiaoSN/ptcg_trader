package main

import (
	"fmt"
	"math/rand"
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
	// import _ "net/http/pprof"
	// go func() {
	// 	log.Printf("listen http port 6060: %+v", http.ListenAndServe("localhost:6060", nil))
	// }()

	rootCmd.AddCommand(
		cmd.ServerCmd,
		cmd.MatcherCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
