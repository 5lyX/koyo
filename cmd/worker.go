/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"koyo/worker"
	"log"

	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Worker Command to operate a Koyo worker node.",
	Long: `koyo worker command.

The worker runs tasks and responds to the manager's requests about task state.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		name, _ := cmd.Flags().GetString("name")
		dbType, _ := cmd.Flags().GetString("dbType")

		log.Printf("Starting worker.")
		w := worker.New(name, dbType)
		api := worker.Api{Address: host, Port: port, Worker: w}
		go w.RunTasks()
		go w.CollectStats()
		go w.UpdateTasks()
		log.Printf("Starting worker API on http://%s:%d", host, port)
		api.Start()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().StringP("host", "H", "0.0.0.0", "Hostname or IP Address")
	workerCmd.Flags().IntP("port", "p", 5556, "Port on which to listen")
	workerCmd.Flags().StringP("name", "n", fmt.Sprintf("worker-%s", uuid.New().String()), "Name of the worker")
	workerCmd.Flags().StringP("dbType", "d", "memory", "Type of datastore to use for tasks (\"memory\" or \"persistent\")")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
