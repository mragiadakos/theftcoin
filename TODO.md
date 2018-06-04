The theftcoin is a blockchain based on the tendermint to exchange coins that taxed on transactions.
The tax will be just a number in percentage  submitted to the validator
All the validators need to start with the same tax.
The tax will be a number with a public key that the tax will go in

The user will request the latest tax from the validator.
The user will add the latest tax in the transaction and the validator will validate it.
After the successful validation, the validator will put the coins of the tax to the submitted public key that came with the tax.

There will be also two other type of administrators listed in an IPFS file
- The inflators will be public keys, that can add coins to the blockchain
- The watchers will be public keys, that can query others people coins and transactions 

POST /Delivery
REQUEST
Signature: signature
Data: {
    From : public key
    To: *public key // will be empty for ADD and REMOVE
    Action: string
    TaxHash : *string // will be empty for ADD and REMOVE
    Date: UTC
}
RESPONSE:
  Error scenarios:
    - the signature is not correct
    For ADD_ACTION
        - the user is not listed in the inflators
    For REMOVE_ACTION
        - the user is not listed in the inflators
    For SEND_ACTION
        - the 'TaxHash' is not correct
        - the coin transfer do not fit with the money that the user has
        - the date of the transaction is old (passed 5 seconds)
        
        
POST /query
REQUEST
Signature: signature
Data: {
   From: public key
   Nonce: string
   Date: UTC
   User: *public key 
   Coins: float64
}
RESPONSE
    Error scenarios:
    - the signature is not correct
    - the date of the transaction is old (passed 5 seconds)
    - the requester is not in the list of watchers to request coins or the transaction of the user

    Output 
    {
        Coins: number
    }

