# stock - get stock infomation
[![GoDoc](https://godoc.org/github.com/z-Wind/stock?status.png)](http://godoc.org/github.com/z-Wind/stock)

## Table of Contents

* [Installation](#installation)
* [Usage](#usage)
* [Example](#example)
* [Include](#Include)

## Installation

    $ go get github.com/z-Wind/stock

To build with two ways

    $ cd $GOPATH/src/github.com/z-Wind/stock
    $ make

(optional) To run unit tests:

    $ cd $GOPATH/src/github.com/z-Wind/stock
    $ make test

(optional) To clean all except source code:

    $ cd $GOPATH/src/github.com/z-Wind/stock
    $ make clean

## Usage

### Start server
    $ stock -addr host:port [-accountID id]

### Quote
    $ GET http://localhost:6060/quote?symbols={symbols}
- symbols: seperate by comma，like VTI,VBR
	
### PriceHistory
    $ GET http://localhost:6060/priceHistory?symbols={symbols}
- symbols: seperate by comma，like VTI,VBR

### PriceAdjHistory
    $ GET http://localhost:6060/priceAdjHistory?symbols={symbols}
- symbols: seperate by comma，like VTI,VBR

### SavedOrder
> just for TD Ameritrade, should set accountID

    $ GET http://localhost:6060/savedOrder
	
	$ DELETE http://localhost:6060/savedOrder?savedOrderID={savedOrderID}
	
	$ POST JSON http://localhost:6060/savedOrder
	JSON Format
    {
      "Symbol": "string",
      "AssetType": "string[1]",
      "Instruction": "string[1]",
      "Price": 0,
      "Qunatity": 0
    }	
[1] [TD API Create Saved Order](https://developer.tdameritrade.com/account-access/apis/post/accounts/%7BaccountId%7D/savedorders-0)

## Example

### Start server

    $ cd $GOPATH/src/github.com/z-Wind/stock
    $ ./stock -addr localhost:6060 

### Simple Demo
    go to http://localhost:6060/

### Quote
    $ curl -X GET http://localhost:6060/quote?symbols=VTI,VBR,0050.tw,6564.two
	
### PriceHistory
    $ curl -X GET http://localhost:6060/priceHistory?symbols=VTI,VBR,0050.tw,6564.two
	
### PriceAdjHistory
    $ curl -X GET http://localhost:6060/priceAdjHistory?symbols=VTI,VBR,0050.tw,6564.two
	
### SavedOrder
> for just TD Ameritrade, should set accountID

    $ curl -X GET http://localhost:6060/savedOrder

## Include
- [gotd](https://github.com/z-Wind/gotd)
- [twse](https://github.com/z-Wind/twse)
- [alphavantage](https://github.com/z-Wind/alphavantage)
