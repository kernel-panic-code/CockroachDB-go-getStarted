package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Post type to mirror db records as Go struct
type Post struct {
	ID        int64
	Name      string
	Category  string
	Author    string
	CreatedAt string
	UpdatedAt string
}

var roachConnection *sql.DB

func init() {
	var err error
	roachConnection, err = sql.Open("postgres", "postgresql://root@localhost:26257/blog_db?sslmode=disable")

	if err != nil {
		panic("Could not establish connection to CockroachDB")
	}
}

func createRecord(post Post) error {
	qryString := fmt.Sprintf(
		"INSERT INTO posts (name, category, author, created_at, updated_at) VALUES ('%s', '%s', '%s', NOW(), NOW())",
		post.Name, post.Category, post.Author)

	fmt.Printf("Query: %v\n", qryString)
	_, err := roachConnection.Exec(qryString)
	return err
}

func readRecords() []Post {
	var result []Post
	rows, err := roachConnection.Query("select * FROM posts;")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Name, &post.Author, &post.Category, &post.CreatedAt, &post.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		result = append(result, post)
	}

	return result
}

func updateRecord(fieldName string, newVal string, condFieldName string, condFieldValue int64) error {
	qryString := fmt.Sprintf("UPDATE posts set %s='%s' WHERE %s=%d", fieldName, newVal, condFieldName, condFieldValue)

	fmt.Printf("Query: %v\n", qryString)
	_, err := roachConnection.Exec(qryString)

	return err
}

func deleteRecord(id int64) error {
	qryString := fmt.Sprintf("DELETE from posts WHERE id = %d", id)
	fmt.Printf("Query: %v\n", qryString)

	_, err := roachConnection.Exec(qryString)
	return err
}

func main() {
	Post1 := Post{Name: "foo", Category: "Category1", Author: "User1"}
	Post2 := Post{Name: "bar", Category: "Category2", Author: "User2"}

	err := createRecord(Post1)
	if err != nil {
		panic("Could not create record")
	}

	createRecord(Post2)
	if err != nil {
		panic("Could not create record")
	}

	result := readRecords()
	fmt.Println(result)

	err = updateRecord("name", "fooBar", "id", result[0].ID)
	if err != nil {
		panic("Could not update")
	}

	fmt.Println("----------------AFTER UPDATE---------------------")

	result = readRecords()
	fmt.Println(result)

	err = deleteRecord(result[0].ID)
	if err != nil {
		panic("Could not update")
	}

	fmt.Println("----------------AFTER DELETE---------------------")
	result = readRecords()
	fmt.Println(result)
}
