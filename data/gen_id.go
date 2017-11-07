package data

import "gopkg.in/mgo.v2/bson"

type RoomIDGen struct {
	ServerID   string `bson:"_id"`
	LastRoomID string `bson:"LastRoomID"`
}

func (this *RoomIDGen) Save() bool {
	return Upsert(RoomIDGens, bson.M{"_id": this.ServerID}, this)
}

func (this *RoomIDGen) Get() {
	Get(RoomIDGens, this.ServerID, this)
}

type UserIDGen struct {
	ServerID   string `bson:"_id"`
	LastUserID string `bson:"LastUserID"`
}

func (this *UserIDGen) Save() bool {
	return Upsert(UserIDGens, bson.M{"_id": this.ServerID}, this)
}

func (this *UserIDGen) Get() {
	Get(UserIDGens, this.ServerID, this)
}
