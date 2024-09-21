package notificationcommonlang

import (
	"strings"

	"github.com/apache/incubator-answer/plugin"
)

const LanguageZhCN Lang = "zh_CN"

type ZhCN struct {
}

func newZhCN() *ZhCN {
	return &ZhCN{}
}

func (z *ZhCN) GetLangType() Lang {
	return LanguageZhCN
}

func (z *ZhCN) TranslateGetCardInfo() string {
	return "查看详情"
}

func (z *ZhCN) TranslateDescription(mixinNotificationMsg *plugin.NotificationMessage) string {
	descriptionTpl := z.getDescriptionTemplate(mixinNotificationMsg.Type)
	description := strings.Replace(descriptionTpl, "%s", mixinNotificationMsg.TriggerUserDisplayName, 1)
	return description
}

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

const (
	zhCnTplUpdateQuestionDescription         = "%s 更新了问题"
	zhCnTplAnswerTheQuestionDescription      = "%s 回答了问题"
	zhCnTplUpVotedTheQuestionDescription     = "%s 赞了问题"
	zhCnTplDownVotedTheQuestionDescription   = "%s 踩了问题"
	zhCnTplUpdateAnswerDescription           = "%s 更新了回答"
	zhCnTplAcceptAnswerDescription           = "%s 采纳了回答"
	zhCnTplCommentQuestionDescription        = "%s 评论了问题"
	zhCnTplCommentAnswerDescription          = "%s 评论了回答"
	zhCnTplUpVotedTheCommentDescription      = "%s 赞了评论"
	zhCnTplReplyToYouDescription             = "%s 回复了你"
	zhCnTplMentionYouDescription             = "%s 提到了你"
	zhCnTplInvitedYouToAnswerDescription     = "%s 邀请了你回答问题"
	zhCnTplYourQuestionIsClosedDescription   = "你的问题已关闭"
	zhCnTplYourQuestionWasDeletedDescription = "你的问题已删除"
	zhCnTplYourAnswerWasDeletedDescription   = "你的回答已删除"
	zhCnTplYourCommentWasDeletedDescription  = "你的评论已删除"
	zhCnTplNewQuestionDescription            = "%s 提出了问题"
)
