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
	"os"

	"github.com/spf13/cobra"
)

// diagramCmd represents the diagram command
var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Create a network diagram from AVD project",
	Long: `This command will create a network diagram from an AVD project.
It uses the same input as the actTopology command and generates a diagram output.
For an example input file run the command with the -e flag`,
	Run: func(cmd *cobra.Command, args []string) {
		example, _ := cmd.Flags().GetBool("example")
		if example {
			fmt.Println("Example:")
			fmt.Println(actExampleConfig)
			os.Exit(0)
		}
		fmt.Println("diagram called")
		folder := cmd.Flag("folder").Value.String()
		output := cmd.Flag("output").Value.String()
		input := cmd.Flag("input").Value.String()
		createDiagram(folder, input, output)
	},
}

func init() {
	rootCmd.AddCommand(diagramCmd)
	diagramCmd.Flags().StringP("folder", "f", "intended/structured_configs", "Folder with the structured configuration files")
	diagramCmd.Flags().StringP("input", "i", "topology.yml", "ACT Topology file")
	diagramCmd.Flags().BoolP("example", "e", false, "Prints an example input file")
	diagramCmd.Flags().StringP("output", "O", "network-diagram.png", "Output diagram file")
}

func createDiagram(folder, inputActTopology, outputDiagram string) {
	// TODO: Implement diagram creation logic
	fmt.Printf("Creating diagram from folder: %s, input: %s, output: %s\n", folder, inputActTopology, outputDiagram)
}
