```
package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

/*
参考: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
*/

type Database struct {
	Mongo *mongo.Client
}

var db *Database

//初始化

func Init(mongoUrl string) {
	db = &Database{
		Mongo: SetConnect(mongoUrl),
	}
}

//连接设置

func SetConnect(mongoUrl string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接池
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl).SetMaxPoolSize(20))
	if err != nil {
		log.Println(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	return client
}

//创建数据库连接信息对象

type mgo struct {
	database   string
	collection string
}

func NewMgo(database, collection string) *mgo {
	return &mgo{
		database,
		collection,
	}
}

//插入单个文档

func (m *mgo) InsertOne(value interface{}) *mongo.InsertOneResult {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	insertResult, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		log.Fatal(err)
	}
	return insertResult
}

//插入多个文档

func (m *mgo) InsertMany(values []interface{}) *mongo.InsertManyResult {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	result, err := collection.InsertMany(context.TODO(), values)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

//删除文档

func (m *mgo) Delete(key string, value interface{}) int64 {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	return count.DeletedCount
}

//删除多个文档

func (m *mgo) DeleteMany(key string, value interface{}) int64 {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}

	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return count.DeletedCount
}

//更新单个文档

func (m *mgo) UpdateOne(filter, update interface{}) int64 {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return result.UpsertedCount
}

//更新多个文档

func (m *mgo) UpdateMany(filter, update interface{}) int64 {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return result.UpsertedCount
}

// 查询单个文档

func (m *mgo) FindOne(key string, value interface{}) *mongo.SingleResult {
	client := db.Mongo
	collection, e := client.Database(m.database).Collection(m.collection).Clone()
	if e != nil {
		log.Fatal(e)
	}
	filter := bson.D{{key, value}}
	singleResult := collection.FindOne(context.TODO(), filter)
	return singleResult
}

//查询多个文档

func (m *mgo) FindMany(filter interface{}) (*mongo.Cursor, error) {
	client := db.Mongo
	collection, e := client.Database(m.database).Collection(m.collection).Clone()
	if e != nil {
		log.Fatal(e)
	}
	return collection.Find(context.TODO(), filter)
}

//多条件查询

func (m *mgo) FindManyByFilters(filter interface{}) (*mongo.Cursor, error) {
	client := db.Mongo
	collection, e := client.Database(m.database).Collection(m.collection).Clone()
	if e != nil {
		log.Fatal(e)
	}
	return collection.Find(context.TODO(), bson.M{"$and": filter})
}

//查询集合里有多少数据

func (m *mgo) CollectionCount() (string, int64) {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	name := collection.Name()
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

//按选项查询集合
// Skip 跳过
// Limit 读取数量
// sort 1 ，-1 . 1 为升序 ， -1 为降序

func (m *mgo) CollectionDocuments(Skip, Limit int64, sort int, key string, value interface{}) *mongo.Cursor {
	client := db.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	SORT := bson.D{{"_id", sort}}
	filter := bson.D{{key, value}}
	findOptions := options.Find().SetSort(SORT).SetLimit(Limit).SetSkip(Skip)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}

``
