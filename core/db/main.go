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

	"github.com/libr-forum/Libr/core/db/config"
	"github.com/libr-forum/Libr/core/db/internal/keycache"
	peer "github.com/libr-forum/Libr/core/db/internal/network/peers"
	"github.com/libr-forum/Libr/core/db/internal/routing"
	"github.com/libr-forum/Libr/core/db/internal/utils"
	"github.com/joho/godotenv"
)

var JS_API_key string
var JS_ServerURL string

func main() {
	keycache.InitKeys()
	godotenv.Load()
	JS_API_key = config.Cfg.JSAPIKey
	JS_ServerURL = config.Cfg.JSServerURL
	if JS_API_key == "" || JS_ServerURL == "" {
		fmt.Println("[DEBUG] Missing JS API key or server URL")
		return
	}
	//utils.SetupMongo("mongodb+srv://peer:peerhehe@cluster0.vswojqe.mongodb.net/")
	//utils.SetupMongo(JS_ServerURL)
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