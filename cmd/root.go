/*
Copyright © 2022 Mateusz Urbanek <mateusz.urbanek.98@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/shanduur/cli/cmd/version"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "FIXME: cli-template is a example app using Cobra",
	Long: `FIXME: this app is product of building the app from a template repo
	
The repo is meant to simplify setting up new repos for CLI (and not only) apps.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initRootCmd()
	initSubCmd()
	cobra.OnInitialize(initConfig, initLogger)
}

func initLogger() {
	timeFormat := viper.GetString("log_timestamp_format")

	switch strings.ToUpper(timeFormat) {
	case zerolog.TimeFormatUnix, "UNIX":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	case zerolog.TimeFormatUnixMicro:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro

	case zerolog.TimeFormatUnixNano:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	case zerolog.TimeFormatUnixMs:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	case "RFC3339":
		zerolog.TimeFieldFormat = time.RFC3339

	case "RFC3339Nano":
		zerolog.TimeFieldFormat = time.RFC3339Nano

	case "RFC1123":
		zerolog.TimeFieldFormat = time.RFC1123

	case "RFC1123Z":
		zerolog.TimeFieldFormat = time.RFC1123Z

	case "RFC822":
		zerolog.TimeFieldFormat = time.RFC822

	case "RFC822Z":
		zerolog.TimeFieldFormat = time.RFC822Z

	case "RFC850":
		zerolog.TimeFieldFormat = time.RFC850

	default:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	lvl := viper.GetString("log_level")

	level, err := zerolog.ParseLevel(lvl)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Debug().Str("log_level", zerolog.GlobalLevel().String()).Send()
}

func initRootCmd() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default \"$HOME/.cli.toml\")")

	rootCmd.PersistentFlags().String("log-level", "info", "configure default log level for application")
	_ = viper.BindPFlag("log_level", rootCmd.Flag("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("log_level", "info")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".cli")
	}

	viper.SetEnvPrefix("cli")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Trace().Str("file", viper.ConfigFileUsed()).Msg("config file loaded")
	}
}

func initSubCmd() {
	rootCmd.AddCommand(version.NewVersionCmd())
}
