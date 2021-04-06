## AlgoeDB
A lightweight, persistent, NoSQL database written in Go. 

Inspired by the Deno project [AloeDB](https://github.com/Kirlovon/AloeDB). Many thanks to [@Kirlovon](https://github.com/Kirlovon) for the inspiration!

## Features
* 🎉 Simple to use API, similar to [MongoDB](https://www.mongodb.com/)!
* 🚀 Optimized for a large number of operations.
* 📁 Stores data in readable JSON file.
* ⚖  No dependencies outside of the standard library.

## Examples Usage

```go
config := AlgoeDB.DatabaseConfig{Path: "/path/to/file/people.json"}
db, err := AlgoeDB.NewDatabase(&config)
if err != nil {
    log.Fatal(err)
}

people := []map[string]interface{}{}
people = append(people, map[string]interface{}{"name": "Billy", "age": 27})
people = append(people, map[string]interface{}{"name": "Carisa", "age": 26})

err = db.InsertMany(people)
if err != nil {
    log.Fatal(err)
}

query := map[string]interface{}{"name": "Carisa"}
results := db.FindOne(query)

if results != nil {
    fmt.Println("results:", results) //results: [map[age:26 name:Carisa]]
} else {
    fmt.Println("no documents found")
}

query = map[string]interface{}{"age": AlgoeDB.MoreThan(25)}
results = db.FindMany(query)

if results != nil {
    fmt.Println("results:", results) //results: [map[age:27 name:Billy] map[age:26 name:Carisa]]
} else {
    fmt.Println("no documents found")
}
```