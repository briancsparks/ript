package cmd

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"

  "github.com/spf13/cobra"
)

// slurpCmd represents the slurp command
var slurpCmd = &cobra.Command{
  Use:   "slurp",
  Short: "Grab a dir tree and copy it to become a template",

  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("slurp called")
  },
}

func init() {
  rootCmd.AddCommand(slurpCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // slurpCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // slurpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
