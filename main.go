package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
)

func init() {

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

// define schema, with our rootQuery and rootMutation
var schema, schemaErr = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: RootMutation,
})

func main() {

	fmt.Println("[schema error]", schemaErr)
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[in handler]", r.URL.Query())
		result := executeQuery(r.URL.Query()["query"][0], schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Get single todo: curl -g 'http://localhost:8080/graphql?query={todo(id:\"b\"){id,text,done}}'")
	fmt.Println("Create new todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{createEmp(name:\"Jyoti+R\", job : \"dev\", mgr : \"1\" ){EMPNO,ENAME,JOB,SALARY,DEPT}}'")
	fmt.Println("Load todo list: curl -g 'http://localhost:8080/graphql?query={empList{EMPNO,ENAME,JOB, SALARY,DEPT}}'")
	http.ListenAndServe(":8080", nil)
}
