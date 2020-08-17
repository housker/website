package database

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type helpers interface {
	GetPatterns()
	GetResponses()
}

type IntentProvider struct {
	host             string
	user             string
	port             string
	name             string
	connectionString string
	pool             *sqlx.DB
}

type Intents struct {
	Intents []Intent
}

type Intent struct {
	Tag       string
	Patterns  []string
	Responses []Response
}

type Response struct {
	Message string	// `json:"myName"`
	El      string
}

func connectToPool(connectionString string) *sqlx.DB {
	return sqlx.MustConnect("postgres", connectionString)
}

// func (ip *database.IntentProvider) connectToPool() *sqlx.DB {
// 	connectionString := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", ws.dbConfig.host, ws.dbConfig.name)
// 	return sqlx.MustConnect("postgres", connectionString)
// }

func SetIntentProvider() *IntentProvider {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", dbHost, dbName)
	pool := connectToPool(connectionString)
	return &IntentProvider{host: dbHost, user: os.Getenv("DB_USER"), port: os.Getenv("DB_PORT"), name: dbName, connectionString: connectionString, pool: pool}
}

// func setintentProvider() *database.IntentProvider {
// 	dbHost := os.Getenv("DB_HOST")
// 	dbName := os.Getenv("DB_NAME")
// 	connectionString := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", dbHost, dbName)
// 	return &intentProvider{host: dbHost, user: os.Getenv("DB_USER"), port: os.Getenv("DB_PORT"), name: dbName, connectionString: connectionString}
// }

func (ip *IntentProvider) TagHandler(w http.ResponseWriter, r *http.Request) {
	tags, _ := ip.GetTags()

	tagBytes, err := json.Marshal(tags)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(tagBytes)
}

func (ip *IntentProvider) GetTags() ([]string, error) {
	// db := ip.connectToPool()
	var tags []string
	statement := `SELECT name FROM tags`
	err := ip.pool.Select(&tags, statement)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		return nil, err
	}
	return tags, nil
}

// HELPER FUNCTION _____________
// func (ip *IntentProvider) generateIntents() Intents {
// 	var groups []Intent

// 	tags, _ := ip.getTags()

// 	for _, tag := range tags {
// 		patterns, _ := ip.getPatterns(tag)
// 		responses, _ := ip.getResponses(tag)
// 		group := Intent{Tag: tag, Patterns: patterns, Responses: responses}
// 		groups = append(groups, group)
// 	}
// 	return Intents{Intents: groups}
// }

func (ip *IntentProvider) GetResponses(tag string) ([]Response, error) {
	// db := ip.connectToPool()

	var (
		responses []Response
		message   string
		el        string
	)

	stmt, err := ip.pool.Prepare("SELECT message, el FROM responses WHERE tag = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(tag)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&message, &el)
		if err != nil {
			log.Fatal(err)
		}
		response := Response{Message: message, El: el}
		responses = append(responses, response)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return responses, nil
}

func (ip *IntentProvider) GetPatterns(tag string) ([]string, error) {
	// db := ip.connectToPool()
	var patterns []string
	statement := "SELECT pattern FROM patterns WHERE tag = $1"
	err := ip.pool.Select(&patterns, statement, tag)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		return nil, err
	}
	return patterns, nil
}

// _____________________________

//  Helper functions ---------------------------------------------------------------
// func (ip *database.IntentProvider) getTags() ([]string, error) {
// 	db := ws.connectToPool()
// 	var tags []string
// 	statement := `SELECT name FROM tags`
// 	err := db.Select(&tags, statement)
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("Error: %v", err))
// 		return nil, err
// 	}
// 	return tags, nil
// }

// func (ip *database.IntentProvider) getPatterns(tag string) ([]string, error) {
// 	db := ws.connectToPool()
// 	var patterns []string
// 	statement := "SELECT pattern FROM patterns WHERE tag = $1"
// 	err := db.Select(&patterns, statement, tag)
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("Error: %v", err))
// 		return nil, err
// 	}
// 	return patterns, nil
// }

// func (ip *database.IntentProvider) getResponses(tag string) ([]Response, error) {
// 	db := ws.connectToPool()

// 	var (
// 		responses []Response
// 		message   string
// 		el        string
// 	)

// 	stmt, err := db.Prepare("SELECT message, el FROM responses WHERE tag = $1")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	rows, err := stmt.Query(tag)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		err = rows.Scan(&message, &el)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		response := Response{Message: message, El: el}
// 		responses = append(responses, response)
// 	}
// 	if err = rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return responses, nil
// }

// ___________________________________

// func init() {
// 	if os.Getenv("DB_HOST") {
// 		return
// 	}
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// }

// func database() {
// 	ip := setIntentProvider()

// }
