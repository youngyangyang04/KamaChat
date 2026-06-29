package dao

import (
	"gochat/internal/dao"
	"gochat/internal/model"
	"gochat/pkg/util/random"
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	if err := dao.Init(); err != nil {
		t.Skipf("skip integration test without database: %v", err)
	}

	userInfo := &model.UserInfo{
		Uuid:      "U" + strconv.Itoa(random.GetRandomInt(11)),
		Nickname:  "apylee",
		Telephone: "18032353211",
		Email:     "1212312312@qq.com",
		Password:  "123456",
		CreatedAt: time.Now(),
		IsAdmin:   1,
	}
	if err := dao.GormDB.Create(userInfo).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
}
