package mixinbotlang

import (
	"strings"

	"github.com/apache/incubator-answer/plugin"
)

const LanguageEnUS Lang = "en_US"

type EnUS struct {
}

func newEnUS() *EnUS {
	return &EnUS{}
}

func (e *EnUS) GetLangType() Lang {
	return LanguageEnUS
}
func (e *EnUS) TranslateGetCardInfo() string {
	return "Details"
}

func (e *EnUS) TranslateDescription(mixinNotificationMsg *plugin.NotificationMessage) string {
	descriptionTpl := e.getDescriptionTemplate(mixinNotificationMsg.Type)
	description := strings.Replace(descriptionTpl, "%s", mixinNotificationMsg.TriggerUserDisplayName, 1)
	return description
}

func (e *EnUS) getDescriptionTemplate(msgType plugin.NotificationType) string {
	switch msgType {
	case plugin.NotificationUpdateQuestion:
		return enUsTplUpdateQuestionDescription
	case plugin.NotificationAnswerTheQuestion:
		return enUsTplAnswerTheQuestionDescription
	case plugin.NotificationUpVotedTheQuestion:
		return enUsTplUpVotedTheQuestionDescription
	case plugin.NotificationDownVotedTheQuestion:
		return enUsTplDownVotedTheQuestionDescription
	case plugin.NotificationUpdateAnswer:
		return enUsTplUpdateAnswerDescription
	case plugin.NotificationAcceptAnswer:
		return enUsTplAcceptAnswerDescription
	case plugin.NotificationUpVotedTheAnswer:
		return enUsTplUpVotedTheAnswerDescription
	case plugin.NotificationDownVotedTheAnswer:
		return enUsTplDownVotedTheAnswerDescription
	case plugin.NotificationCommentQuestion:
		return enUsTplCommentQuestionDescription
	case plugin.NotificationCommentAnswer:
		return enUsTplCommentAnswerDescription
	case plugin.NotificationUpVotedTheComment:
		return enUsTplUpVotedTheCommentDescription
	case plugin.NotificationReplyToYou:
		return enUsTplReplyToYouDescription
	case plugin.NotificationMentionYou:
		return enUsTplMentionYouDescription
	case plugin.NotificationYourQuestionIsClosed:
		return enUsTplYourQuestionIsClosedDescription
	case plugin.NotificationYourQuestionWasDeleted:
		return enUsTplYourAnswerWasDeletedDescription
	case plugin.NotificationYourCommentWasDeleted:
		return enUsTplYourCommentWasDeletedDescription
	case plugin.NotificationInvitedYouToAnswer:
		return enUsTplInvitedYouToAnswerDescription
	case plugin.NotificationNewQuestion:
		return enUsTplNewQuestionDescription
	case plugin.NotificationNewQuestionFollowedTag:
		return enUsTplNewQuestionFollowedTagDescription
	default:
		return ""
	}
}

const (
	enUsTplUpdateQuestionDescription         = "%s updated the topic"
	enUsTplAnswerTheQuestionDescription      = "%s replied the topic"
	enUsTplUpVotedTheQuestionDescription     = "%s upvoted the topic"
	enUsTplDownVotedTheQuestionDescription   = "%s downvoted the topic"
	enUsTplUpdateAnswerDescription           = "%s updated the reply"
	enUsTplAcceptAnswerDescription           = "%s accepted the reply"
	enUsTplUpVotedTheAnswerDescription       = "%s upvoted the reply"
	enUsTplDownVotedTheAnswerDescription     = "%s downvoted the reply"
	enUsTplCommentQuestionDescription        = "%s commented the topic"
	enUsTplCommentAnswerDescription          = "%s commented the reply"
	enUsTplUpVotedTheCommentDescription      = "%s upvoted the comment"
	enUsTplReplyToYouDescription             = "%s replied to you"
	enUsTplMentionYouDescription             = "%s mentioned you"
	enUsTplInvitedYouToAnswerDescription     = "%s invited you to reply the topic"
	enUsTplYourQuestionIsClosedDescription   = "Your topic is closed"
	enUsTplYourQuestionWasDeletedDescription = "Your topic was deleted"
	enUsTplYourAnswerWasDeletedDescription   = "Your reply was deleted"
	enUsTplYourCommentWasDeletedDescription  = "Your comment was deleted"
	enUsTplNewQuestionDescription            = "%s created a topic"
	enUsTplNewQuestionFollowedTagDescription = "%s followed the topic"
)
