package privacy

import "github.com/big0t/constant-chain/common"

type PrivacyLogger struct {
	Log common.Logger
}

func (logger *PrivacyLogger) Init(inst common.Logger) {
	logger.Log = inst
}

// Global instant to use
var Logger = PrivacyLogger{}
