package cmd

import (
	"fmt"
	"go-tree/constant"
	"go-tree/internal"
	"os"

	"github.com/spf13/cobra"
)

var flags map[string]interface{}

var goTree = &cobra.Command{
	Use:   "./main",
	Short: "unix command \"tree\" implementation in go",
	Long:  "go-tree is a cli tool which draws a tree of directory structure",
	Run: func(cmd *cobra.Command, args []string) {
		internal.DrawTree(flags)
	},
}

func init() {
	flags = map[string]interface{}{}
	flags[constant.Root] = goTree.PersistentFlags().StringP(constant.Root, "r", ".", "Root path of the tree")
	flags[constant.Path] = goTree.PersistentFlags().BoolP(constant.Path, "f", false, "Flag to show fullpaths")
	flags[constant.Dir] = goTree.PersistentFlags().BoolP(constant.Dir, "d", false, "Flag to only list directories")
	flags[constant.Level] = goTree.PersistentFlags().IntP(constant.Level, "L", 0, "Max level of tree depth")
	flags[constant.Mode] = goTree.PersistentFlags().BoolP(constant.Mode, "p", false, "Flag to show permission modes")
	flags[constant.Time] = goTree.PersistentFlags().BoolP(constant.Time, "t", false, "Flag to sort output by modified time")
	flags[constant.JSON] = goTree.PersistentFlags().BoolP(constant.JSON, "J", false, "Prints tree in JSON format")
	flags[constant.XML] = goTree.PersistentFlags().BoolP(constant.XML, "X", false, "Prints tree in XML format")
	flags[constant.Indent] = goTree.PersistentFlags().BoolP(constant.Indent, "i", false, "Prints tree without indentation lines")
}

func Execute() {
	if err := goTree.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
