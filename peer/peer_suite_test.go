package peer

import (
	"os"
	"path"
	"testing"

	"github.com/ellcrys/druid/configdir"
	"github.com/ellcrys/druid/util"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/ellcrys/druid/util/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var log = logger.NewLogrusNoOp()
var cfg *configdir.Config

func setTestCfg() error {
	var err error
	dir, _ := homedir.Dir()
	cfgDir := path.Join(dir, ".ellcrys_test")
	os.MkdirAll(cfgDir, 0700)
	cfg, err = util.LoadCfg(cfgDir)
	cfg.Peer.Dev = true
	cfg.Peer.MaxAddrsExpected = 5
	cfg.Peer.Test = true
	return err
}

func removeTestCfgDir() error {
	dir, _ := homedir.Dir()
	err := os.RemoveAll(path.Join(dir, ".ellcrys_test"))
	return err
}
func TestPeer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Peer Suite")
}
