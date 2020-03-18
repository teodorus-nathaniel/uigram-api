package helper

import "go.mongodb.org/mongo-driver/bson"

func ConvertStructToBSON(data interface{}) (bson.M, error) {
	bytes, _ := bson.Marshal(data)
	bsonData := bson.M{}
	err := bson.Unmarshal(bytes, &bsonData)
	
	return bsonData, err
}
