package notificationcommonlang

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

/*
	func (z *ZhCN) getDescriptionTemplate(msgType plugin.NotificationType) string {
		switch msgType {
		case plugin.NotificationUpdateQuestion:
			return zhCnTplUpdateQuestionDescription
		case plugin.NotificationAnswerTheQuestion:
			return zhCnTplAnswerTheQuestionDescription
		case plugin.NotificationUpVotedTheQuestion:
			return zhCnTplUpVotedTheQuestionDescription
		case plugin.NotificationDownVotedTheQuestion:
			return zhCnTplDownVotedTheQuestionDescription
		case plugin.NotificationUpdateAnswer:
			return zhCnTplUpdateAnswerDescription
		case plugin.NotificationAcceptAnswer:
			return zhCnTplAcceptAnswerDescription
		case plugin.NotificationCommentQuestion:
			return zhCnTplCommentQuestionDescription
		case plugin.NotificationUpVotedTheComment:
			return zhCnTplUpVotedTheCommentDescription
		case plugin.NotificationReplyToYou:
			return zhCnTplReplyToYouDescription
		case plugin.NotificationMentionYou:
			return zhCnTplMentionYouDescription
		case plugin.NotificationYourQuestionIsClosed:
			return zhCnTplYourQuestionIsClosedDescription
		case plugin.NotificationYourQuestionWasDeleted:
			return zhCnTplYourQuestionWasDeletedDescription
		case plugin.NotificationYourAnswerWasDeleted:
			return zhCnTplYourAnswerWasDeletedDescription
		case plugin.NotificationYourCommentWasDeleted:
			return zhCnTplYourCommentWasDeletedDescription
		case plugin.NotificationCommentAnswer:
			return zhCnTplCommentAnswerDescription
		case plugin.NotificationInvitedYouToAnswer:
			return zhCnTplInvitedYouToAnswerDescription
		case plugin.NotificationNewQuestion, plugin.NotificationNewQuestionFollowedTag:
			return zhCnTplNewQuestionDescription
		default:
			return ""
		}
	}
*/
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
	enUsTplUpdateQuestionDescription         = "%s updated the question"
	enUsTplAnswerTheQuestionDescription      = "%s answered the question"
	enUsTplUpVotedTheQuestionDescription     = "%s upvoted the question"
	enUsTplDownVotedTheQuestionDescription   = "%s downvoted the question"
	enUsTplUpdateAnswerDescription           = "%s updated the answer"
	enUsTplAcceptAnswerDescription           = "%s accepted the answer"
	enUsTplUpVotedTheAnswerDescription       = "%s upvoted the answer"
	enUsTplDownVotedTheAnswerDescription     = "%s downvoted the answer"
	enUsTplCommentQuestionDescription        = "%s commented the question"
	enUsTplCommentAnswerDescription          = "%s commented the answer"
	enUsTplUpVotedTheCommentDescription      = "%s upvoted the comment"
	enUsTplReplyToYouDescription             = "%s replied to you"
	enUsTplMentionYouDescription             = "%s mentioned you"
	enUsTplInvitedYouToAnswerDescription     = "%s invited you to answer the question"
	enUsTplYourQuestionIsClosedDescription   = "Your question is closed"
	enUsTplYourQuestionWasDeletedDescription = "Your question was deleted"
	enUsTplYourAnswerWasDeletedDescription   = "Your answer was deleted"
	enUsTplYourCommentWasDeletedDescription  = "Your comment was deleted"
	enUsTplNewQuestionDescription            = "%s asked a question"
	enUsTplNewQuestionFollowedTagDescription = "%s asked a question"
)
