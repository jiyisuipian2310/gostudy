package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// 数据库连接信息
	DB_IP       = "192.168.104.100"
	DB_PORT     = "27017"
	DB_NAME     = "testing"
	DB_USER     = "yull"
	DB_PASSWORD = "123456"
)

func main() {
	//URI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", DB_USER, DB_PASSWORD, DB_IP, DB_PORT, DB_NAME)
	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s", DB_USER, DB_PASSWORD, DB_IP, DB_PORT)
	fmt.Printf("URI: %s\n", URI)

	strBwdaToken := "123456789"

	cliopts := options.Client().ApplyURI(URI)
	cliopts.SetBwdaToken(strBwdaToken)
	client, err := mongo.Connect(context.Background(), cliopts)
	if err != nil {
		panic(err)
	}

	if false {
		if err := client.Ping(context.TODO(), nil); err != nil {
			panic(err)
		}
	}

	// 插入数据
	if false {
		dbcollection := client.Database("testing").Collection("teachers")
		teacher1 := bson.D{{"name", "zhangsan"}, {"age", 120}}
		res, err := dbcollection.InsertOne(context.TODO(), teacher1)
		if err != nil {
			fmt.Printf("InsertOne failed: %s\n", err.Error())
		} else {
			fmt.Println(res.InsertedID)
		}
	}

	// 查询数据
	if true {
		dbcollection := client.Database("testing").Collection("teachers")

		// 查询数据 1:升序  -1:降序
		findOptions := options.Find().SetSort(bson.D{{"age", -1}, {"name", 1}})
		findOptions.SetLimit(3)
		//filter := bson.M{"age": bson.M{"$gt": 20.11}}
		filter := bson.M{"age": bson.M{"$gt": 20.11}, "name": bson.M{"$regex": "张三"}}
		//filter := bson.M{"age": bson.M{"$gt": 20.11}, "name": "teacher3"}
		cur, err := dbcollection.Find(context.TODO(), filter, findOptions)
		//cur, err := dbcollection.Find(context.TODO(), bson.M{}, findOptions)
		if err != nil {
			log.Fatal(err)
		}

		// 关闭游标
		defer cur.Close(context.TODO())

		// 遍历结果
		for cur.Next(context.TODO()) {
			var elem bson.M
			if err = cur.Decode(&elem); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", elem) // 打印每个文档
		}
	}

	if false {
		dbcollection := client.Database("testing").Collection("teachers")
		teacher1 := bson.D{{"name", "teacher1"}}
		res, err := dbcollection.DeleteMany(context.TODO(), teacher1)
		if err != nil {
			fmt.Printf("InsertOne failed: %s\n", err.Error())
		} else {
			fmt.Println("delete count: ", res.DeletedCount)
		}
	}

	client.Disconnect(context.TODO())
}
