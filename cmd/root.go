package cmd

import (
	"github.com/llamadeus/keyval-server/internal"
	"github.com/spf13/cobra"
	"log"
)

var (
	address         string
	storageFilePath string

	rootCmd = &cobra.Command{
		Use:  "keyval-server",
		Long: "adasd",
		Run: func(cmd *cobra.Command, args []string) {
			internal.StartKeyValServer(address, storageFilePath)
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&address, "address", "a", ":3000", "The server's bind address")
	rootCmd.Flags().StringVarP(&storageFilePath, "storage", "s", "", "The storage file path")

	_ = rootCmd.MarkFlagRequired("storage")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
