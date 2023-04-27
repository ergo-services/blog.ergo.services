package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"

	"blockinfo/apps/blockinfoapp"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"

	log "github.com/sirupsen/logrus"
)

var (
	OptionNodeName   string
	OptionNodeCookie string
)

func init() {
	// generate random value for cookie
	buff := make([]byte, 12)
	rand.Read(buff)
	randomCookie := hex.EncodeToString(buff)

	flag.StringVar(&OptionNodeName, "name", "blockinfo@localhost", "node name")
	flag.StringVar(&OptionNodeCookie, "cookie", randomCookie, "a secret cookie for interaction within the cluster")

}

func main() {
	var options node.Options

	flag.Parse()

	// Create applications that must be started
	apps := []gen.ApplicationBehavior{
		blockinfoapp.CreateBlockInfoApp(),
	}
	options.Applications = apps

	// Starting node
	blockinfoNode, err := ergo.StartNode(OptionNodeName, OptionNodeCookie, options)
	if err != nil {
		panic(err)
	}
	log.Infof("Node %q is started\n", blockinfoNode.Name())

	// starting process Storage
	_, err = blockinfoNode.Spawn("storage", gen.ProcessOptions{}, createStorage())
	if err != nil {
		panic(err)
	}

	// starting process Web
	_, err = blockinfoNode.Spawn("web", gen.ProcessOptions{}, createWeb())
	if err != nil {
		panic(err)
	}

	blockinfoNode.Wait()
}
