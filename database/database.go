package database

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// type questionPost struct {
// 	inputString string
// }

type IntentProvider struct {
	policy					*bluemonday.Policy
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
	Message string	`json:"message"`
	El      string	`json:"el"`
}

func connectToPool(connectionString string) *sqlx.DB {
	fmt.Println("connecting to pool")
	return sqlx.MustConnect("postgres", connectionString)
}

func SetIntentProvider() *IntentProvider {
	p := bluemonday.StrictPolicy()
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", dbUser, dbHost, dbName)
	fmt.Println("connectionString", connectionString)
	pool := connectToPool(connectionString)
	return &IntentProvider{policy: p, host: dbHost, user: dbUser, port: os.Getenv("DB_PORT"), name: dbName, connectionString: connectionString, pool: pool}
}

// func setintentProvider() *database.IntentProvider {
// 	dbHost := os.Getenv("DB_HOST")
// 	dbName := os.Getenv("DB_NAME")
// 	connectionString := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", dbHost, dbName)
// 	return &intentProvider{host: dbHost, user: os.Getenv("DB_USER"), port: os.Getenv("DB_PORT"), name: dbName, connectionString: connectionString}
// }

func (ip *IntentProvider) PredictionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var question string
	err := decoder.Decode(&question)
	if err != nil {
			panic(err)
	}
	log.Println(question)


	answer, _ := ip.getPrediction(question)
	fmt.Println("answer: ", answer)

	answerBytes, err := json.Marshal(answer)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(answerBytes)
}

// NEVER CALL THIS FUNCTION WITHOUT FIRST SANITIZING IT
func (ip *IntentProvider) getPrediction(question string) ([]byte, error) {
	fmt.Println("question  inside", question)
	question = ip.policy.Sanitize(question)
	fmt.Println("question: ", question)

	cmd := exec.Command("python", "-c", fmt.Sprintf("from ai.prediction import chatbot_response; print(chatbot_response('%v'))", question))
	out, _ := cmd.CombinedOutput()
	fmt.Println("out", out)
	fmt.Println("string out", string(out))
	// if err != nil {
	// 	fmt.Println("Error executing prediction", err)
	// }
	re := regexp.MustCompile(`{[^\\"]+}`)
	return re.Find([]byte(string(out))), nil
}

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
