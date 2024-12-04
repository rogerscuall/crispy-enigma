/*
Copyright © 2024 Roger Gomez rogerscuall@gmail.com

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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rogerscuall/crispy-enigma/internal/act"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var actCVP act.ActCVPConfigs

// actCleanCmd represents the actClean command
var actCleanCmd = &cobra.Command{
	Use:   "actClean",
	Short: "Cleans a production AVD designed configuration to be used with ACT",
	Long: `Arista Cloud Test (ACT) is runs in virtualized devices, that do not support all the physical devices features.
Also we might now want to expose sensitive information like passwords or other secrets to ACT.
This command will clean a production AVD designed configuration to be used with ACT.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("actClean called")
		folder := cmd.Flag("folder").Value.String()
		outputFolder := cmd.Flag("update").Value.String()
		fullOutputFolder := filepath.Join(outputFolder, folder)
		fmt.Println("Folder: ", fullOutputFolder)
		if _, err := os.Stat(fullOutputFolder); os.IsNotExist(err) {
			err := os.MkdirAll(fullOutputFolder, os.ModePerm)
			if err != nil {
				log.Fatalf("Error creating output folder: %v", err)
			}
		}
		if err := viper.UnmarshalKey("act_cvp", &actCVP); err != nil {
			fmt.Println("Error unmarshaling provider: " + err.Error())
		}
		if len(actCVP) != 1 {
			cobra.CheckErr("Only one configuration for CVP is supported")
		}

		log.Print("Folder: ", folder)
		files, err := getConfigFiles(folder)
		if err != nil {
			log.Fatalf("Error reading folder: %v", err)
		}
		// Example configuration string
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				cobra.CheckErr(err)
			}
			content, err := io.ReadAll(f)
			if err != nil {
				cobra.CheckErr(err)
			}
			config := string(content)
			cvpHost := actCVP[0].Host + ":" + actCVP[0].Port
			updatedConfig := act.CleanConfig(config, cvpHost, actCVP[0].VRF)
			fileName := filepath.Base(file)
			outputFile := filepath.Join(fullOutputFolder, fileName)
			err = os.WriteFile(outputFile, []byte(updatedConfig), 0644)
			if err != nil {
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(actCleanCmd)
	actCleanCmd.Flags().StringP("folder", "f", "intended/configs", "Folder where the running config files are located")
	actCleanCmd.Flags().StringP("update", "u", "updated", "Updated configuration for ACT")
}
