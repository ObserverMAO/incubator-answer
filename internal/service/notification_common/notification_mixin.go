// notification mixin
package notificationcommon

import (
	"context"
	"encoding/base64"
	"math/rand"
	"regexp"
	"time"

	notificationcommonlang "github.com/apache/incubator-answer/internal/service/notification_common/lang"
	"github.com/apache/incubator-answer/plugin"
	"github.com/goccy/go-json"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/segmentfault/pacman/log"
)

const (
	maxCardTitleLength   = 28
	maxCardContentLength = 2 << 8
	prefixTitle          = "💬"
	ellipsis             = "..."
)

func (ns *NotificationCommon) sendMixinNotification(notificationMsgDescription string, msg *plugin.NotificationMessage) error {
	notificationMsgTitle, notificationMsgContent := prefixTitle+msg.QuestionTitle, msg.Content

	titleRunes := []rune(notificationMsgTitle)
	if len(titleRunes) > maxCardTitleLength {
		notificationMsgTitle = string(titleRunes[:maxCardTitleLength-len([]rune(ellipsis))]) + ellipsis
	}

	re := regexp.MustCompile(`<.*?>`)
	notificationMsgContent = re.ReplaceAllString(notificationMsgContent, "")

	contentRunes := []rune(notificationMsgContent)
	if len(contentRunes) > maxCardContentLength {
		notificationMsgContent = string(contentRunes[:maxCardContentLength-len([]rune(ellipsis))]) + ellipsis
	}

	lastContent := "\n\n\n" + notificationMsgDescription
	totalContentRunes := []rune(notificationMsgContent + lastContent)
	if len(totalContentRunes) <= maxCardContentLength {
		notificationMsgContent += lastContent
	} else {
		availableLength := maxCardContentLength - len([]rune(lastContent)) - len([]rune(ellipsis))
		notificationMsgContent = string(contentRunes[:availableLength]) + ellipsis + lastContent
	}

	card := &mixin.AppCardMessage{
		AppID:       ns.mixinbotService.Config.ClientID,
		Title:       notificationMsgTitle,
		Description: notificationMsgContent,
		Shareable:   true,
	}
	ns.fillCardAction(card, msg)

	cardBytes, err := json.Marshal(card)
	if err != nil {
		return err
	}

	cardBase64code := base64.StdEncoding.EncodeToString(cardBytes)
	err = ns.mixinbotService.Client.SendMessage(context.Background(), &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(ns.mixinbotService.Config.ClientID, msg.ReceiverExternalID),
		RecipientID:    msg.ReceiverExternalID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           cardBase64code,
	})
	if err != nil {
		log.Errorf("send mixin notification failed: %v", err)
		return err
	}
	return nil
}

func (ns *NotificationCommon) fillCardAction(card *mixin.AppCardMessage, msg *plugin.NotificationMessage) {
	btnMsg := mixin.AppButtonMessage{
		Label:  ns.langPicker.Pick(notificationcommonlang.GetLanguage(msg.ReceiverLang)).TranslateGetCardInfo(),
		Action: msg.QuestionUrl,
		Color:  randomCardColor(),
	}

	switch msg.Type {
	case plugin.NotificationUpdateQuestion, plugin.NotificationInvitedYouToAnswer, plugin.NotificationNewQuestion, plugin.NotificationNewQuestionFollowedTag:
		btnMsg.Action = msg.QuestionUrl
	case plugin.NotificationAnswerTheQuestion, plugin.NotificationUpdateAnswer, plugin.NotificationAcceptAnswer:
		btnMsg.Action = msg.AnswerUrl
	case plugin.NotificationCommentQuestion, plugin.NotificationCommentAnswer, plugin.NotificationReplyToYou, plugin.NotificationMentionYou:
		btnMsg.Action = msg.CommentUrl
	default:
		log.Debugf("this type of notification will be drop, the type is %s", msg.Type)
	}
	card.Actions = []mixin.AppButtonMessage{btnMsg}
}

var cardColorList = []string{
	"#7983C2", "#8F7AC5", "#C5595A", "#C97B46", "#76A048", "#3D98D0",
	"#5979F0", "#8A64D0", "#B76753", "#AA8A46", "#9CAD23", "#6BC0CE",
	"#6C89D3", "#AA66C3", "#C8697D", "#C49B4B", "#5FB05F", "#52A98B",
	"#75A2CB", "#A75C96", "#9B6D77", "#A49373", "#6AB48F", "#93B289",
}

func randomCardColor() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return cardColorList[r.Intn(len(cardColorList))]
}
