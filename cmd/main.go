package main

import (
	"log"
	"os"
	"strconv"

	"github.com/SEG-UNIBE/artio-miner/pkg/miner"
	"github.com/SEG-UNIBE/artio-miner/pkg/storage"
	"github.com/joho/godotenv"
)

/*
main example function for just fetching the data from the relays
*/
func main() {
	log.Println("artio-miner is starting up for you :)")
	startingRelays := []string{"wss://relay.artiostr.ch/", "wss://relay.artio.inf.unibe.ch/"}
	miners := make([]*miner.RelayMiner, 0)
	for _, relay := range startingRelays {
		miners = append(miners, miner.NewMiner(relay))
	}
	err := godotenv.Load(".env")

	uri := os.Getenv("NEO4J_URI")
	password := os.Getenv("NEO4J_PASSWORD")
	db := os.Getenv("NEO4J_DB")
	username := os.Getenv("NEO4J_USERNAME")
	maxRecursion, _ := strconv.ParseInt(os.Getenv("MAX_RECURSION"), 10, 64)
	maxRunners, _ := strconv.ParseInt(os.Getenv("MAX_RUNNERS"), 10, 64)
	pushUsers, _ := strconv.ParseBool(os.Getenv("PUSH_USERS"))
	exportCommand := os.Getenv("EXPORT_COMMAND")
	maxMessageLimit, _ := strconv.ParseInt(os.Getenv("MAX_MESSAGE_FETCH_LIMIT"), 10, 64)

	if maxMessageLimit == 0 {
		maxMessageLimit = 10000
	}

	neo := storage.Neo4jInstance{Username: username, Password: password, URI: uri, DBName: db}
	err = neo.Init()
	if err != nil {
		log.Fatalf("Error on neo4j init: %v", err)
		return
	}
	_ = neo.Clean()

	defer neo.Close()
	defer log.Println("Thanks for running the artio-miner :)")
	manager := miner.Manager{Neo: &neo, MaxRecursion: int(maxRecursion), MaxRunners: int(maxRunners), PushUsers: pushUsers, ExportCommand: exportCommand, MaxMessageCount: maxMessageLimit}

	manager.Run(startingRelays)

}
