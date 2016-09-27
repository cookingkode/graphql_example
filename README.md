# GraphQL based employee CRUD service example 

## Queries

* create employee
```
curl -g 'http://localhost:8080/graphql?query=mutation+_{createEmp(name:"Tester",job:"dev",mgr:"1",deptno:"1",sal:"100"){EMPNO,ENAME,JOB,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'
```
NOTE : 
- deptno needs to be a valid department id in Department Table
- empno is auto-incremented field

* update employee
```
curl -g 'http://localhost:8080/graphql?query=mutation+_{updateEmp(name:"Tester",job:"dev",mgr:"2",deptno:"1",sal:"50",empno:"1"){EMPNO,ENAME,JOB,SALARY,DEPT{DEPTNO,DNAME,LOC}}}'
```
