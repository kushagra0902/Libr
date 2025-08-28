package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/libr-forum/Libr/core/db/internal/models"
	// "github.com/libr-forum/Libr/core/db/internal/node"
	// "go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

// SetupMongo initializes the global MongoClient
// func SetupMongo(uri string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to MongoDB: %w", err)
// 	}
//
// 	if err := client.Ping(ctx, nil); err != nil {
// 		return fmt.Errorf("failed to ping MongoDB: %w", err)
// 	}
//
// 	MongoClient = client
// 	log.Println("✅ MongoDB connected")
// 	return nil
// }

// DisconnectMongo gracefully closes the MongoDB connection
// func DisconnectMongo() {
// 	if MongoClient != nil {
// 		if err := MongoClient.Disconnect(context.Background()); err != nil {
// 			log.Println("⚠️ Error disconnecting MongoDB:", err)
// 		} else {
// 			log.Println("🛑 MongoDB disconnected")
// 		}
// 	}
// }

// 🚀 Uses global MongoClient and ctx
// func GetDbAddr() ([]*models.Node, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	collection := MongoClient.Database("Addrs").Collection("nodes") // replace with actual DB & collection
// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)
//
// 	var nodeList []*models.Node
// 	for cursor.Next(ctx) {
// 		var doc struct {
// 			NodeId string `bson:"node_id"`
// 			PeerId string `bson:"peer_id"`
// 		}
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, err
// 		}
//
// 		nodeId, _ := node.DecodeNodeID(doc.NodeId)
// 		node := &models.Node{
// 			NodeId: nodeId,
// 			PeerId: doc.PeerId,
// 		}
// 		nodeList = append(nodeList, node)
// 	}
// 	return nodeList, nil
// }

// func GetOnlineMods() ([]*models.Mod, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	collection := MongoClient.Database("Addrs").Collection("mods")
// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)
//
// 	var mods []*models.Mod
// 	for cursor.Next(ctx) {
// 		var doc struct {
// 			IP        string `bson:"ip"`
// 			Port      string `bson:"port"`
// 			PublicKey string `bson:"public_key"`
// 		}
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, err
// 		}
// 		mods = append(mods, &models.Mod{
// 			IP:        doc.IP,
// 			Port:      doc.Port,
// 			PublicKey: doc.PublicKey,
// 		})
// 	}
// 	fmt.Println("Online mods:", mods)
// 	return mods, nil
// }

// func GetRelayAddr() ([]string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	collection := MongoClient.Database("Addrs").Collection("relays")
// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch relay addresses: %w", err)
// 	}
// 	defer cursor.Close(ctx)
//
// 	var relayList []string
// 	for cursor.Next(ctx) {
// 		var doc struct {
// 			Address string `bson:"address"`
// 		}
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, fmt.Errorf("failed to decode relay document: %w", err)
// 		}
// 		if strings.HasPrefix(doc.Address, "/") {
// 			relayList = append(relayList, strings.TrimSpace(doc.Address))
// 		}
// 	}
//
// 	return relayList, nil
// }

// -------------------- ACTIVE CODE --------------------

// type relayResp struct {
// 	RelayList []relays `json:"relaylist"`
// }

type modResp struct {
	ModLists []models.Mods `json:"mod_list"`
}

// type relays struct {
// 	Address string `json:"address"`
// }

type NodeResp struct {
	NodesLists []models.Node `json:"boot_list"`
}

// ✅ Fetch relay addresses
func GetRelayAddrFromJSServer() ([]string, error) {
	serverURL := os.Getenv("JS_ServerURL")

	req, err := http.NewRequest("GET", serverURL+"/api/getrelay", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	var payload relayResp
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	fmt.Println(resp.Body)
	fmt.Println(payload)

	var addresses []string
for _, relay := range payload.RelayList.Relays {
    addresses = append(addresses, relay.Address)
}
	fmt.Println(addresses)
	return addresses, nil
}


type relayResp struct {
    RelayList struct {
        Relays []struct {
            Address string `json:"address"`
        } `json:"relaylist"`
    } `json:"relay_list"`
}

// ✅ Fetch mods
func GetModsFromJSServer() ([]*models.Mods, error) {
	fmt.Println("Getting Mods from JS Server...")

	serverURL := "https://libr-q0ok.onrender.com"
	req, err := http.NewRequest("GET", serverURL+"/api/getmod", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	fmt.Println("mod res body = ", string(bodyBytes))

	// Temporary struct for decoding raw data
	type RawMod struct {
		PeerId    string `json:"peer_id"`
		PublicKey string `json:"public_key"`
	}

	var rawResp struct {
		ModLists []RawMod `json:"mod_list"`
	}

	if err := json.Unmarshal(bodyBytes, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	var modReturnList []*models.Mods
	for _, raw := range rawResp.ModLists {
		modReturnList = append(modReturnList, &models.Mods{
			Peerid:    raw.PeerId,
			PublicKey: raw.PublicKey,
		})
	}

	return modReturnList, nil
}

// ✅ Fetch DB nodes
func GetDBFromJSServer() ([]*models.Node, error) {
	fmt.Println("Getting DB from JS Servers...")

	serverURL := "https://libr-q0ok.onrender.com"
	req, err := http.NewRequest("GET", serverURL+"/api/getboot", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status code: %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("res body = ", string(bodyBytes))

	// Temporary struct for decoding
	type RawNode struct {
		NodeId string `json:"node_id"`
		PeerId string `json:"peer_id"`
	}

	var rawResp struct {
		NodesLists []RawNode `json:"boot_list"`
	}

	if err := json.Unmarshal(bodyBytes, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	var result []*models.Node
	for _, raw := range rawResp.NodesLists {
		decoded, err := base64.StdEncoding.DecodeString(raw.NodeId)
		if err != nil || len(decoded) != 20 {
			fmt.Printf("⚠ Skipping node with invalid node_id: %s\n", raw.NodeId)
			continue
		}

		var nodeId [20]byte
		copy(nodeId[:], decoded)

		result = append(result, &models.Node{
			NodeId: nodeId,
			PeerId: raw.PeerId,
		})
	}

	return result, nil
}
