/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"wyag/internal"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path] Where to create the repository.",
	Short: "Initialize a new, empty repository.",
	Long: `This command creates an empty Git repository - basically a .git directory with subdirectories for objects, refs/heads, refs/tags, and template files. An initial HEAD file that references the HEAD of the master
	branch is also created.

	Running git init in an existing repository is safe. It will not overwrite things that are already there. The primary reason for rerunning git init is to pick up newly added templates.`,
	Args:                  cobra.MaximumNArgs(1),
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		var path string
		if len(args) == 1 {
			path = args[0]
		} else {
			path = "./"
		}

		err := internal.RepoCreate(path)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
