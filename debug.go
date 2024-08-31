package main

import (
	"flag"
	"fmt"
	"os"
)

func checkForDebug() {

	// Define the --debug flag
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Delete the database.json file if the debug flag is set
	if *dbg {
		fmt.Println("Debug mode enabled: Deleting database.json")
		// Delete database.json
		err := os.Remove("database.json")
		// Check if database.json was deleted successfully
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error deleting database.json: %s\n", err)
		} else if os.IsNotExist(err) {
			fmt.Println("database.json does not exist, nothing to delete.")
		} else {
			fmt.Println("database.json successfully deleted.")
		}
	}
}
