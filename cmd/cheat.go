package cmd

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"github.com/briancsparks/ript/ript"
	"log"

	"github.com/spf13/cobra"
)

var prj, destDir string

// cheatCmd represents the cheat command
var cheatCmd = &cobra.Command{
	Use:   "cheat",
	Short: "Get away with cheating by generating from template",

	Run: func(cmd *cobra.Command, args []string) {
		//err := ript.Cheat("gocli", "./one")

    if len(args) == 0 {
      log.Fatalf("Need template name");
    }

    // TODO: Call Cheat2
		err := ript.Cheat(args[0], destDir)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cheatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cheatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cheatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cheatCmd.Flags().StringVar(&prj, "projectname", "projektzero", "help")
	cheatCmd.Flags().StringVar(&destDir, "dest", "nodestgiven", "The Dir")

	bindFlags(cheatCmd)
}
