package logger

import (
	"github.com/go-logr/logr"
	"k8s.io/klog/v2/klogr"
)

func New() logr.Logger {
	return klogr.New().WithName("MyName")
}
