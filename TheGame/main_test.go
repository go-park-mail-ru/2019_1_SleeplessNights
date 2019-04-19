package TheGame

import (
	"github.com/xlab/closer"
	"testing"
)

func TestMain(m *testing.M) {
	defer closer.Close()
	m.Run()
}
