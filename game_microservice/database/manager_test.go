package database_test

import (
	"github.com/xlab/closer"
	"testing"
)

func TestMain(m *testing.M){
	defer closer.Close()
	m.Run()
}
