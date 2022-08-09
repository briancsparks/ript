package cmd

/* Copyright © 2022 RIPTENV_IMFO_USER_IRL_NAME <RIPTENV_GITHUB_USER_EMAIL> -- MIT (see LICENSE file) */

import (
  "fmt"
  "github.com/fsnotify/fsnotify"
  "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/pflag"
  "github.com/spf13/viper"
  "os"
  "strings"
)

var cfgFile string
var programRootName = "riptprojectname"

// ---------------------------------------------------------------------------------------------------------------------
// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".programRootName" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigType("yaml")
    viper.SetConfigName("." + programRootName)
  }

  viper.AutomaticEnv()                // read in environment variables that match
  viper.SetEnvPrefix(programRootName) // Automatically upper-cased

  bindFlags(rootCmd)

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }

  viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
  })
  viper.WatchConfig()
}

// =====================================================================================================================
// roothelp.go
// Put bindFlags into roothelp.go as-is

// ---------------------------------------------------------------------------------------------------------------------
// Stolen from: https://github.com/carolynvs/stingoftheviper/blob/main/main.go
// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command /*, v *viper.Viper*/) {
  cmd.Flags().VisitAll(func(f *pflag.Flag) {

    // Environment variables can't have dashes in them, so bind them to their equivalent
    // keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
    if strings.Contains(f.Name, "-") {
      envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
      viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", strings.ToUpper(programRootName), envVarSuffix))
    }

    // Apply the viper config value to the flag when the flag is not set and viper has a value
    if !f.Changed && viper.IsSet(f.Name) {
      val := viper.Get(f.Name)
      cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
    }
  })

  cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {

    // Environment variables can't have dashes in them, so bind them to their equivalent
    // keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
    if strings.Contains(f.Name, "-") {
      envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
      viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", strings.ToUpper(programRootName), envVarSuffix))
    }

    // Apply the viper config value to the flag when the flag is not set and viper has a value
    if !f.Changed && viper.IsSet(f.Name) {
      val := viper.Get(f.Name)
      cmd.PersistentFlags().Set(f.Name, fmt.Sprintf("%v", val))
    }
  })
}
