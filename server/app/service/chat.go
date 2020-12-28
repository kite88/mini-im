package service

import (
	"encoding/json"
	"log"
	"miniIM/server/app/dao"
	"miniIM/server/app/model"
	"miniIM/server/conf"
	"miniIM/server/db"
	"miniIM/server/util/EmojiUtil"
	"strconv"
	"time"
)

func FriendList(uid int64) []model.UserModel {
	var list []model.FriendModel
	dao.FriendDao{}.D().Where("user_id = ? or friend_id = ?", uid, uid).Find(&list)
	var uids []int64
	for _, item := range list {
		if isUidInUidList(item.UserId, uids) == false {
			uids = append(uids, item.UserId)
		}
		if isUidInUidList(item.FriendId, uids) == false {
			uids = append(uids, item.FriendId)
		}
	}
	var users []model.UserModel
	dao.UserDao{}.D().Where("id in(?)", uids).Where("id <> ?", uid).Find(&users)
	return users
}

func isUidInUidList(uid int64, uids []int64) bool {
	for _, value := range uids {
		if value == uid {
			return true
		}
	}
	return false
}

func ChatList(uid int64, fuid int64) []model.MessageModel {
	var users []model.UserModel
	var userMap = make(map[string]model.UserModel)
	dao.UserDao{}.D().Where("id = ? or id = ?", uid, fuid).Find(&users)
	for _, value := range users {
		userMap[strconv.FormatInt(value.ID, 10)] = value
	}
	var list []model.MessageModel
	dao.MessageDao{}.D().
		Where("(user_id = ? AND friend_id = ?) or (friend_id = ? AND user_id = ?)", uid, fuid, uid, fuid).
		Order("create_time desc").
		Find(&list)
	for key, value := range list {
		if value.UserId == uid {
			list[key].Type = 1
		}
		if value.FriendId == uid {
			list[key].Type = 2
		}
		list[key].User = userMap[strconv.FormatInt(value.UserId, 10)]
		list[key].Content = EmojiUtil.UnicodeEmojiDecode(value.Content)
	}
	return list
}

type Messages struct {
	Time        int64           `json:"time"`
	Message     string          `json:"message"`
	ContentType int64           `json:"content_type"`
	Uid         int64           `json:"uid"`
	User        model.UserModel `json:"user"`
}

func MessageList(uid int64) []Messages {
	var msgList []model.MessageModel
	dao.MessageDao{}.D().Where("user_id = ?", uid).
		Or("friend_id = ?", uid).
		Order("create_time desc").
		Find(&msgList)
	var uids []int64
	var msgIds []int64
	for _, item := range msgList {
		if isUidInUidList(item.UserId, uids) == false {
			uids = append(uids, item.UserId)
			msgIds = append(msgIds, item.ID)
		}
		if isUidInUidList(item.FriendId, uids) == false {
			uids = append(uids, item.FriendId)
			msgIds = append(msgIds, item.ID)
		}
	}
	var mList []model.MessageModel
	if len(msgIds) > 0 {
		dao.MessageDao{}.D().Where("id IN(?)", msgIds).Find(&mList)
	}
	var messages []Messages
	var uids2 []int64
	for _, item := range mList {
		var tempUid int64
		if item.UserId == uid {
			tempUid = item.FriendId
		}
		if item.FriendId == uid {
			tempUid = item.UserId
		}
		messages = append(messages, Messages{
			Time:        item.CreateTime,
			Message:     item.Content,
			Uid:         tempUid,
			ContentType: item.ContentType,
		})
		uids2 = append(uids2, tempUid)
	}
	var users []model.UserModel
	var userMap = make(map[string]model.UserModel)
	dao.UserDao{}.D().Where("id IN(?)", uids2).Find(&users)
	for _, value := range users {
		userMap[strconv.FormatInt(value.ID, 10)] = value
	}
	for k, item := range messages {
		messages[k].User = userMap[strconv.FormatInt(item.Uid, 10)]
		messages[k].Message = EmojiUtil.UnicodeEmojiDecode(item.Message)
	}
	return sortMessageList(messages)
}

func sortMessageList(list []Messages) []Messages {
	for i := 0; i < len(list)-1; i++ {
		for j := 0; j < len(list)-1-i; j++ {
			if list[j].Time < list[j+1].Time {
				temp := list[j]
				list[j] = list[j+1]
				list[j+1] = temp
			}
		}
	}
	return list
}

func SaveToList(m model.MessageModel) {
	byte, _ := json.Marshal(m)
	err := db.Redis().LPush(conf.ImMsgQueueKey, string(byte)).Err()
	if err != nil {
		log.Println("消息存入队列失败")
	}
}

func SaveToErrorList(m model.MessageModel) {
	byte, _ := json.Marshal(m)
	err := db.Redis().LPush(conf.ImMsgQueueKey+"_err", string(byte)).Err()
	if err != nil {
		log.Println("消息存入队列2失败")
	}
}

func ReadInList() (bool, model.MessageModel) {
	var m model.MessageModel
	result, err := db.Redis().BRPop(10*time.Second, conf.ImMsgQueueKey).Result()
	if err != nil {
		log.Println("消息读取失败")
		return false, m
	}
	if len(result) == 0 {
		log.Println("队列里没有东西啊")
		return false, m
	}
	json.Unmarshal([]byte(result[1]), &m)
	return true, m
}

func SavePersistentData(data model.MessageModel) bool {
	r := dao.MessageDao{}.D().Create(&model.MessageModel{
		ContentType: data.ContentType,
		Content:     EmojiUtil.UnicodeEmojiCode(data.Content),
		CreateTime:  data.CreateTime,
		UserId:      data.UserId,
		FriendId:    data.FriendId,
	})
	if r.Error != nil {
		log.Println("数据永久保存失败")
		return false
	}
	if r.RowsAffected == 0 {
		log.Println("数据永久保存失败")
		return false
	}
	return true
}

//保存离线消息
func SaveOfflineMsgList(uid string, value string) {
	if IsUserLine(uid) == true {
		return
	}
	err := db.Redis().LPush(conf.ImOfflineMsg+uid, value).Err()
	if err != nil {
		log.Println("离线消息存入队列失败")
	}
}

//发送离线消息
func ReadOfflineMsg(uid string) []string {
	var list []string
	resultLen, e := db.Redis().LLen(conf.ImOfflineMsg + uid).Result()
	if e == nil {
		var i int64
		for i = 0; i < resultLen; i++ {
			r, _ := db.Redis().RPop(conf.ImOfflineMsg + uid).Result()
			list = append(list, r)
		}
	}
	return list
}
