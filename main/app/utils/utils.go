package utils

import (
	"github.com/itang/gotang"
	"github.com/nu7hatch/gouuid"
)

func Uuid() string {
	u4, err := uuid.NewV4()
	gotang.AssertNoError(err)
	return u4.String()
}
