package main

import (
	"flag"
	"log"
	"time"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/entity"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	dbPath = flag.String("db", "test.db", "db path")
)

func main() {

	flag.Parse()
	engine, err := xorm.NewEngine("sqlite", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()

	users := make([]entity.User, 0)
	err = engine.Find(&users)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		var configs []entity.UserNotificationConfig
		err = engine.Where("user_id = ?", user.ID).Find(&configs)
		if err != nil {
			log.Fatal(err)
		}
		has := false
		for _, config := range configs {
			if config.Source == string(constant.AllNewQuestionForFollowingTagsSource) {
				has = true
			}
			if has {
				config.Channels = `[{"key":"email","enable":true}]`
				_, err = engine.ID(config.ID).Update(&config)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		if !has {
			configs = append(configs, entity.UserNotificationConfig{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Enabled:   true,
				UserID:    user.ID,
				Source:    string(constant.AllNewQuestionForFollowingTagsSource),
				Channels:  `[{"key":"email","enable":true}]`,
			})
			_, err = engine.Insert(&configs)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
