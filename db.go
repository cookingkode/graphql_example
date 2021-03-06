package main

import (
	"database/sql"
	"errors"
)

import _ "github.com/go-sql-driver/mysql"

var (
	ERR_DEPT_NOT_FOUND = errors.New("Department not found")
	ERR_EMP_NOT_FOUND  = errors.New("Employee not found")
)

var (
	db                   *sql.DB
	getEmployeeStmt      *sql.Stmt
	listEmployesDeptStmt *sql.Stmt
	createEmployeeStmt   *sql.Stmt
	listAllEmployesStmt  *sql.Stmt
	getDepartmentStmt    *sql.Stmt
	updateEmpStmt        *sql.Stmt
	delEmpStmt           *sql.Stmt
)

func dbInit(dbstring string) {
	var err error

	if db, err = sql.Open("mysql", dbstring); err != nil {
		panic(err)
	}

	if createEmployeeStmt, err = db.Prepare("INSERT INTO  EMPLOYEE ( ENAME, JOB, SALARY, MGR, DEPTNO) VALUES ( ?, ?, ?, ?, ?)"); err != nil {
		panic(err)
	}

	if getEmployeeStmt, err = db.Prepare(
		"SELECT EMPNO, ENAME, JOB, MGR, SALARY , D.DEPTNO, D.DNAME, D.LOC FROM EMPLOYEE E LEFT JOIN DEPARTMENT D ON E.DEPTNO = D.DEPTNO  WHERE EMPNO =?"); err != nil {
		panic(err)
	}

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

	if delEmpStmt, err = db.Prepare(
		"DELETE FROM EMPLOYEE   WHERE EMPNO = ?"); err != nil {
		panic(err)
	}

}

func getDepartment(deptno string) (*Department, error) {
	var dept Department

	if err := getDepartmentStmt.QueryRow(deptno).Scan(
		&dept.DEPTNO,
		&dept.DNAME,
		&dept.LOC); err != nil {
		return nil, ERR_DEPT_NOT_FOUND
	}

	return &dept, nil
}

func getEmployee(empno string) (*Employee, error) {
	var emp Employee

	if err := getEmployeeStmt.QueryRow(empno).Scan(
		&emp.EMPNO,
		&emp.ENAME,
		&emp.JOB,
		&emp.MGR,
		&emp.SALARY,
		&emp.DEPT.DEPTNO,
		&emp.DEPT.DNAME,
		&emp.DEPT.LOC); err != nil {
		return nil, ERR_EMP_NOT_FOUND
	}

	return &emp, nil

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

func delEmployee(empno string) error {
	_, err := delEmpStmt.Exec(empno)
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
		toret = append(toret, emp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return toret, nil
}
