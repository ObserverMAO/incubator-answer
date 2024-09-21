package mixinbot

import (
	"context"
	"os"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
)

type MixinBotService struct {
	Config   MixinBotConfig
	SpendKey mixinnet.Key

	User   *mixin.User
	Client *mixin.Client
}

func NewMixinBotService(config *MixinBotConfig) (*MixinBotService, error) {
	mixinBot := &MixinBotService{}

	client, err := mixin.NewFromKeystore(&mixin.Keystore{
		ClientID:          config.ClientID,
		AppID:             config.ClientID,
		SessionID:         config.SessionID,
		ServerPublicKey:   config.ServerPublicKey,
		SessionPrivateKey: config.SessionPrivateKey,
	})
	if err != nil {
		return nil, err
	}

	user, err := client.UserMe(context.Background())
	if err != nil {
		return nil, err
	}

	if config.SpendKey == "" {
		config.SpendKey = os.Getenv("SPEND_KEY")
	}

	if config.SpendKey != "" {
		spendKey, err := mixinnet.ParseKeyWithPub(config.SpendKey, user.SpendPublicKey)
		if err != nil {
			return nil, err
		}
		mixinBot.SpendKey = spendKey
	}

	mixinBot.Config = *config
	mixinBot.User = user
	mixinBot.Client = client

	// err = mixinBot.Client.SendMessage(context.Background(), &mixin.MessageRequest{
	// 	ConversationID: mixin.UniqueConversationID(config.ClientID, mixinBot.User.App.CreatorID),
	// 	RecipientID:    mixinBot.User.App.CreatorID,
	// 	MessageID:      mixin.RandomTraceID(),
	// 	Category:       mixin.MessageCategoryPlainText,
	// 	Data:           base64.StdEncoding.EncodeToString([]byte("您已连接机器人")),
	// })
	// if err != nil {
	// 	log.Errorf("send message to mixin bot failed: %+v", err)
	// 	return nil, err
	// }

	return mixinBot, nil
}
