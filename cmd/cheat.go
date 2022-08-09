package cmd

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"github.com/briancsparks/ript/ript"

	"github.com/spf13/cobra"
)

var prj string

// cheatCmd represents the cheat command
var cheatCmd = &cobra.Command{
	Use:   "cheat",
	Short: "Get away with cheating by generating from template",

	Run: func(cmd *cobra.Command, args []string) {
		//ript.Cheat2(
		//  "/home/sparksb/go/src/bcs/tryouts/__go-project-template/one",
		//  "/home/sparksb/go/src/bcs/tryouts/ript/scratch/one",
		//)
		ript.Cheat("gocli", "./one")
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

	bindFlags(cheatCmd)
}
