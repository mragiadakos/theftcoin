package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/mragiadakos/theftcoin/server/confs"
	"github.com/mragiadakos/theftcoin/server/ctrls"
	absrv "github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"
	tmlog "github.com/tendermint/tmlibs/log"
)

func main() {
	logger := tmlog.NewTMLogger(kitlog.NewSyncWriter(os.Stdout))
	flagAbci := "socket"
	ipfsDaemon := flag.String("ipfs", "127.0.0.1:5001", "the URL for the IPFS's daemon")
	node := flag.String("node", "tcp://0.0.0.0:46658", "the TCP URL for the ABCI daemon")
	ipfsInflatorsHash := flag.String("inflators", "", "the IPFS hash with the JSON list of public keys for inflators")
	ipfsWatchersHash := flag.String("watchers", "", "the IPFS hash with the JSON list of public keys for watchers")
	ipfsTaxHash := flag.String("tax", "", "the IPFS hash with the JSON for the tax")
	waitSec := flag.Int("wait", 5, "the seconds for an acceptable query")
	createDemoKeys := flag.Bool("create-demo-keys", false, "Create the first demo keys.")
	flag.Parse()

	if *createDemoKeys {
		CreateInflatorsPublicKey()
		CreateTaxPublicKey()
		CreateWatchersPublicKey()
		return
	}
	if len(*ipfsInflatorsHash) == 0 {
		fmt.Println("Error ", errors.New("The IPFS hash for inflators is missing"))
		return
	}

	confs.Conf.IpfsInflators = *ipfsInflatorsHash

	err := confs.Conf.SubmitInflators()
	if err != nil {
		fmt.Println("Error ", err.Error())
		return
	}

	if len(*ipfsWatchersHash) == 0 {
		fmt.Println("Error ", errors.New("The IPFS hash for watchers is missing"))
		return
	}

	confs.Conf.IpfsWatchers = *ipfsWatchersHash

	err = confs.Conf.SubmitWatchers()
	if err != nil {
		fmt.Println("Error ", err.Error())
		return
	}

	if len(*ipfsTaxHash) == 0 {
		fmt.Println("Error ", errors.New("The IPFS hash for tax is missing"))
		return
	}

	confs.Conf.IpfsTax = *ipfsTaxHash

	err = confs.Conf.SubmitTax()
	if err != nil {
		fmt.Println("Error ", err.Error())
		return
	}

	confs.Conf.AbciDaemon = *node
	confs.Conf.IpfsConnection = *ipfsDaemon
	confs.Conf.WaitingRequestTime = *waitSec

	app := ctrls.NewTCApplication()
	srv, err := absrv.NewServer(confs.Conf.AbciDaemon, flagAbci, app)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		fmt.Println("Error ", err)
		return
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}
