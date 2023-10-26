/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	generate "gitlab.com/EysteinnSig/stackmap-mapserver/cmd/generate"

	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/backends/mapfile"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/fetch"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/logger"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/rabbitmq"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Mapserver",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := run() //format.Dotest1()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sidecar.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringP("product", "p", "", "Product to base mapfile on.")
	rootCmd.AddCommand(generate.MakeGenerate())
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func run() error {
	passw := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")

	rabbit := rabbitmq.RabbitMQ{}
	err := rabbit.Init(host, passw)
	if err != nil {
		return err
	}

	bend := mapfile.WMS{ApiHost: "stackmap.clouds.is:8080"}
	bend.SQLData.FillFromEnviron()

	/*msg := types.ProductModified{
		Product: "hrit-ash",
		Project: "vedur",
	}
	return bend.OnReceived(&msg)*/
	err = rabbit.Listen(&bend)
	return err
}

func run2() {
	sqldata := types.SQLData{}

	sqldata.SQLUser = getEnv("PSQLUSER", "postgres")
	sqldata.SQLPass = getEnv("PSQLPASS", "")
	sqldata.SQLDB = getEnv("PSQLDB", "postgres")
	sqldata.SQLHost = getEnv("PSQLHOST", "psqlapi-service.default.svc.cluster.local")
	apiHost := getEnv("APIHOST", "")
	mapfilesdir := getEnv("MAPFILESDIR", "/mapfiles/")

	projects := types.UniqueProjects{}
	err := fetch.GetData(apiHost+"/api/v1/projects", &projects)
	if err != nil {
		logger.GetLogger().Fatal(err)
		log.Print(err)
	}
	log.Print("Found ", len(projects.Projects), "projects")

	for _, project := range projects.Projects {
		//err := fetch.FetchAllProducts("./mapfiles/", apiHost, sqldata)
		logger.GetLogger().Info("Fetching products for: ", project)
		outdir := filepath.Join(mapfilesdir, project)
		/*err = os.MkdirAll(outdir, 0655) // 0700)
		if err != nil {
			logger.GetLogger().Fatal(err)
		}*/
		err := fetch.FetchAllProducts(outdir, apiHost, project, sqldata)
		if err != nil {
			logger.GetLogger().Panic(err)
		}
	}

	/*fmt.Println("Going to sleep")
	for {
		time.Sleep(time.Second)
	}*/

	err = rabbitmq.DoRabbitMQ("./mapfiles/", apiHost, sqldata)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
}
