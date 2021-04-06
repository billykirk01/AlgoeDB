## AlgoeDB
A lightweight, persistent, NoSQL database written in Go. 

Inspired by the Deno project [AloeDB](https://github.com/Kirlovon/AloeDB). Many thanks to [@Kirlovon](https://github.com/Kirlovon) for the inspiration!

## Features
* 🎉 Simple to use API, similar to [MongoDB](https://www.mongodb.com/)!
* 📁 Stores data in readable JSON file.
* 🚀 Optimized for a large number of operations.
* ⚖  No dependencies outside of the standard library.

## Examples Usage

```go
type People []map[string]interface{}

type Person map[string]interface{}

config := AlgoeDB.DatabaseConfig{Path: "./people.json"}
db, err := AlgoeDB.NewDatabase(&config)
if err != nil {
    log.Fatal(err)
}

people := People{}
people = append(people, Person{"name": "Billy", "age": 27})
people = append(people, Person{"name": "Carisa", "age": 26})

err = db.InsertMany(people)
if err != nil {
    log.Fatal(err)
}

query := Person{"name": "Carisa"}
result := db.FindOne(query)
if result != nil {
    fmt.Println(result) // [map[age:26 name:Carisa]]
} else {
    fmt.Println("no documents found")
}

query = Person{"age": AlgoeDB.MoreThan(25)}
results := db.FindMany(query)
if results != nil {
    fmt.Println(results) //[map[age:27 name:Billy] map[age:26 name:Carisa]]
} else {
    fmt.Println("no documents found")
}
```