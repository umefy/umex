/*
Copyright © 2025 LeoZhao0709<leo.zhao.real@gmail.com>

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
package cmd

import (
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/umefy/umex/flags"
	"github.com/umefy/umex/tpl"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new umefy project",
	Run: func(cmd *cobra.Command, args []string) {
		template, _ := cmd.Flags().GetString(flags.TEMPLATE)
		module, _ := cmd.Flags().GetString(flags.MODULE)
		debugMode, _ := cmd.Flags().GetBool(flags.DEBUG)
		localMode, _ := cmd.Flags().GetBool(flags.LOCAL)

		flagsModel := flags.Model{
			Template:  template,
			Module:    module,
			DebugMode: debugMode,
			LocalMode: localMode,
		}

		currDir, err := os.Getwd()
		if err != nil {
			color.Redf("failed to get current directory: %s", err)
			return
		}

		projectDir := currDir
		if len(args) > 0 {
			projectDir = filepath.Join(currDir, args[0])
		}

		switch template {
		case "goWebApp":
			err := tpl.CreateGoWebApp(flagsModel, projectDir)
			if err != nil {
				color.Redf("failed to create project: %s", err)
			}
		case "":
			color.Redf("template is required")
		default:
			color.Redf("template %s not found", template)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().StringP(flags.TEMPLATE, flags.TEMPLATE_SHORT, "", "template to use for the project, support: goWebApp")
	newCmd.Flags().StringP(flags.MODULE, flags.MODULE_SHORT, "", "go mod module name, required for go project")
}
