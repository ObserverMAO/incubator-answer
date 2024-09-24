package notification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"regexp"

	"github.com/apache/incubator-answer/internal/schema"
	"github.com/apache/incubator-answer/internal/service/mixinbot"
	mixinbotlang "github.com/apache/incubator-answer/internal/service/mixinbot/lang"
	"github.com/apache/incubator-answer/pkg/display"
	"github.com/fox-one/mixin-sdk-go/v2"
)

// 生成新问题的message
func (es *ExternalNotificationService) sendMixinNewQuestionMessage(ctx context.Context, description, receiverExternalID string, msg *schema.ExternalNotificationMsg) {
	rawData := msg.NewQuestionTemplateRawData
	siteInfo, err := es.siteInfoService.GetSiteGeneral(ctx)
	if err != nil {
		return
	}
	seoInfo, err := es.siteInfoService.GetSiteSeo(ctx)
	if err != nil {
		return
	}
	questionUrl := display.QuestionURL(
		seoInfo.Permalink, siteInfo.SiteUrl, rawData.QuestionID, rawData.QuestionTitle)

	notificationMsgTitle, notificationMsgContent := mixinbot.PrefixTitle+rawData.QuestionTitle, rawData.QuestionParsedText

	titleRunes := []rune(notificationMsgTitle)
	if len(titleRunes) > mixinbot.MaxCardTitleLength {
		notificationMsgTitle = string(titleRunes[:mixinbot.MaxCardTitleLength-len([]rune(mixinbot.Ellipsis))]) + mixinbot.Ellipsis
	}

	re := regexp.MustCompile(`<.*?>`)
	notificationMsgContent = re.ReplaceAllString(notificationMsgContent, "")

	contentRunes := []rune(notificationMsgContent)
	if len(contentRunes) > mixinbot.MaxCardContentLength {
		notificationMsgContent = string(contentRunes[:mixinbot.MaxCardContentLength-len([]rune(mixinbot.Ellipsis))]) + mixinbot.Ellipsis
	}

	lastContent := "\n\n\n" + description
	totalContentRunes := []rune(notificationMsgContent + lastContent)
	if len(totalContentRunes) <= mixinbot.MaxCardContentLength {
		notificationMsgContent += lastContent
	} else {
		availableLength := mixinbot.MaxCardContentLength - len([]rune(lastContent)) - len([]rune(mixinbot.Ellipsis))
		notificationMsgContent = string(contentRunes[:availableLength]) + mixinbot.Ellipsis + lastContent
	}

	card := &mixin.AppCardMessage{
		AppID:       es.mixinBotService.Config.ClientID,
		Title:       notificationMsgTitle,
		Description: notificationMsgContent,
		Shareable:   true,
	}
	es.fillCardAction(card, msg, questionUrl)

	cardBytes, err := json.Marshal(card)
	if err != nil {
		return
	}

	cardBase64code := base64.StdEncoding.EncodeToString(cardBytes)
	messageRequest := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(es.mixinBotService.Config.ClientID, receiverExternalID),
		RecipientID:    receiverExternalID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           cardBase64code,
	}

	es.mixinBotService.SendMessage(ctx, messageRequest)
}

func (es *ExternalNotificationService) fillCardAction(card *mixin.AppCardMessage, msg *schema.ExternalNotificationMsg, questionUrl string) {
	btnMsg := mixin.AppButtonMessage{
		Label:  es.langPicker.Pick(mixinbotlang.GetLanguage(msg.ReceiverLang)).TranslateGetCardInfo(),
		Action: questionUrl,
		Color:  mixinbot.RandomCardColor(),
	}
	btnMsg.Action = questionUrl
	card.Actions = []mixin.AppButtonMessage{btnMsg}
}
