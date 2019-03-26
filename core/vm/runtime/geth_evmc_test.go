// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package runtime

import (
	"math/big"
	"testing"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)



var (
	ewasmfile string
	input string
	expected string
)

// The VM config for state tests that accepts --vm.* command line arguments.
var testVMConfig = func() vm.Config {
	vmconfig := vm.Config{}
	flag.StringVar(&vmconfig.EVMInterpreter, utils.EVMInterpreterFlag.Name, utils.EVMInterpreterFlag.Value, utils.EVMInterpreterFlag.Usage)
	flag.StringVar(&vmconfig.EWASMInterpreter, utils.EWASMInterpreterFlag.Name, utils.EWASMInterpreterFlag.Value, utils.EWASMInterpreterFlag.Usage)
	flag.StringVar(&ewasmfile, "ewasmfile", "", "ewasm file to run")
	flag.StringVar(&input, "input", "", "input to ewasm module, read with getCallData")
	flag.StringVar(&expected, "expected", "", "expected return data from ewasm module, compared to return data from finish")
	flag.Parse()
	return vmconfig
}()


func TestCallEwasm(t *testing.T) {
	state, _ := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	address := common.HexToAddress("0x0a")

	var (
		wasmbytes  []byte
		err  error
	)

	fmt.Println("ewasmfile:", ewasmfile)

	if len(ewasmfile) > 0 {
		wasmbytes, err = ioutil.ReadFile(ewasmfile)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Need to pass --ewasmfile arg!")
	}

	// fmt.Println("read wasmfile:", wasmbytes)
	
	// getCallDataSize.json
	//var code = common.Hex2Bytes("0061736d01000000010d036000017f60027f7f0060000002340208657468657265756d0f67657443616c6c4461746153697a65000008657468657265756$

	// return data test
	//var code = common.Hex2Bytes("0061736d0100000001090260027f7f0060000002130108657468657265756d0666696e6973680000030201010503010001071102066d656d6f72790200046d61696e00010a0a0108004100411410000b0b1a010041000b1400000000000000000000000000000000efbeadde")

	state.SetCode(address, wasmbytes)
	ewasmChainConfig := &params.ChainConfig{
		ByzantiumBlock: new(big.Int),
		EWASMBlock:	new(big.Int),
	};

	ret, _, err := Call(address, common.Hex2Bytes(input), &Config{ChainConfig: ewasmChainConfig, State: state, EVMConfig: testVMConfig, GasLimit: 1000000})
	if err != nil {
	t.Fatal("didn't expect error", err)
	}

	//var expected = "00000000000000000000000000000000efbeadde"

	fmt.Println("got return bytes:", common.Bytes2Hex(ret))
	if common.Bytes2Hex(ret) != expected {
		t.Error(fmt.Sprintf("Expected %v, got %v", expected, common.Bytes2Hex(ret)))
		return
	}

}



func BenchmarkCallEwasm(b *testing.B) {

	state, _ := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	address := common.HexToAddress("0x0a")

	// getCallDataSize.json
	//var code = common.Hex2Bytes("0061736d01000000010d036000017f60027f7f0060000002340208657468657265756d0f67657443616c6c4461746153697a65000008657468657265756$

	var code = common.Hex2Bytes("0061736d0100000001090260027f7f0060000002130108657468657265756d0666696e6973680000030201010503010001071102066d656d6f72790200046d61696e00010a0a0108004100411410000b0b1a010041000b1400000000000000000000000000000000efbeadde")

	state.SetCode(address, code)
	ewasmChainConfig := &params.ChainConfig{
		ByzantiumBlock: new(big.Int),
		EWASMBlock:	new(big.Int),
	};

	var (
		ret  []byte
		err  error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret, _, err = Call(address, nil, &Config{ChainConfig: ewasmChainConfig, State: state, EVMConfig: testVMConfig, GasLimit: 1000000})
	}
	b.StopTimer()
	//Check if it is correct
	if err != nil {
		b.Error(err)
		return
	}

	var expected = "00000000000000000000000000000000efbeadde"

	if common.Bytes2Hex(ret) != expected {
		b.Error(fmt.Sprintf("Expected %v, got %v", expected, common.Bytes2Hex(ret)))
		return
	}
	//fmt.Println("gas used:", startGas - contract.Gas)

}
