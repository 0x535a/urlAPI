package error

// 1xx: 表示请求有问题

const (
	ServerOtherError = iota + 100
	ServerInvalidParams
	ServerNotificationFailed
	ServerPermissionDenied
)
