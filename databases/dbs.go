package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(""))
	if err != nil {
		log.SetPrefix("CONN SETUP ::: ")
		log.SetFlags(0)
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println("DISCONNECTED .....")
			log.Fatal(err)
		}
	}()
	fmt.Println("CONNECTION IS UP .....")

	//list database
	result, err := client.ListDatabaseNames(
		nil,
		bson.D{})
	fmt.Println(result)
	if err != nil {
		log.Fatal(err)
	}

	for _, db := range result {
		fmt.Println(db)
	}

	//get database by name
	db := client.Database("test")
	//dbResults, err := db.ListCollectionNames(
	//	nil,
	//	bson.D{{}})
	//fmt.Println(dbResults)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, col := range dbResults {
	//	fmt.Println("DATABASE :::: " + col)
	//}

	col := db.Collection("users")

	ops := options.Find()
	cursor, err := col.Find(nil, bson.D{}, ops)
	if err != nil {
		log.SetPrefix("Collection ::: ")
		log.SetFlags(0)
		log.Fatal(err)
	}

	var colResults []bson.M

	if err = cursor.All(nil, &colResults); err != nil {
		log.SetPrefix("Collection CURSOR ERROR ::: ")
		log.SetFlags(0)
		log.Fatal(err)
	}

	data, err := json.Marshal(colResults)
	write, err := log.Writer().Write(data)
	if err != nil {
		return
	}

	fmt.Println(write)

}
