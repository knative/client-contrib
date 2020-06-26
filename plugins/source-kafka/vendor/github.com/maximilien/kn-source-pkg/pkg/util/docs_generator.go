// Copyright © 2020 The Knative Authors
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

package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	readmePath     = "."
	readmeFilename = "README.md"
	headerFilepath = "./docs/header.md"
	footerFilepath = "./docs/footer.md"
)

// ReadmeGenerator is to generate README.md for a plugin command,
// including header.md and footer.md under doc folder and command usages
func ReadmeGenerator(cmd *cobra.Command) error {
	mdFilepath := filepath.Join(readmePath, readmeFilename)
	return readmeGenerator(cmd, mdFilepath)
}

//Private

func readmeGenerator(cmd *cobra.Command, filepath string) error {
	var (
		header string
		footer string
		err    error
	)
	writer := bytes.NewBufferString("")
	err = docGenerator(cmd, writer)
	if err != nil {
		return err
	}

	mdFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer mdFile.Close()

	if header, err = getHeader(cmd); err != nil {
		return err
	}
	if _, err := io.WriteString(mdFile, header+"\n"); err != nil {
		return err
	}

	if _, err = io.WriteString(mdFile, "## Usage\n\n"); err != nil {
		return err
	}
	reader := bufio.NewReader(writer)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if strings.Contains(line, "##") {
			line = strings.Replace(line, "##", "###", 1)
		}
		_, err = io.WriteString(mdFile, line)
		if err != nil {
			return err
		}
	}

	if footer, err = getFooter(); err != nil {
		return err
	}
	if _, err := io.WriteString(mdFile, footer+"\n"); err != nil {
		return err
	}

	return nil
}

func docGenerator(cmd *cobra.Command, w io.Writer) error {
	if err := doc.GenMarkdownCustom(cmd, w, linkConverter); err != nil {
		return err
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := docGenerator(c, w); err != nil {
			return err
		}
	}
	return nil
}

func linkConverter(ins string) string {
	outs := strings.TrimSuffix(ins, ".md")
	outs = strings.Replace(outs, "_", "-", -1)
	return "#" + outs
}

func getHeader(cmd *cobra.Command) (string, error) {
	if dirExists(headerFilepath) {
		return readFileContent(headerFilepath)
	}
	name := cmd.CommandPath()
	return fmt.Sprintf("# %s\n\nKnative Client plugin `%s`\n", name, name), nil
}

func getFooter() (string, error) {
	if dirExists(footerFilepath) {
		return readFileContent(footerFilepath)
	}
	return `## More information
	
* [Knative Client](https://github.com/knative/client)
* [How to contribute a plugin](https://github.com/knative/client-contrib#how-to-contribute-a-plugin)
`, nil
}

func readFileContent(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func dirExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
