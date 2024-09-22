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

	return mixinBot, nil
}
