package biz

import (
	"at-kratos/internal/pkg/util"
	"errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewSignupLogUseCase,
)

const PerSignupLog = 2000

var WayOfSignupLog util.ArrayString = []string{"pc", "wap"}
var errClosed = errors.New("collector closed")
var errTimeout = errors.New("collector timeout")
