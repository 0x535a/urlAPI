package error

// 2xx: 表示op内有问题
const (
	OpOtherError = iota + 200
	OpConfigNotImplemented
	OpWxClientError
	OpWxResponseError
	OpWxTokenExpired
	OpJSONError
	OpURLParseError
	_OpMachineIsNotAvailable
	OpOrderNotAssignedToWorker
	_OpStringToIntFailed
	OpUserWXIDNotFound
	OpMachineNotAvailable
	OpUserAlreadyBinded
	OpMachineNotAssignedToOrder
	OpMachineAlreadyInUse
)
