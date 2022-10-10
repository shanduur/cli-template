/*
Copyright © 2022 Mateusz Urbanek <mateusz.urbanek.98@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package version

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	Version  string
	Revision string

	versionShort = "Show the version information"
	versionLong  = `Show the version information`

	runVersion = func(cmd *cobra.Command, args []string) int {
		log.Info().Interface("client_version", buildVersion(Version, Revision)).Send()

		return 0
	}
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: versionShort,
		Long:  versionLong,
		Run:   func(cmd *cobra.Command, args []string) { os.Exit(runVersion(cmd, args)) },
	}

	return cmd
}

func buildVersion(version, revision string) map[string]string {
	return map[string]string{
		"version":  version,
		"revision": revision,
	}
}
