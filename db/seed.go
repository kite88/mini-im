package db

import (
	"errors"
	"log"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mini-im/conf"
	"mini-im/model"
)

// AvatarURL 生成内置头像地址（由 /avatar 接口动态渲染 SVG）
func AvatarURL(username, nickname string) string {
	q := url.Values{}
	q.Set("name", username)
	if nickname != "" {
		q.Set("label", nickname)
	}
	return "/avatar?" + q.Encode()
}

// Seed 初始化演示数据：
//  1. 按配置写入账号（已存在则只补齐昵称/头像，不改密码）
//  2. 将所有账号两两建立好友关系，保证开箱即可聊天
//
// 关闭方式：config.yaml -> seed.enabled = false
func Seed() error {
	cfg := conf.C.Seed
	if !cfg.Enabled || len(cfg.Accounts) == 0 {
		return nil
	}
	if DB == nil {
		return errors.New("数据库尚未初始化")
	}

	password := cfg.Password
	if password == "" {
		password = "123123"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	ids := make([]int64, 0, len(cfg.Accounts))

	for _, item := range cfg.Accounts {
		if item.Username == "" {
			continue
		}

		var user model.User
		err := DB.Where("username = ?", item.Username).First(&user).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			user = model.User{
				Username:   item.Username,
				Nickname:   item.Nickname,
				Avatar:     AvatarURL(item.Username, item.Nickname),
				Password:   string(hash),
				CreateTime: now,
				UpdateTime: now,
			}
			if user.Nickname == "" {
				user.Nickname = item.Username
			}
			if err = DB.Create(&user).Error; err != nil {
				return err
			}
			log.Printf("[seed] 创建账号 %s(id=%d)", user.Username, user.ID)
		case err != nil:
			return err
		default:
			avatar := user.Avatar
			if avatar == "" {
				avatar = AvatarURL(item.Username, item.Nickname)
			}
			if item.Nickname != "" && user.Nickname != item.Nickname {
				_ = DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
					"nickname":    item.Nickname,
					"avatar":      avatar,
					"update_time": now,
				})
			}
		}
		ids = append(ids, user.ID)
	}

	return seedFriendRelations(ids, now)
}

// seedFriendRelations 两两建立好友关系（全连通）
func seedFriendRelations(ids []int64, now int64) error {
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			var count int64
			err := DB.Model(&model.Friend{}).
				Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
					ids[i], ids[j], ids[j], ids[i]).
				Count(&count).Error
			if err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			if err = DB.Create(&model.Friend{
				UserID:   ids[i],
				FriendID: ids[j],
				AddTime:  now,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
