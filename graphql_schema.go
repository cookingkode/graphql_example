package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

type Department struct {
	DEPTNO string `json:"DEPTNO"`
	DNAME  string `json:"DNAME'`
	LOC    string `json:"LOC'`
}

type Employee struct {
	EMPNO  string     `json:"EMPNO"`
	ENAME  string     `json:"ENAME"`
	JOB    string     `json:"JOB"`
	SALARY string     `json:"SALARY"`
	MGR    string     `json:"MGR"`
	DEPT   Department `json:"DEPT"`
}

var deptType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Dept",
	Fields: graphql.Fields{
		"DEPTNO": &graphql.Field{
			Type: graphql.String,
		},
		"DNAME": &graphql.Field{
			Type: graphql.String,
		},
		"LOC": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// define custom GraphQL ObjectType `empType` for our Golang struct `Employee`
// Note that
// - the fields  map with the json tags for the fields in our struct
// - the field types match the field type in our struct
var empType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Emp",
	Fields: graphql.Fields{
		"EMPNO": &graphql.Field{
			Type: graphql.String,
		},
		"ENAME": &graphql.Field{
			Type: graphql.String,
		},
		"JOB": &graphql.Field{
			Type: graphql.String,
		},
		"SALARY": &graphql.Field{
			Type: graphql.String,
		},
		"DEPT": &graphql.Field{
			Type: graphql.NewNonNull(deptType),
		},
	},
})

var tmpList []Employee

// root mutation
var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		/*
		   curl -g 'http://localhost:8080/graphql?query=mutation+_{createEmp(name:"Jyoti",job:"dev",mgr:"1"){EMPNO,ENAME,JOB,SALARY,DEPT}}
		*/
		"createEmp": &graphql.Field{
			Type: empType, // the return type for this field
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"job": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"mgr": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				fmt.Println("[in resolve]")

				// marshall and cast the argument value
				name, _ := params.Args["name"].(string)
				job, _ := params.Args["job"].(string)
				mgr, _ := params.Args["mgr"].(string)

				// perform mutation operation here
				// for e.g. create a Todo and save to DB.
				newEmp := Employee{
					EMPNO:  "1",
					JOB:    job,
					ENAME:  name,
					MGR:    mgr,
					SALARY: "0",
					DEPT: Department{
						DEPTNO: "1",
						DNAME:  "dp",
						LOC:    "blr",
					},
				}

				// create in DB
				createEmployee(newEmp)

				//tmpList = append(tmpList, newEmp)

				// return the new Emp object
				// Note here that
				// - we are returning a struct instance here
				// - we previously specified the return Type to be `empType`
				return newEmp, nil
			},
		},
	},
})

// root query
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		/*
		   curl -g 'http://localhost:8080/graphql?query={empList{EMPNO,ENAME,JOB, SALARY,DEPT}}'
		*/
		"empListInDept": &graphql.Field{
			Type:        graphql.NewList(empType),
			Description: "List of employees in a Department",
			Args: graphql.FieldConfigArgument{
				"dname": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dname, _ := p.Args["dname"].(string)
				ret, err := listEmployeesInDept(dname)
				fmt.Println("[list]", ret, err)
				return ret, err
			},
		},
	},
})
