package test

import (
	"search_engine/utils"
	"testing"
)

func TestLog(t *testing.T) {
	utils.Logger.Infoln("this is an info log")
}
