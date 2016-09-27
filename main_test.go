package main

import (
	"flag"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/kr/pretty"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
)

// var testSchema graphql.Schema

func init() {

	rand.Seed(time.Now().UnixNano())

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Diff(a, b interface{}) []string {
	return pretty.Diff(a, b)
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	flag.Parse()
	//init DB
	dbInit(*dbConnectString)

	m.Run()
}

func createEmployeeRecord(name, dept string, t *testing.T) (string, *graphql.Result) {
	createQuery := "mutation _{createEmp(name:\"%s\",job:\"dev\",mgr:\"1\",deptno:\"%s\",sal:\"100\"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}"
	// Deptno 1 exists and is ENGINEERING
	validCreateQuery := fmt.Sprintf(createQuery, name, dept)
	result := executeQuery(validCreateQuery, schema)

	resp, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("CreateEmployeeTest failed , could not retrieve employee")
		return "", result
	}

	dat, ok := resp["createEmp"].(map[string]interface{})
	if !ok {
		return "", result
	}
	emp := dat["EMPNO"].(string)

	return emp, result

}

func checkResult(name, dept, empno, queryname, dname, loc string, result *graphql.Result, t *testing.T) {
	expected := &graphql.Result{
		Data: map[string]interface{}{
			queryname: map[string]interface{}{
				"DEPT": map[string]interface{}{
					"DEPTNO": dept, "DNAME": dname, "LOC": loc,
				},
				"EMPNO":  empno,
				"ENAME":  name,
				"MGR":    "1",
				"JOB":    "dev",
				"SALARY": "100",
			},
		},
		Errors: nil,
	}

	if !reflect.DeepEqual(result.Data, expected.Data) {
		t.Fatalf("test failed %+v expected %+v diff %v", result, expected, Diff(expected, result))
	}

}

func TestCreateEmployee(t *testing.T) {
	name := randStringRunes(8)
	validDept := "1" //dname = "ENGINEERING", loc = "NY"
	empno, result := createEmployeeRecord(name, validDept, t)

	if len(result.Errors) > 0 {
		t.Fatalf("CreateEmployeeTest failed %+v", result.Errors)
	}

	checkResult(name, validDept, empno, "createEmp", "ENGINEERING", "NY", result, t)
}

func TestCreateEmployeeInvalidDept(t *testing.T) {
	name := randStringRunes(8)
	validDept := "100" // NON-Existnet
	_, result := createEmployeeRecord(name, validDept, t)

	// there should be an error
	if len(result.Errors) == 0 {
		t.Fatalf("InvalidDeptCreateEmployeeTest failed %+v", result.Errors)
	}

}

func TestUpdateEmployee(t *testing.T) {

	name := randStringRunes(8)
	dept := "1"
	empno, _ := createEmployeeRecord(name, dept, t)

	// change dept
	dept = "2" // dname = "MARKETING", loc = "LA"
	updateQFmt := "mutation _{updateEmp(empno:\"%s\",name:\"%s\",job:\"dev\",mgr:\"1\",deptno:\"%s\",sal:\"100\"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}"
	updateQuery := fmt.Sprintf(updateQFmt, empno, name, dept)
	updateResult := executeQuery(updateQuery, schema)

	if len(updateResult.Errors) > 0 {
		t.Fatalf("UpdateEmpoyee failed %+v", updateResult.Errors)
	}

	checkResult(name, dept, empno, "updateEmp", "MARKETING", "LA", updateResult, t)

}

func TestDelEmployee(t *testing.T) {
	name := randStringRunes(8)
	dept := "1"
	empno, _ := createEmployeeRecord(name, dept, t)

	// delete
	delQFmt := "mutation _{delEmp(empno:\"%s\")}"
	delQuery := fmt.Sprintf(delQFmt, empno)
	res := executeQuery(delQuery, schema)
	if len(res.Errors) > 0 {
		t.Fatalf("Del failed %+v", res.Errors)
	}

}

func TestGetEmployee(t *testing.T) {
	name := randStringRunes(8)
	dept := "1"
	empno, _ := createEmployeeRecord(name, dept, t)

	getQueryFmt := "{empByNo(empno:\"%s\"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}"
	getRes := executeQuery(fmt.Sprintf(getQueryFmt, empno), schema)
	if len(getRes.Errors) > 0 {
		t.Fatalf("Get failed %+v", getRes.Errors)
	}

	checkResult(name, dept, empno, "empByNo", "ENGINEERING", "NY", getRes, t)

}

func TestGetInvalidEmployee(t *testing.T) {
	// delete all entries
	_, err := db.Exec("DELETE FROM EMPLOYEE")
	if err != nil {
		t.Fatalf("could not delete existing entries")
	}

	// create an reccord
	name := randStringRunes(8)
	dept := "1"
	empno, _ := createEmployeeRecord(name, dept, t)

	// this wont exist since we droped
	nonExistent := strings.Repeat(empno, 2)

	getQueryFmt := "{empByNo(empno:\"%s\"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}"
	getRes := executeQuery(fmt.Sprintf(getQueryFmt, nonExistent), schema)
	// there should be an error
	if len(getRes.Errors) == 0 {
		t.Fatalf("Get failed %+v", getRes.Errors)
	}

}

func ensureNoEmployees(queryname string, result *graphql.Result, no int, t *testing.T) {
	resp, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("\nlistall failed %+v\n ", result.Data)
	}

	emps, ok := resp[queryname].([]interface{})
	if !ok {
		t.Fatalf("\nlistall failed %+v\n ", emps)
	}

	if len(emps) != no {
		t.Fatalf("emp no mismatch ")
	}
}

func TestListAllEmployees(t *testing.T) {
	// delete all entries
	_, err := db.Exec("DELETE FROM EMPLOYEE")
	if err != nil {
		t.Fatalf("could not delete existing entries")
	}

	name1 := randStringRunes(8)
	name2 := randStringRunes(8)
	dept := "1"

	// create two employees
	createEmployeeRecord(name1, dept, t)
	createEmployeeRecord(name2, dept, t)

	listAll := "{empListAll{EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}"
	getRes := executeQuery(listAll, schema)
	if len(getRes.Errors) > 0 {
		t.Fatalf("Get failed %+v", getRes.Errors)
	}

	ensureNoEmployees("empListAll", getRes, 2, t)

}

func TestListAllEmployeesInDept(t *testing.T) {
	// delete all entries
	_, err := db.Exec("DELETE FROM EMPLOYEE")
	if err != nil {
		t.Fatalf("could not delete existing entries")
	}

	name1 := randStringRunes(8)
	name2 := randStringRunes(8)
	dept := "1" // ENGINEERING

	// create two employees
	createEmployeeRecord(name1, dept, t)
	createEmployeeRecord(name2, dept, t)

	listAllInDept := "{empListInDept(dname:\"%s\"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}"
	getRes := executeQuery(fmt.Sprintf(listAllInDept, "ENGINEERING"), schema)
	if len(getRes.Errors) > 0 {
		t.Fatalf("Get failed %+v", getRes.Errors)
	}

	ensureNoEmployees("empListInDept", getRes, 2, t)

}
