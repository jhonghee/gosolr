// Copyright © 2018 Jhonghee Park <jhonghee@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	archiveURL      = "http://archive.apache.org/dist/lucene/solr/"
	dynURL          = "https://www.apache.org/dyn/closer.lua/lucene/solr/"
	mirrorSelector  = "body > pre:nth-child(2) > a"
	archiveSelector = "body > pre:nth-child(4) > a"
)

var cfgFile string
var mirror bool

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gosolr",
	Short: "gosolr helps you get started with Solr from command line",
	Long: `From installation to operation, gosolr will help you get started
	with Solr from your command line.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gosolr.yaml)")

	RootCmd.PersistentFlags().BoolVarP(&mirror, "mirror", "m", true, "Toggle the use of a mirror for the commands")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("archiveURL", archiveURL)
	viper.SetDefault("dynURL", dynURL)

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".gosolr") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getMirrorURL() (string, error) {
	dynURL := viper.GetString("dynURL")
	doc, err := goquery.NewDocument(dynURL)
	if err != nil {
		log.Fatal(err)
	}
	return doc.Find("body > div:nth-child(3) > p:nth-child(3) > a > strong").Html()
}

func getVersions(url string, selector string) []string {
	var versions = make([]string, 0)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(selector).Each(func(index int, item *goquery.Selection) {
		html, _ := item.Html()
		semver, _ := regexp.Compile(`^(\d+)\.(\d+)\.?(\d+)?(.*)/$`)
		submatches := semver.FindSubmatch([]byte(html))
		if len(submatches) > 0 {
			major := string(submatches[1])
			minor := getVersionDigit(submatches[2])
			patch := getVersionDigit(submatches[3])
			versions = append(versions, fmt.Sprintf("%s%s%s%s", major, minor, patch, string(submatches[4])))
		}
	})
	return versions
}
