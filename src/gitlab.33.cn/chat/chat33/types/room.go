package types

const (
	RoomType = 1

	RoomLevelNomal   = 1
	RoomLevelManager = 2
	RoomLevelMaster  = 3

	RoomNodisturbingOff = 2
	RoomNodisturbingOn  = 1

	RoomUncommonUse = 1
	RoomCommonUse   = 2

	RoomOnTop    = 1
	RoomNotOnTop = 2

	RoomNotDeleted = 1
	RoomDeleted    = 2

	CanAddFriend      = 1
	CanNotAddFriend   = 2
	ShouldApproval    = 1
	ShouldNotApproval = 2
	CanNotJoinRoom    = 3

	AdminMuted    = 2
	AdminNotMuted = 1

	MasterMuted    = 2
	MasterNotMuted = 1

	RoomUserDeleted    = 2
	RoomUserNotDeleted = 1

	JoinApplyNo   = 2
	JoinApplyOk   = 3
	JoinApplyWait = 1
)
