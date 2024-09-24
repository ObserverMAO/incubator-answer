package mixinbotlang

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
	case plugin.NotificationUpVotedTheAnswer:
		return zhCnTplUpVotedTheAnswerDescription
	case plugin.NotificationDownVotedTheAnswer:
		return zhCnTplDownVotedTheAnswerDescription
	case plugin.NotificationCommentQuestion:
		return zhCnTplCommentQuestionDescription
	case plugin.NotificationCommentAnswer:
		return zhCnTplCommentAnswerDescription
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
	case plugin.NotificationInvitedYouToAnswer:
		return zhCnTplInvitedYouToAnswerDescription
	case plugin.NotificationNewQuestion:
		return zhCnTplNewQuestionDescription
	case plugin.NotificationNewQuestionFollowedTag:
		return zhCnTplNewQuestionFollowedTagDescription
	default:
		return ""
	}
}

const (
	zhCnTplUpdateQuestionDescription         = "%s 更新了话题"
	zhCnTplAnswerTheQuestionDescription      = "%s 回答了话题"
	zhCnTplUpVotedTheQuestionDescription     = "%s 赞了话题"
	zhCnTplDownVotedTheQuestionDescription   = "%s 踩了话题"
	zhCnTplUpdateAnswerDescription           = "%s 更新了回复"
	zhCnTplUpVotedTheAnswerDescription       = "%s 赞了回复"
	zhCnTplDownVotedTheAnswerDescription     = "%s 踩了回复"
	zhCnTplAcceptAnswerDescription           = "%s 采纳了回复"
	zhCnTplCommentQuestionDescription        = "%s 评论了话题"
	zhCnTplCommentAnswerDescription          = "%s 评论了回复"
	zhCnTplUpVotedTheCommentDescription      = "%s 赞了评论"
	zhCnTplReplyToYouDescription             = "%s 回复了你"
	zhCnTplMentionYouDescription             = "%s 提到了你"
	zhCnTplInvitedYouToAnswerDescription     = "%s 邀请了你回复话题"
	zhCnTplYourQuestionIsClosedDescription   = "你的话题已关闭"
	zhCnTplYourQuestionWasDeletedDescription = "你的话题已删除"
	zhCnTplYourAnswerWasDeletedDescription   = "你的回复已删除"
	zhCnTplYourCommentWasDeletedDescription  = "你的评论已删除"
	zhCnTplNewQuestionDescription            = "%s 发起了话题"
	zhCnTplNewQuestionFollowedTagDescription = "%s 关注了话题"
)
