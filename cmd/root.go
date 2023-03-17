/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hb",
	Short: "Application for interacting with Hyundai and its Bluelink service",
	Long: `
hb allows you to authenticate with your Hyundai account and perform actions like
requesting vehicle information or locking or unlocking your vehicle.

hb will generate a .env file for you the first time you run the bare 'hb' command.
Populate the file with the specified credentials and then you can authenticate
and more.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.OpenFile(".env", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("failed to create .env with error: %v", err)
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			log.Fatalf("failed to read length of .env with error: %v", err)
		}

		if !(fi.Size() > 0) {
			f.WriteString("USERNAME=\n")
			f.WriteString("PASSWORD=\n")
			f.WriteString("PIN=")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hyundai-bluelink-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


