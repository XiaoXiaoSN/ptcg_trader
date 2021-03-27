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
	rootCmd.AddCommand(cmd.ServerCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
