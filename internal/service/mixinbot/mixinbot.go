package mixinbot

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
	"google.golang.org/appengine/log"
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

func (m *MixinBotService) SendMessage(ctx context.Context, message *mixin.MessageRequest) error {
	sendMessageToUser := func(retry int) error {
		err := m.Client.SendMessage(ctx, message)
		if err != nil {
			// try create conversation
			_, err = m.Client.CreateContactConversation(ctx, message.RecipientID)
			if err != nil {
				log.Errorf(ctx, "create contact conversation failed: %v, data: %v, retry: %d", err, message, retry)
				return err
			}
			err = m.Client.SendMessage(ctx, message)
			if err != nil {
				log.Errorf(ctx, "send mixin notification failed: %v, data: %v, retry: %d", err, message, retry)
				return err
			}
		}
		return nil
	}

	var err error
	for retry := 0; retry < maxRetries; retry++ {
		err = sendMessageToUser(retry)
		if err == nil {
			break
		}
		delay := baseDelay * time.Duration(1<<retry)
		time.Sleep(delay)
	}

	return err
}

var cardColorList = []string{
	"#7983C2", "#8F7AC5", "#C5595A", "#C97B46", "#76A048", "#3D98D0",
	"#5979F0", "#8A64D0", "#B76753", "#AA8A46", "#9CAD23", "#6BC0CE",
	"#6C89D3", "#AA66C3", "#C8697D", "#C49B4B", "#5FB05F", "#52A98B",
	"#75A2CB", "#A75C96", "#9B6D77", "#A49373", "#6AB48F", "#93B289",
}

func RandomCardColor() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return cardColorList[r.Intn(len(cardColorList))]
}
