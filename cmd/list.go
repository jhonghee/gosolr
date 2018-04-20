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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getVersionDigit(input []byte) string {
	if len(input) == 0 {
		return ".0"
	} else {
		return "." + string(input)
	}
}

func showVersions(url string, versions []string) {
	fmt.Printf("From mirror, %s\n", url)
	for _, v := range versions {
		fmt.Println(v)
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all versions of Solr available",
	Long: `List the versions of Solr available from a mirror:

By default, versions available from a chosen mirror will be displayed.
With -m=false, all versions from Solr archive will be displayed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mirror {
			url, err := getMirrorURL()
			if err != nil {
				log.Fatal(err)
			}
			versions := getVersions(url, mirrorSelector)
			showVersions(url, versions)
			return
		}

		archiveURL := viper.GetString("archiveURL")
		versions := getVersions(archiveURL, archiveSelector)
		showVersions(archiveURL, versions)
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
