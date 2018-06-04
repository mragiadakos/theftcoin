Here is the guide on how to enable the server using command lines

First we will create IPFS hashes for tax, watchers and inflators

$ ./server -create-demo-keys
Creating the inflators files:
The inflator's private key inflator_priv.json
The inflators' file for the server is inflators.json
IPFS Hash for inflators: QmQU8m7SuRKCbZ2kbbm2jwsRJyEcu8KEZAWXryt64UN4EV
Creating the tax files:
The taxer's private key file tax_priv.json
The taxer's public key file tax.json
IPFS Hash for tax: QmVnExTWSTb4eiaZzhFobPdxQFXNmEVQuauQyKtEyBXLuQ
Creating the watchers files:
The watcher's private key watcher_priv.json
The watchers' file for the server is watchers.json
IPFS Hash for watchers: QmdqadWayqNAkBNj24EcbQCsxcheFTDXtS2KLFUipGAJ3E


and we will use these hashes to start the server

$ ./server -inflators=QmQU8m7SuRKCbZ2kbbm2jwsRJyEcu8KEZAWXryt64UN4EV -watchers=QmdqadWayqNAkBNj24EcbQCsxcheFTDXtS2KLFUipGAJ3E -tax=QmVnExTWSTb4eiaZzhFobPdxQFXNmEVQuauQyKtEyBXLuQ
I[06-07|20:57:28.931] Starting ABCIServer                          module=abci-server impl=ABCIServer
I[06-07|20:57:28.931] Waiting for new connection...                module=abci-server 
^Ccaptured interrupt, exiting...
I[06-07|20:57:31.027] Stopping ABCIServer                          module=abci-server impl=ABCIServer

Now you need to enable the tendermint daemon so all the transaction saved in the blockchain.

To start the transactions you need to use the client.

Lets use the inflator's private key to add some coins
$ ./client a --key inflator_priv.json --coins 1000
The add was successful

To make sure that the coins added, we will look if the coins added
$ ./client q --key inflator_priv.json 
Coins:  1000

Lets generate a new user to send the money from the inflator to the user
$ ./client g --filename receiver.json
The generate was successful

To send the money, we will use the public key
$ cat receiver.json 
{"PublicKey":"08011220982feb614689a49874f39de38b47300dd52a69257c94bd052fade955310cd46c","PrivateKey":"CAESYDtePfzg2ZV6y2W54GRMV7vKKCZj8vRdeUorp5kB+LX3mC/rYUaJpJh0853ji0cwDdUqaSV8lL0FL63pVTEM1GyYL+thRomkmHTzneOLRzAN1SppJXyUvQUvrelVMQzUbA=="}

We will send money
$ ./client send --key inflator_priv.json  --receiver='08011220982feb614689a49874f39de38b47300dd52a69257c94bd052fade955310cd46c' --coins 100
Error: The tax is not included.

Ah! we need to include the tax's IPFS hash that the validators will validate
$ ./client send --key inflator_priv.json  --receiver='08011220982feb614689a49874f39de38b47300dd52a69257c94bd052fade955310cd46c' --coins 100 --tax QmVnExTWSTb4eiaZzhFobPdxQFXNmEVQuauQyKtEyBXLuQ
The send was successful

Now the receiver has 90 coins only because the tax took 10%
$ ./client q --key receiver.json 
Coins:  90


The taxer now has 10 coins
$ ./client q --key tax_priv.json 
Coins:  10


Now the receiver want to know how much the thief,... sorry I meant the taxer, how much he has accumulated
$ ./client q --key receiver.json  --user 080112203d722de979182ad5137370dd511d2de009fd9ffb274ea834f246378031abf892
Error:You are not a watcher.

Ah! only the watcher can see other's people money
Now we use the private key of the watcher to see the taxes
$ ./client q --key watcher_priv.json  --user 080112203d722de979182ad5137370dd511d2de009fd9ffb274ea834f246378031abf892
Coins:  10