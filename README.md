
# GoNotes-App
GoNotes is a simple and efficient notes application built using the Gofr framework in the Go programming language. It provides a robust set of features to manage your notes seamlessly. 


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file
```bash
APP_NAME=Notes-app-using-gofr

HTTP_PORT=8000

DB_HOST=localhost

DB_USER=root

DB_PASSWORD=root123

DB_NAME=test_db

DB_PORT=3306

DB_DIALECT=mysql
```
## Setup MySQL

You will have to run the mysql server and create a database locally using the following docker command:

```bash
docker run --name gofr-mysql -e MYSQL_ROOT_PASSWORD=root123 -e MYSQL_DATABASE=test_db -p 3306:3306 -d mysql:8.0.30
```

Access test_db database and create table notes with columns note_id,title and content

```bash
docker exec -it gofr-mysql mysql -uroot -proot123 test_db -e "CREATE TABLE notes ( note_id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255) NOT NULL UNIQUE, content TEXT NOT NULL);"
```
## Run Locally

Clone the project

```bash
  git clone https://github.com/amanlalwani2/GoNotes-app.git
```

Go to the project directory

```bash
  cd GoNotes-app
```

Install dependencies along with gofr

```bash
  go get ./...
```

Start the server

```bash
  go run main.go
```


## Running Tests

To run all the tests, run the following command(**I have achieved Unit Test Coverage≥90%**)

```bash
  go test ./...
```


## API Reference**(Import the postman collection provided in the repo)**

#### To test run APIS you have to use path:
```bash
  http://localhost:3000
```


#### Get all notes

```
  GET /notes
```
#### Insert note

```
  POST /notes
```

#### Update note

```
  PUT /notes/{id}
```
| Parameter | Description                      |
| :-------- |:-------------------------------- |
| `id`      |**Required**. Id of note to Update |


#### Delete note

```
  DELETE /notes/{id}
```
| Parameter | Description                      |
| :-------- |:-------------------------------- |
| `id`      |**Required**. Id of note to Delete |



## Published Postman Collection

```bash
  https://documenter.getpostman.com/view/31295823/2s9Ykn81p8#5b48a552-4252-45fd-b6d4-5143c6f912bd
```
## Diagrams
### **Use-Case Diagram** 
<img src="https://iili.io/JuXPHhu.png" height="800" width="400" />

### **Sequence Diagram** 
<img src="https://iili.io/JuXU5Ml.md.png" height="800" width="400" />
