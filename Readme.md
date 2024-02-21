docker pull skaliarman/appscode-api-server:latest

- API link:
- [http://localhost:8080](http://localhost:8080/)
- http://localhost:8080/employee

Data Structure
```json
{
        "id": 1,
        "name": "Arman",
        "salary": "5000"
 }
```
cobra-cli:

cobra-cli init

cobra-cli start


kubectl create secret generic my-secret -n database --from-literal=root-passwo
rd="arman" --dry-run=client -oyaml>dbsecret.yaml

kubectl get endpoints -n database

kubectl get pods -n database -owide

kubectl expose deployment appscode -n appscode --por
t 8080 --dry-run=client -oyaml>appscode-service.yaml


kubectl exec -it -n database pods/mysql-5d8f95d66d-w65jr -- bash

create database appscode;

CREATE TABLE appscode.employee (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100), salary VARCHAR(100));

INSERT INTO appscode.employee ( name, salary) VALUES ( 'Arman', '5000');

INSERT INTO appscode.employee ( name, salary) VALUES ( 'sourav', '5000');

kubectl port-forward -n appscode svc/appscode 8080




