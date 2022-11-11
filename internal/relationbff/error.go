package relationbff

import "errors"

var (
	ErrSourceUndefined = errors.New("source undefined")
	ErrFollowLimit     = errors.New("follow limit")
	ErrBlacked         = errors.New("blacked")
)
