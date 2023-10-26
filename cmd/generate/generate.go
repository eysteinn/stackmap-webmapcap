/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/backends/mapfile"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"
)

var project string = ""
var product string = ""
var rootdir string = ""

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backendstr := args[0]
		fmt.Println("Arg: ", backendstr)
		if backendstr == "mapfile" {
			bend := mapfile.WMS{ApiHost: "stackmap.clouds.is:8080", MapfilesRoot: rootdir}
			bend.SQLData.FillFromEnviron()
			err := bend.OnReceived(&types.ProductModified{
				Product: product,
				Project: project,
			})
			if err != nil {
				log.Fatal(err)
			}

		}
		fmt.Println("generate called")
	},
}

func MakeGenerate() *cobra.Command {
	generateCmd.Flags().StringVar(&project, "project", "", "project name")
	generateCmd.Flags().StringVar(&product, "product", "", "product name")
	generateCmd.Flags().StringVar(&rootdir, "rootdir", "", "root directory")

	return generateCmd
}
func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Generate() {

	bend := mapfile.WMS{ApiHost: "stackmap.clouds.is:8080"}
	bend.SQLData.FillFromEnviron()

	/*msg := types.ProductModified{
		Product: "hrit-ash",
		Project: "vedur",
	}
	return bend.OnReceived(&msg)*/
}
