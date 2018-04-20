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
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

const lastest = "latest"

var version string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a version of Solr under HOME directory by default",
	Long: `A mirror will be automatically chosen and the zip file of
give version of Solr will be downloaded and extracted under HOME directory:

The available versions can be obtained from command "list".`,
	Run: func(cmd *cobra.Command, args []string) {
		if mirror {
			versions := []string{}
			var mirrorURL string
			for len(versions) == 0 {
				u, err := getMirrorURL()
				if err != nil {
					log.Fatal(err)
				}
				versions = getVersions(u, mirrorSelector)
				if len(versions) > 0 {
					mirrorURL = u
				}
			}
			if version == lastest {
				version = versions[len(versions)-1]
			}
			filename := fmt.Sprintf("solr-%s.zip", version)
			u, _ := url.Parse(mirrorURL)
			u.Path = path.Join(u.Path, version, filename)
			surl := u.String()

			download(filename, surl)

			unzip(filename, "solr-installation")
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&version, "version", "v", lastest, "The version of Solr to install")
}

func download(filename, surl string) {
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	resp, err := http.Get(surl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Downloading %s\n", surl)
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes has been downloaded\n", n)
}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}
