# MySQL Database


## go-sql-driver/mysql package

In order to use a specific database we need to install the adecuate driver via `go get -u github.com/go-sql-driver/mysql`.

## Connecting to a MySQL Database

- If you do not have a running database, you can get up one up and running with Docker.

```go
import "database/sql"
import _ "go-sql-driver/mysql"

// Configure the database connection (always check errors)
db, err := sql.Open("mysql", "username:password@(127.0.0.1:3306)/dbname?parseTime=true")

if err != nil {
  // handle it
}

// Initialize the first connection to the database to see if everything works correctly.
// Make sure to check the error
err = db.Ping()
if err != nil {
  fmt.Println(err)
  panic
}
```

## Creating our first database

- Every data entry in our database is stored in a specific table.
- A database table consists of columns and rows.
- The columns give each data entry a label and specify the type of it.
- The rows are the inserted data values.

```sql
CREATE TABLE users (
  id INT AUTO_INCREMENT,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at DATETIME,
  PRIMARY KEY (id)
);
```

To execute this query use:
```go
query := `
    CREATE TABLE users (
      id INT AUTO_INCREMENT,
      username TEXT NOT NULL,
      password TEXT NOT NULL,
      created_at DATETIME,
      PRIMARY KEY (id)
    );`
// Executes the SQL query in our database.
// Check err to ensure there was no error.
_, err := db.Exec(query)
```

## Inserting our first user

- By default Go uses prepared statements for inserting dynamic data into our SQL queries.
- Prepared statements are considered a secure way to interact with the database without risk of a SQL injection.
- `INSERT INTO users (username, password, created_at) VALUES (?,?,?)`

To use this query:
```go
import "time"

username := "johndoe"
password := "secret"
createdAt := time.Now()

// Inserts our data into the users table and returns with the result and a possible error.
// The result contains information about the last inserted id and the count of rows this query affects.
INSERTION_QUERY := `INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`
result, err := db.Exec(INSERTION_QUERY, username, password, createdAt)

// To grab the newly generated id for your user
userId, err := result.LastInsertedId()
```

## Querying our users table

- Now that we have a user in our table, we want to query it and get back all of its information.
- Two possibilities in Go
  + `db.Query` which can query multiple rows for us to iterate over.
  + `db.QueryRow` in case we only want to query a specific row.
- Querying a specific row works basically like every other SQL command.

```sql
SELECT id, username, password, created_at FROM users WHERE id = ?
```

In go we first declare some variables to store our data in and then query a single database row:
```go
var (
  id int
  username string
  password string
  createdAt time.Time
)

// Query the database and scan the values into our variables.
// Check for errors.
query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt)

```

## Querying all users

- We can use the SQL command above and trim the `WHERE` clause.

```sql
SELECT id, username, password, created_at FROM users
```

```go
type user struct {
    id        int
    username  string
    password  string
    createdAt time.Time
}

rows, err := db.Query(`SELECT id, username, password, created_at FROM users`) // check err
defer rows.Close()

var users []user
for rows.Next() {
    var u user
    // This will yield a JSON structure
    err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt) // check err
    users = append(users, u)
}
err := rows.Err() // check err
```

## Deleting a user from the table.


```go
_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1) // check err
```

## Code

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/root?parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    { // Create a new table
        query := `
            CREATE TABLE users (
                id INT AUTO_INCREMENT,
                username TEXT NOT NULL,
                password TEXT NOT NULL,
                created_at DATETIME,
                PRIMARY KEY (id)
            );`

        if _, err := db.Exec(query); err != nil {
            log.Fatal(err)
        }
    }

    { // Insert a new user
        username := "johndoe"
        password := "secret"
        createdAt := time.Now()

        result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
        if err != nil {
            log.Fatal(err)
        }

        id, err := result.LastInsertId()
        fmt.Println(id)
    }

    { // Query a single user
        var (
            id        int
            username  string
            password  string
            createdAt time.Time
        )

        query := "SELECT id, username, password, created_at FROM users WHERE id = ?"
        if err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt); err != nil {
            log.Fatal(err)
        }

        fmt.Println(id, username, password, createdAt)
    }

    { // Query all users
        type user struct {
            id        int
            username  string
            password  string
            createdAt time.Time
        }

        rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        var users []user
        for rows.Next() {
            var u user

            err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
            if err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("%#v", users)
    }

    {
        _, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1)
        if err != nil {
            log.Fatal(err)
        }
    }
}
```
