package router

const groupPrefix = "Group-"
const roomPrefix = "Room-"

func GetGroupRouteById(str string) string {
	return groupPrefix + str
}

func GetRouteByToid(str interface{}) string {
	return "default"
}

func GetRoomRouteById(str string) string {
	return roomPrefix + str
}
