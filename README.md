# GraphQL based employee CRUD service example 

## Queries

* create employee
```
curl -g 'http://localhost:8080/graphql?query=mutation+_{createEmp(name:"Tester",job:"dev",mgr:"1",deptno:"1",sal:"100"){EMPNO,ENAME,JOB,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'

    NOTE : 
    1) deptno needs to be a valid department id in Department Table
    2) empno is auto-incremented field
```

* update employee
```
curl -g 'http://localhost:8080/graphql?query=mutation+_{updateEmp(name:"Tester",job:"dev",mgr:"2",deptno:"1",sal:"50",empno:"1"){EMPNO,ENAME,JOB,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'
```

* delete employee
```
curl -g 'http://localhost:8080/graphql?query=mutation+_{delEmp(empno:"3")}'
```

* list employees in a deparment
```
curl -g 'http://localhost:8080/graphql?query={empListInDept(dname:"ENGINEERING"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'
```

* get employee by no
```
curl -g 'http://localhost:8080/graphql?query={empByNo(empno:"1"){EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'
```

* list all employees
```
curl -g 'http://localhost:8080/graphql?query={empListAll{EMPNO,ENAME,JOB,MGR,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'
```