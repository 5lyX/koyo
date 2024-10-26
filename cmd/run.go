/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a new task.",
	Long: `koyo run command.

The run command starts a new task.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, _ := cmd.Flags().GetString("manager")
		filename, _ := cmd.Flags().GetString("filename")

		fullFilePath, err := filepath.Abs(filename)
		if err != nil {
			log.Fatal(err)
		}
		if !fileExists(fullFilePath) {
			log.Fatalf("File %s does not exist.", filename)
		}

		log.Printf("Using manager: %v\n", manager)
		log.Printf("Using file: %v\n", fullFilePath)
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Unable to read file: %v", filename)
		}
		log.Printf("Data: %v\n", string(data))

		url := fmt.Sprintf("http://%s/tasks", manager)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Panic(err)
		}

		if resp.StatusCode != http.StatusCreated {
			log.Printf("Error sending request: %v", resp.StatusCode)
		}

		defer resp.Body.Close()
		log.Println("Successfully sent task request to manager")
	},
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, fs.ErrNotExist)
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("manager", "m", "localhost:5555", "Manager to talk to")
	runCmd.Flags().StringP("filename", "f", "task.json", "Task specification file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
