/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// parseOutputCmd represents the parseOutput command
var parseOutputCmd = &cobra.Command{
	Use:   "parseOutput",
	Short: "Find the number of occurrences of a keyword in a file",
	Long: `Find the number of occurrences of a keyword in a file.
If the number of occurrences is greater than a given number, it will return an error.
It will print to the console the number of occurrences and the file name and the line.
The match will be case insensitive.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("parseOutput called")
		var flag bool
		count, err := cmd.Flags().GetUint("count")
		if err != nil {
			cobra.CheckErr(err)
		}
		log.Print("Count number is: ", count)
		inputKeywords := cmd.Flag("keyword").Value.String()
		keywords := strings.Split(inputKeywords, ",")
		if len(keywords) == 0 {
			cobra.CheckErr("Keyword is required")
		}
		log.Println("Keywords: ", keywords)
		filePath := cmd.Flag("file").Value.String()
		file, err := os.Open(filePath)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer file.Close()

		// Maps for counting occurrences and storing lines
		keywordCounts := make(map[string]int)
		keywordLines := make(map[string][]string)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Get the original line and its lowercase version for comparison
			line := scanner.Text()
			lowerLine := strings.ToLower(line)

			// Check if the line contains any of the keywords
			for _, keyword := range keywords {
				lowerKeyword := strings.ToLower(keyword)
				if strings.Contains(lowerLine, lowerKeyword) {
					// Increment count and store the line
					keywordCounts[lowerKeyword]++
					keywordLines[lowerKeyword] = append(keywordLines[lowerKeyword], line)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		// Print the lines for keywords that meet the threshold
		for keyword, times := range keywordCounts {
			if times > int(count) {
				flag = true
				fmt.Printf("Keyword '%s' appears %d times, exceeding the threshold of %d. Matching lines:\n", keyword, times, count)
				for _, line := range keywordLines[keyword] {
					fmt.Println(line)
				}
				fmt.Println() // Add a blank line for readability
			}
		}
		if flag {
			cobra.CheckErr("One or more keywords exceeded the threshold")
		}
	},
}

func init() {
	rootCmd.AddCommand(parseOutputCmd)

	parseOutputCmd.Flags().UintP("count", "c", 0, "Number of occurrences to count")
	parseOutputCmd.Flags().StringP("keyword", "k", "", "Keyword to search")
	parseOutputCmd.Flags().StringP("file", "f", "", "File to search")
}
