package main

import (
	"database/sql"
	"fmt"
)

import _ "github.com/go-sql-driver/mysql"

const (
	LIST_ALL_EMPS_IN_DEPT = "SELECT * FROM EMPLOYEE WHERE DEPTNO= (SELECT DEPTNO FROM  DEPARTMENT WHERE DNAME =?)"
	LIST_ALL_EMPS         = "SELECT * FROM EMPLOYEE"
	GET_EMPLOYEE          = "SELECT * FROM EMPLOYEE WHERE EMPNO = ?"
)

var (
	db *sql.DB
	//getEmployee    *sql.Stmt
	//setEmployee    *sql.Stmt
	listEmployesDeptStmt *sql.Stmt
	createEmployeeStmt   *sql.Stmt
	listAllEmployesStmt  *sql.Stmt
	getDepartmentStmt    *sql.Stmt
	updateEmpStmt        *sql.Stmt
)

func init() {
	var err error

	if db, err = sql.Open("mysql", "root@/company"); err != nil {
		panic(err)
	}

	if createEmployeeStmt, err = db.Prepare("INSERT INTO  EMPLOYEE ( ENAME, JOB, SALARY, MGR, DEPTNO) VALUES ( ?, ?, ?, ?, ?)"); err != nil {
		panic(err)
	}
	/*
			if getEmployee, err := db.Prepare("SELECT * FROM EMPLOYEE WHERE EMPNO = ?"); err != nil {
				panic(err)
			}

		if listEmployesDeptStmt, err = db.Prepare("SELECT EMPNO, ENAME, JOB, MGR, SALARY FROM EMPLOYEE WHERE DEPTNO= (SELECT DEPTNO FROM  DEPARTMENT WHERE DNAME =?)"); err != nil {
			panic(err)
		}*/

	if listEmployesDeptStmt, err = db.Prepare(
		"SELECT EMPNO, ENAME, JOB, MGR, SALARY , D.DEPTNO, D.DNAME, D.LOC FROM EMPLOYEE E LEFT JOIN DEPARTMENT D ON E.DEPTNO = D.DEPTNO  WHERE DNAME =?"); err != nil {
		panic(err)
	}

	if listAllEmployesStmt, err = db.Prepare(
		"SELECT EMPNO, ENAME, JOB, MGR, SALARY , D.DEPTNO, D.DNAME, D.LOC FROM EMPLOYEE E LEFT JOIN DEPARTMENT D ON E.DEPTNO = D.DEPTNO"); err != nil {
		panic(err)
	}

	if getDepartmentStmt, err = db.Prepare(
		"SELECT DEPTNO, DNAME, LOC  FROM  DEPARTMENT  WHERE DEPTNO = ?"); err != nil {
		panic(err)
	}

	if updateEmpStmt, err = db.Prepare(
		"UPDATE EMPLOYEE SET ENAME = ?, JOB = ?, SALARY = ?, MGR = ?, DEPTNO = ?  WHERE EMPNO = ?"); err != nil {
		panic(err)
	}

}

/*
func createEmployee(emp Employee) error {
	_, err := createEmployeeStmt.Exec(emp.EMPNO, emp.ENAME, emp.JOB, emp.SALARY, emp.MGR, emp.DEPT.DEPTNO)
	return err
}
*/

func getDepartment(deptno string) (*Department, error) {
	var dept Department

	if err := getDepartmentStmt.QueryRow(deptno).Scan(
		&dept.DEPTNO,
		&dept.DNAME,
		&dept.LOC); err != nil {
		return nil, err
	}

	return &dept, nil
}

func createEmployee(name, job, salary, mgr, deptno string) (int64, error) {

	res, err := createEmployeeStmt.Exec(name, job, salary, mgr, deptno)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateEmployee(empno, name, job, salary, mgr, deptno string) error {
	_, err := updateEmpStmt.Exec(name, job, salary, mgr, deptno, empno)
	return err
}

func listAllEmployees() ([]Employee, error) {
	var toret []Employee

	rows, err := listAllEmployesStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emp Employee
		if err := rows.Scan(
			&emp.EMPNO,
			&emp.ENAME,
			&emp.JOB,
			&emp.MGR,
			&emp.SALARY,
			&emp.DEPT.DEPTNO,
			&emp.DEPT.DNAME,
			&emp.DEPT.LOC); err != nil {
			return nil, err
		}
		fmt.Printf(" row %v \n", emp)
		toret = append(toret, emp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return toret, nil
}

func listEmployeesInDept(DNAME string) ([]Employee, error) {

	var toret []Employee

	rows, err := listEmployesDeptStmt.Query(DNAME)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emp Employee
		if err := rows.Scan(
			&emp.EMPNO,
			&emp.ENAME,
			&emp.JOB,
			&emp.MGR,
			&emp.SALARY,
			&emp.DEPT.DEPTNO,
			&emp.DEPT.DNAME,
			&emp.DEPT.LOC); err != nil {
			return nil, err
		}
		fmt.Printf(" row %v \n", emp)
		toret = append(toret, emp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return toret, nil
}
