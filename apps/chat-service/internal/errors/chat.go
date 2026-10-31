package customErrors

import "errors"

var DirectChatExist = errors.New("direct chat already exists")

var IncorrectChatType = errors.New("incorrect chat type")

var CannotJoinDirectChat = errors.New("cannot join to direct chat")

var ChatNotFound = errors.New("chat not found")
