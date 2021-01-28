package client

import (
	"root/pkg/expect"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/internal/message"
	"strconv"
)

func (this *Client) CMD_Login(param []string) {
	expect.True(len(param) == 1)

	uuid := param[0]
	pb := &message.LoginReq{
		Passport:     uuid,
		PassportType: "uuid",
		OS:           "xxxx",
		Version:      "1.x.x",
		PackageCode:  "golang",
	}
	this.uuid = uuid
	msg := network.NewPbMessage(pb, message.MSG_LOGIN_REQ.Int32())
	this.conn.NetSession().SendMsg(msg.Buffer())
	log.Info("login send")
}

func (this *Client) CMD_Enter(param []string) {
	expect.True(len(param) == 1)

	pb := &message.EnterGameReq{
		RID:          this.RID,
		UID:          this.UID,
		SID:          this.SID,
		SecurityCode: this.SecurityCode,
		Passport:     this.uuid,
		NewPlayer:    this.newPlayer,
		LoginTime:    this.timestamp,
		PassportType: "uuid",
		OS:           "xxxx",
		PackageCode:  "golang",
	}
	msg := network.NewPbMessage(pb, message.MSG_ENTER_GAME_REQ.Int32())
	this.conn.NetSession().SendMsg(msg.Buffer())
	log.Info("enterGame send")
}

// slot押注
func (this *Client) CMD_SlotBet(param []string) {
	expect.True(len(param) == 1)

	betVal, _ := strconv.Atoi(param[0])
	betMsg := network.NewPbMessage(&message.SlotBetReq{BetLv: int64(betVal)}, message.MSG_SLOT_BET_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

// slot选择玩法
func (this *Client) CMD_SlotSelect(param []string) {
	expect.True(len(param) == 1)

	betVal, _ := strconv.Atoi(param[0])
	betMsg := network.NewPbMessage(&message.SlotSelectReq{Select: int64(betVal)}, message.MSG_SLOT_SELECT_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

// slot押注
func (this *Client) CMD_UseItem(param []string) {
	expect.True(len(param) == 2)

	itemId, _ := strconv.Atoi(param[0])
	itemCount, _ := strconv.Atoi(param[1])
	betMsg := network.NewPbMessage(&message.UseItemReq{ItemId: int32(itemId), Count: int64(itemCount)}, message.MSG_USE_ITEM_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

// gm addItem 10002 100000
func (this *Client) CMD_GM(param []string) {
	cmd := ""
	for _, v := range param {
		cmd += v
		cmd += " "
	}
	betMsg := network.NewPbMessage(&message.GmCommandReq{CMD: cmd}, message.MSG_GM_COMMAND_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

//eg :blu 30001 3

func (this *Client) CMD_BuildingLevelUp(param []string) {
	expect.True(len(param) == 2)
	buildingId, _ := strconv.Atoi(param[0])
	upToLevel, _ := strconv.Atoi(param[1])
	betMsg := network.NewPbMessage(&message.BuildingLevelUpReq{Id: int32(buildingId), UptoLevel: int32(upToLevel)}, message.MSG_BUILDING_LEVEL_UP_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

//eg :bsu 30001 3
func (this *Client) CMD_BuildingStarUp(param []string) {
	expect.True(len(param) == 2)

	buildingId, _ := strconv.Atoi(param[0])
	upToStar, _ := strconv.Atoi(param[1])
	betMsg := network.NewPbMessage(&message.BuildingStarUpReq{Id: int32(buildingId), UptoStar: int32(upToStar)}, message.MSG_BUILDING_STAR_UP_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

func (this *Client) CMD_Place(param []string) {
	expect.True(len(param) == 1)

	id, _ := strconv.Atoi(param[0])
	islandMap := map[int32]int32{}
	islandMap[1] = 100
	islandMap[2] = 200
	islandMap[3] = 300
	islandMap[4] = 400
	islandMap[5] = 500
	betMsg := network.NewPbMessage(&message.PlaceBuildingReq{Land: &message.Island{
		Id:          int32(id),
		BuildingMap: islandMap,
	}}, message.MSG_PLACE_BUILDING_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}
func (this *Client) CMD_CheckPointFinish(param []string) {
	expect.True(len(param) == 1)

	glod, _ := strconv.Atoi(param[0])

	betMsg := network.NewPbMessage(&message.CheckpointFinishReq{
		Gold:         int64(glod),
		IslandId:     123,
		CheckpointId: 1,
		Success:      1,
	}, message.MSG_CHECKPOINT_FINISH_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}

func (this *Client) CMD_BuildingPropUp(param []string) {
	expect.True(len(param) == 3)
	bId, _ := strconv.Atoi(param[0])
	pId, _ := strconv.Atoi(param[1])
	pLv, _ := strconv.Atoi(param[2])

	betMsg := network.NewPbMessage(&message.BuildingPropUpReq{
		BuildingId: int32(bId),
		PropId:     int32(pId),
		UptoLevel:  int32(pLv),
	}, message.MSG_BUILDING_PROP_UP_REQ.Int32())
	this.conn.NetSession().SendMsg(betMsg.Buffer())
}
