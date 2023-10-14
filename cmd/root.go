/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	cfgFile string

	fileServer http.Handler

	verbose   bool
	port      string
	ipaddress string
	directory string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "serve-api",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
		)

		// Parsing or setting default for directory
		if directory == "" {
			directory = "./"
		} else {
			directory = path.Clean(directory)
		}

		if verbose {
			log.Println("INF verbose=true")
			// log.Println("INF logfile=" + logfile)
			log.Println("INF port=" + port)
			log.Println("INF ipaddress=" + ipaddress)
			log.Println("INF directory=" + directory)
		}

		// Adding handler for api requests to webserver
		http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
			if !r.URL.Query().Has("table") {
				http.Error(w, "Bad request, query for `table` is missing", 400)
				r.Body.Close()
				return
			}
			table := r.URL.Query().Get("table")

			if !r.URL.Query().Has("columns") {
				http.Error(w, "Bad request, query for `columns` is missing", 400)
				return
			}
			columns := r.URL.Query().Get("columns")

			if !r.URL.Query().Has("offset") {
				http.Error(w, "Bad request, query for `offset` is missing", 400)
				return
			}
			offsetString := r.URL.Query().Get("offset")
			offset, err := strconv.Atoi(offsetString)
			if err != nil {
				http.Error(w, "Bad request - offset should be an integer.", 400)
				return
			}

			if !r.URL.Query().Has("limit") {
				http.Error(w, "Bad request, query for `limit` is missing", 400)
				return
			}
			limitString := r.URL.Query().Get("limit")
			limit, err := strconv.Atoi(limitString)
			if err != nil {
				http.Error(w, "Bad request - offset should be an integer.", 400)
				return
			}

			db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
			if err != nil {
				http.Error(w, "Server error - couldn't open database connection", 500)
				log.Fatal(err)
			}

			var results []map[string]interface{}
			db.Table(table).Select(columns).Offset(offset).Limit(limit).Find(&results)

			bResults, err := json.Marshal(results)
			if err != nil {
				http.Error(w, "Bad request - Your data couldn't be retrieved: "+err.Error(), 400)
				return
			}

			_, err = w.Write(bResults)
			if err != nil {
				fmt.Println(err)
			}
		})

		// Adding handler for static files to webserver
		fileServer = http.FileServer(http.Dir(directory))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if verbose {
				log.Println("INF Serving ", r.RequestURI)
			}
			fileServer.ServeHTTP(w, r)
		})

		// Starting Webserver
		fmt.Println("Listening on " + ipaddress + ":" + port + " ...")
		err = http.ListenAndServe(ipaddress+":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.serve.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Every request will be printed.")
	// rootCmd.PersistentFlags().StringVarP(&logfile, "logfile", "l", "", "Output to file instead of stdout.")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "3000", "Bind to specific port. Default is ':3000'.")
	rootCmd.PersistentFlags().StringVarP(&ipaddress, "ip-address", "i", "0.0.0.0", "Bind to specific ip-address. Default is '0.0.0.0'.")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Serve another directory. Default is './'.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".serve" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".serve")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
