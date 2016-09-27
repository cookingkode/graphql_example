package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
)

// define schema, with our rootQuery and rootMutation
var schema, schemaErr = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: RootMutation,
})

var dbConnectString = flag.String("db", "root@/company", "sql db name connect string, defaults to local root")

func main() {

	if schemaErr != nil {
		panic(schemaErr)
	}

	//init DB
	dbInit(*dbConnectString)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[in handler]", r.URL.Query())
		result := executeQuery(r.URL.Query()["query"][0], schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Graphql server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
