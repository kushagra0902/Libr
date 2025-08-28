package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libr-forum/Libr/core/db/internal/keycache"
	peer "github.com/libr-forum/Libr/core/db/internal/network/peers"
	"github.com/libr-forum/Libr/core/db/internal/routing"
	"github.com/libr-forum/Libr/core/db/internal/utils"
)

var JS_API_key string
var JS_ServerURL string

func main() {
	keycache.InitKeys()
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }
	JS_API_key = os.Getenv("JS_API_KEY")
	JS_ServerURL = os.Getenv("JS_ServerURL")
	relayAddrs, err := utils.GetRelayAddrFromJSServer()
	
	if err != nil {
		fmt.Println("Error while getting relay address, ", err)
	}

	fmt.Println(relayAddrs)
	peer.StartNode(relayAddrs)	
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("Interrupt received. Exiting gracefully.")
	// if(config.DBtype=="boot"){
	// deleteFromJSServer()
	// }
	deleteFromJSServer()
	fmt.Println("Sent delete request to JS server")
	if routing.GlobalRT != nil {
		routing.GlobalRT.SaveToDBAsync()
		time.Sleep(1 * time.Second)
	}
}


func deleteFromJSServer() error {
	deleteData := map[string]string{
		"peer_id" : peer.PeerID,
	}

	jsonData, err := json.Marshal(deleteData)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("DELETE", JS_ServerURL+"/api/deleteboot", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", JS_API_key)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	fmt.Printf("Successfully deleted relay")
	return nil
}

// func CreateDbConfig() error {
// 	   path := config.GetDBConfigPath()

// 	   if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
// 		   return fmt.Errorf("failed to create modconfig directory: %w", err)
// 	   }

// 	   // Only create the file if it does not already exist
// 	   if _, err := os.Stat(path); os.IsNotExist(err) {
// 		   f, err := os.Create(path)
// 		   if err != nil {
// 			   return fmt.Errorf("failed to create config file: %w", err)
// 		   }
// 		   defer f.Close()

// 		   // Write an empty JSON object if file is new
// 		   _, err = f.WriteString("{}")
// 		   if err != nil {
// 			   return fmt.Errorf("failed to write empty JSON to config file: %w", err)
// 		   }
// 	   }
// 	   return nil
// }
