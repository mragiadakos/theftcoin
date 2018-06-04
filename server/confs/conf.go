package confs

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/ipfs/go-ipfs-api"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

type configuration struct {
	IpfsConnection     string
	AbciDaemon         string
	WaitingRequestTime int
	inflators          map[string]int
	watchers           map[string]int
	IpfsTax            string
	IpfsInflators      string
	IpfsWatchers       string
	Tax                Tax
	TaxReceiver        crypto.PubKey
}

type Tax struct {
	Percentage   int
	PublicKeyHex string
}

func (t *Tax) Bytes() ([]byte, error) {
	return hex.DecodeString(t.PublicKeyHex)
}

func (t *Tax) SetPublic(b []byte) {
	t.PublicKeyHex = hex.EncodeToString(b)
}

type Inflator struct {
	PublicKeyHex string
}

func (i *Inflator) SetPublic(b []byte) {
	i.PublicKeyHex = hex.EncodeToString(b)
}

func (i *Inflator) Bytes() ([]byte, error) {
	return hex.DecodeString(i.PublicKeyHex)
}

type Watcher struct {
	PublicKeyHex string
}

func (w *Watcher) SetPublic(b []byte) {
	w.PublicKeyHex = hex.EncodeToString(b)
}

func (w *Watcher) Bytes() ([]byte, error) {
	return hex.DecodeString(w.PublicKeyHex)
}

func (c *configuration) InflatorExists(inlf string) bool {
	_, ok := c.inflators[inlf]
	return ok
}

func (c *configuration) WatcherExists(watcher string) bool {
	_, ok := c.watchers[watcher]
	return ok
}

func (c *configuration) SubmitTax() error {
	sh := shell.NewShell(c.IpfsConnection)
	b, err := sh.BlockGet(c.IpfsTax)
	if err != nil {
		return errors.New("The hash for the tax is not correct: " + err.Error())
	}
	cleaned := cleanJsonFromFileBytesOfIpfs(string(b))
	tax := Tax{}
	err = json.Unmarshal([]byte(cleaned), &tax)
	if err != nil {
		return errors.New("The json for the tax is not correct: " + err.Error())
	}

	pubB, err := tax.Bytes()
	if err != nil {
		return errors.New("The json for the tax is not correct: " + err.Error())
	}

	pubk, err := crypto.UnmarshalPublicKey(pubB)
	if err != nil {
		return errors.New("The tax receiver's public key is not correct")
	}
	c.TaxReceiver = pubk
	c.Tax = tax
	return nil
}

func cleanArrayJsonFromFileBytesOfIpfs(str string) string {
	example := []map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &example)
	if err != nil {
		re := regexp.MustCompile("\\[(.*)\\]")
		strs := re.FindAllString(str, -1)
		ret := strs[0]
		return ret[:len(ret)-2]
	} else {
		return str
	}
}

func cleanJsonFromFileBytesOfIpfs(str string) string {
	example := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &example)
	if err != nil {
		re := regexp.MustCompile("\\{(.*)\\}")
		strs := re.FindAllString(str, -1)
		ret := strs[0]
		return ret
	} else {
		return str
	}
}

func (c *configuration) SubmitInflators() error {
	sh := shell.NewShell(c.IpfsConnection)
	b, err := sh.BlockGet(c.IpfsInflators)
	if err != nil {
		return errors.New("The hash for the inflators is not correct: " + err.Error())
	}
	cleaned := cleanArrayJsonFromFileBytesOfIpfs(string(b))
	inflators := []Inflator{}
	err = json.Unmarshal([]byte(cleaned), &inflators)
	if err != nil {
		return errors.New("The json for the inflators is not correct: " + err.Error())
	}
	c.inflators = map[string]int{}
	for _, v := range inflators {
		pubB, err := v.Bytes()
		if err != nil {
			return errors.New("The inflator's public key " + string(v.PublicKeyHex) + " is not correct," + err.Error())
		}
		pub, err := crypto.UnmarshalPublicKey(pubB)
		if err != nil {
			return errors.New("The inflator's public key is not correct," + err.Error())
		}
		pubB, _ = pub.Bytes()
		c.inflators[string(pubB)] = 0
	}
	return nil
}

func (c *configuration) SubmitWatchers() error {
	sh := shell.NewShell(c.IpfsConnection)
	b, err := sh.BlockGet(c.IpfsWatchers)
	if err != nil {
		return errors.New("The hash for the inflators is not correct: " + err.Error())
	}
	cleaned := cleanArrayJsonFromFileBytesOfIpfs(string(b))
	watchers := []Watcher{}
	err = json.Unmarshal([]byte(cleaned), &watchers)
	if err != nil {
		return errors.New("The json for the inflators is not correct: " + err.Error())
	}
	c.watchers = map[string]int{}
	for _, v := range watchers {
		pubB, err := v.Bytes()
		if err != nil {
			return errors.New("The inflator's public key " + string(v.PublicKeyHex) + " is not correct," + err.Error())
		}
		_, err = crypto.UnmarshalPublicKey(pubB)
		if err != nil {
			return errors.New("The watcher's public key  is not correct")
		}
		c.watchers[string(pubB)] = 0
	}
	return nil
}

var Conf = configuration{}

func init() {
	Conf.IpfsConnection = "127.0.0.1:5001"
	Conf.AbciDaemon = "tcp://0.0.0.0:46658"
	Conf.WaitingRequestTime = 5
	Conf.inflators = map[string]int{}
	Conf.watchers = map[string]int{}
	Conf.Tax = Tax{}
	Conf.IpfsTax = ""
}
