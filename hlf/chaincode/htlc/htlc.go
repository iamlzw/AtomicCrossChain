/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/sha256"
	"encoding/json"
	//"flag"
	"fmt"
	//"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	//contractapi.Contract
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "InitLedger"  {
		return s.InitLedger(stub)
	} else if function == "GetBalance" {
		// Make payment of X units from A to B
		return s.GetBalance(stub, args[0])
	} else if function == "MintToken" {
		// Deletes an entity from its state
		return s.MintToken(stub, args[0],args[1])
	} else if function == "BurnToken" {
		// the old "Query" is now implemtned in invoke
		return s.BurnToken(stub, args[0],args[1])
	} else if function == "transfer" {
		// the old "Query" is now implemtned in invoke
		return s.Transfer(stub, args[0],args[1],args[2])
	} else if function == "conditional" {
		// the old "Query" is now implemtned in invoke
		return s.TransferConditional(stub, args[0],args[1],args[2],args[3],args[4],args[5])
	} else if function == "GetHashTimeLock" {
		// the old "Query" is now implemtned in invoke
		return s.GetHashTimeLock(stub, args[0])
	} else if function == "commit" {
		// the old "Query" is now implemtned in invoke
		return s.Commit(stub, args[0],args[1])
	} else if function == "revert" {
		// the old "Query" is now implemtned in invoke
		return s.Revert(stub, args[0])
	}

	return shim.Error("Invalid invoke function name")
}

// InitLedger
func (s *SmartContract) InitLedger(stub shim.ChaincodeStubInterface) pb.Response {

	dataByteBob := []byte("300")

	err := stub.PutState("Bob", dataByteBob)
	if err != nil {
		return shim.Error("Init Bob balance fail")
	}

	dataByte := []byte("200")

	err = stub.PutState("Alice", dataByte)

	if err != nil {
		return shim.Error("Init Alice balance fail")
	}

	return shim.Success(nil)
}

// Get Balance of an account
func (s *SmartContract) GetBalance (stub shim.ChaincodeStubInterface, id string) pb.Response {

	retValue, _ := stub.GetState(id)
	//intVar , _ := strconv.Atoi(string(retValue))

	//	retValueString := string(retValue)

	return shim.Success(retValue)
}


// Mint token
func (s *SmartContract) MintToken (stub shim.ChaincodeStubInterface, id string, amount string) pb.Response {

	balanceInt,_ := strconv.Atoi(string(s.GetBalance(stub, id).Payload))

	intAmount , _ := strconv.Atoi(amount)

	balanceInt = balanceInt + intAmount

	balanceString := strconv.Itoa(balanceInt)

	dataByte := []byte(balanceString)

	stub.PutState(id, dataByte)

	return shim.Success([]byte("success"))
}

// Mint token
func (s *SmartContract) BurnToken (stub shim.ChaincodeStubInterface, id string, amount string)  pb.Response {

	balanceInt,_ := strconv.Atoi(string(s.GetBalance(stub, id).Payload))

	intAmount , _ := strconv.Atoi(amount)

	if intAmount >= balanceInt {
		return shim.Success([]byte("error: not enough balance"))
	}

	balanceInt = balanceInt - intAmount

	balanceString := strconv.Itoa(balanceInt)

	dataByte := []byte(balanceString)

	stub.PutState(id, dataByte)

	return shim.Success([]byte("success"))
}


// Transfer from one accountto another one
func (s *SmartContract) Transfer (stub shim.ChaincodeStubInterface, fromId string, toId string, amount string) pb.Response {

	s.BurnToken(stub, fromId, amount)

	s.MintToken(stub, toId, amount)

	return shim.Success([]byte("success"))
}

// structure for the timelock
type HashTimeLock struct {
	LockID string `json:"lockid"`
	FromID string `json:"fromid"`
	ToID string `json:"toid"`
	Amount  string `json:"amount"`
	HashLock string `json:"hashlock"`
	TimeLock  string `json:"timelock"`
}

// Transfer from one accountto another one
func (s *SmartContract) TransferConditional (stub shim.ChaincodeStubInterface, lock_id string, from_id string, to_id string, amount string, hashlock string, timelock string) pb.Response {

	// decrease from the from amount
	s.BurnToken(stub, from_id, amount)

	// create HashTimeLock
	hashTimeLock := HashTimeLock {
		LockID:   lock_id,
		FromID:  from_id,
		ToID: to_id,
		Amount:  amount,
		HashLock: hashlock,
		TimeLock: timelock,
	}

	hashTimeLockAsBytes, _ := json.Marshal(hashTimeLock)

	err := stub.PutState(lock_id, hashTimeLockAsBytes)
	if err != nil {
		return shim.Error("put hash time lock fail")
	}
	return shim.Success([]byte("success"))
}

// Getting the created Hash time lock
func (s *SmartContract) GetHashTimeLock (stub shim.ChaincodeStubInterface, lock_id string) pb.Response {

	hashTimeLockAsBytes , _ := stub.GetState(lock_id)

	//hashTimeLock := new(HashTimeLock)
	//_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	return shim.Success(hashTimeLockAsBytes)
}


// Commiting the HTLC
func (s *SmartContract) Commit (stub shim.ChaincodeStubInterface, lock_id string, preimage string) pb.Response {

	hashTimeLockAsBytes , _ := stub.GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	hash := sha256.Sum256([]byte(preimage))

	hashString := fmt.Sprintf("%x", hash)

	fmt.Println("Hash String:", hashString)

	// condition 1 hash pre image

	if hashTimeLock.HashLock != hashString {

		fmt.Println("Invalid password:", hashString, hashTimeLock.HashLock)
		fmt.Println("Transaction reverted:")

		return shim.Error("invalid password")

	}

	// condition 2 time
	timestamp , _ := stub.GetTxTimestamp()
	timestampInt := timestamp.Seconds

	timelock , _ := time.Parse(time.RFC3339, hashTimeLock.TimeLock)

	if  timelock.Unix() < timestampInt {

		fmt.Println("Timelock already activated")
		fmt.Println("Actual transaction timestamp:", timestampInt)
		fmt.Println("Actual timelock:", timelock.Unix())
		fmt.Println("Transaction reverted")

		return shim.Error("Transaction reverted:Timelock already activated")
	}

	// increase amount to
	s.MintToken(stub, hashTimeLock.ToID, hashTimeLock.Amount)

	// delete lock
	stub.DelState(lock_id)

	// rasie event
	stub.SetEvent("Commit", []byte(preimage))

	fmt.Println("success commit")
	return shim.Success([]byte("success commit"))
}

// Revert the HTLC
func (s *SmartContract) Revert (stub shim.ChaincodeStubInterface, lock_id string) pb.Response {

	hashTimeLockAsBytes , _ := stub.GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	// condition 1 hash pre image - DOES NOT MATTER

	// condition 2 time

	timestamp , _ := stub.GetTxTimestamp()
	timestampInt := timestamp.Seconds

	timelock , _ := time.Parse(time.RFC3339, hashTimeLock.TimeLock)

	if  timelock.Unix() > timestampInt {

		fmt.Println("Timelock not yet activated")
		fmt.Println("Actual transaction timestamp:", timestampInt)
		fmt.Println("Actual timelock:", timelock.Unix())
		fmt.Println("Transaction reverted")

		return shim.Error("Transaction reverted:Timelock not yet activated")
	}


	// increase amount from

	s.MintToken(stub, hashTimeLock.FromID, hashTimeLock.Amount)

	fmt.Println("success commit")
	return shim.Success([]byte("success commit"))
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}