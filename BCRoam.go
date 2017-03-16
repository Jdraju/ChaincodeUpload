/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License .
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var msisdnList[] String
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// This is our structure for the broadcaster creating bulk inventory

type rsDetailBlock struct {
	PublicKey   string    `json:"publickey"`
	MSISDN      string    `json:"msisdn"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	HO          string    `json:"ho"`
	RP          string    `json:"rp"`
	Roaming     string    `json:"roaming"`
	Location    string    `json:"location"`
	RateType    string    `json:"ratetype"`
	Action      string    `json:"action"`
	TransType   string    `json:"transtype"`
	Destination string    `json:"destination"`
	Duration    string    `json:"duration"`
	Charges     string    `json:"charges"`
	Flag        string    `json:"flag"`
	Time        time.Time `json:"time"`
}

type rsDetail struct {
	MSISDN     string    `json:"msisdn"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	HO         string    `json:"ho"`
	RP         string    `json:"rp"`
	Roaming    string    `json:"roaming"`
	Location   string    `json:"location"`
	Plan       string    `json:"plan"`
	VoinceOutL string    `json:"voinceOutL"`
	VoinceInL  string    `json:"voinceInL"`
	DataL      string    `json:"float64"`
	VoiceOutR  string    `json:"voiceOutR"`
	VoiceInR   string    `json:"voiceInR"`
	DataR      time.Time `json:"dataR"`
}


//This is a helper structure to point to allPeers
type AllPeers struct {
	PeerName []string `json:"peerName"`
}

//For Debugging
func showArgs(args []string) {

	for i := 0; i < len(args); i++ {
		fmt.Printf("\n %d) : [%s]", i, args[i])
	}
	fmt.Printf("\n")
}

// Init function
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//var err error

	fmt.Println("Launching Init Function")
	//To add Time Stamp
	currentDateStr := time.Now().Format(time.RFC822)
	currtime, _ := time.Parse(time.RFC822, currentDateStr)
	//Inventory hard coded here
	rs1 := rsDetailBlock{"rs1", "14691234567", "A", "DALLAS", "AT&T", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}
	rs2 := rsDetailBlock{"rs2", "14691234568", "A", "DALLAS", "AT&T", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}
	rs3 := rsDetailBlock{"rs3", "14691234569", "A", "DALLAS", "AT&T", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}
	rs4 := rsDetailBlock{"rs4", "14691234570", "A", "DALLAS", "AT&T", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}
	rs5 := rsDetailBlock{"rs5", "349091234567", "A", "BARCELONA", "VODAFONE", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}
	rs6 := rsDetailBlock{"rs6", "349091234568", "A", "BARCELONA", "VODAFONE", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}
	rs7 := rsDetailBlock{"rs7", "349091234569", "A", "BARCELONA", "VODAFONE", "", "FALSE", "DALLAS", "", "", "", "", "", "", "", currtime}

	//Create array for all adspots in ledger
	//var AllPeersArray AllPeers

	t.putMSIDN(stub, rs1, rs1.PublicKey)
	t.putMSIDN(stub, rs2, rs2.PublicKey)
	t.putMSIDN(stub, rs3, rs3.PublicKey)
	t.putMSIDN(stub, rs4, rs4.PublicKey)
	t.putMSIDN(stub, rs5, rs5.PublicKey)
	t.putMSIDN(stub, rs6, rs6.PublicKey)
	t.putMSIDN(stub, rs7, rs7.PublicKey)

	fmt.Println("Init Function Complete")
	return nil, nil
}

//Invoke function

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function :%v", function)

	showArgs(args)
	var key, sp, loc, destmsisdn string

	// Handle different functions
	if function == "discoverRP" {
		fmt.Printf("Function is discoverRP")
		key = args[0]
		sp = args[1]
		loc = args[2]
		return t.discoverRP(stub, key, sp, loc)
	} else if function == "enterData" {
		fmt.Printf("Function is enterData")
		return t.enterData(stub, args)
	} else if function == "authentication" {
		fmt.Printf("Function is authentication")
		key = args[0]
		return t.authentication(stub, key)
	} else if function == "updateRates" {
		fmt.Printf("Function is updateRates")
		key = args[0]
		return t.updateRates(stub, key)
	} else if function == "CallOut" {
		fmt.Printf("Function is CallOut")
		key = args[0]
		destmsisdn = args[1]
		return t.CallOut(stub, key, destmsisdn)
	} else if function == "CallEnd" {
		fmt.Printf("Function is CallEnd")
		key = args[0]
		return t.CallEnd(stub, key)
	} else if function == "CallPay" {
		fmt.Printf("Function is CallPay")
		key = args[0]
		return t.CallPay(stub, key)
	}
	return nil, errors.New("Received unknown function invocation")
}

//QUERY FUNCTION
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("======== Query called, determining function")

	showArgs(args)

	if function == "queryMSISDN" {
		fmt.Printf("Function is queryPeers")
		return t.queryMSISDN(stub, args)
	} else {
		fmt.Printf("Invalid Function!")
	}

	return nil, nil
}

////////////////////////////////////////////////////

//Redirect FUNCTIONS

//Query MSISDN in our network
func (t *SimpleChaincode) queryMSISDN(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("queryMSISDN called")
	var key string
	key = args[0]
	fmt.Println("Key: %v", key)
	bytes, _ := stub.GetState(key)
	fmt.Println(string(bytes))
	fmt.Printf("%x", bytes)
	/*var rsDetailObj rsDetail
	err := json.Unmarshal(bytes, &rsDetailobj)
	if err != nil{
		fmt.Printf("Error in Unmarshalling")
	} else {
		fmt.Printf("Name:%v",rsDetailObj.name)
		fmt.Printf("Addr:%v",rsDetailObj.address)
		fmt.Printf("HO:%v",rsDetailObj.ho)
		fmt.Printf("RP:%v",rsDetailObj.rp)
		fmt.Printf("Roaming:%v",rsDetailObj.roaming)
		fmt.Printf("Location:%v",rsDetailObj.location)
		fmt.Printf("Plan:%v",rsDetailObj.plan)
		fmt.Printf("VoiceOutLocal:%v",rsDetailObj.voinceOutL)
		fmt.Printf("VoiceInLocal:%v",rsDetailObj.voinceInL)
		fmt.Printf("DataLocal:%v",rsDetailObj.dataL)
		fmt.Printf("VoiceOutRoam:%v",rsDetailObj.voiceOutR)
		fmt.Printf("VoiceInRoam:%v",rsDetailObj.voiceInR)
		fmt.Printf("DataRoam:%v",rsDetailObj.dataR)
	}*/
	return bytes, nil
}

func (t *SimpleChaincode) enterData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var rsDetailObj rsDetailBlock
	rsDetailObj.PublicKey = args[0]
	rsDetailObj.MSISDN = args[1]
	rsDetailObj.Name = args[2]
	rsDetailObj.Address = args[3]
	rsDetailObj.HO = args[4]
	rsDetailObj.RP = args[5]
	rsDetailObj.Roaming = args[6]
	rsDetailObj.Location = args[7]
	rsDetailObj.RateType = args[9]
	rsDetailObj.Action = args[10]
	rsDetailObj.TransType = args[11]
	rsDetailObj.Destination = args[12]
	rsDetailObj.Duration = args[13]
	rsDetailObj.Charges = args[14]
	rsDetailObj.Flag = args[15]
	msisdnList = append(msisdnList,rsDetailObj.MSISDN)
	//Get Current Time
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailObj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	
	fmt.Println(rsDetailObj)
	bytes, _ := json.Marshal(rsDetailObj)
	fmt.Println(string(bytes))
	//fmt.Printf("%x",bytes)
	err2 := stub.PutState(rsDetailObj.PublicKey, bytes)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in rsDetailObj")
	} else {
		fmt.Println("Success -  works")
	}

	//showArgs(args)

	return nil, nil
}

//putNetworkPeers: To put an array containing pointers to all blocks for a particular user(or peer) on the ledger
func (t *SimpleChaincode) putMSIDN(stub shim.ChaincodeStubInterface, rs rsDetailBlock, key string) ([]byte, error) {
	//marshalling
	fmt.Println(" Initializing msisdn: ", key)
	fmt.Printf("put details: %+v ", rs)
	fmt.Printf("\n")
	bytes, _ := json.Marshal(rs)
	fmt.Println(string(bytes))
	err2 := stub.PutState(key, bytes)
	msisdnList = append(msisdnList,rs.MSISDN)
	
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
		//return nil, err
	} else {
		fmt.Println("Success - Marshall in msisdn details")
	}
	return nil, nil
}

//Remote Partner Discovery
func (t *SimpleChaincode) discoverRP(stub shim.ChaincodeStubInterface, key string, sp string, loc string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.RP = sp
	rsDetailobj.Location = loc
	rsDetailobj.Action = "Discovery"
	rsDetailobj.TransType = "Setup"
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

//Authentication
func (t *SimpleChaincode) authentication(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var ho, rp, msisdn string

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	ho = rsDetailobj.HO
	rp = rsDetailobj.RP
	msisdn = rsDetailobj.MSISDN
	//ADDING LOGI FOR FRAUD:
	for x := 0; x < len(msisdnList);x++ { 
		if msisdn == msisdnList[x]{
			rsDetailobj.Flag="Fraud"
		}
	}
	
	////// Add logic for authentication here
	if rp == "VODAFONE" {
		if ho == "AT&T" {
			rsDetailobj.Roaming = "True"
			rsDetailobj.Action = "Authentication"
			rsDetailobj.TransType = "Setup"
			fmt.Println("Authentication suucessfull")
		}
	} else if rp == "AT&T" {
		if ho == "VODAFONE" {
			rsDetailobj.Roaming = "True"
			rsDetailobj.Action = "Authentication"
			rsDetailobj.TransType = "Setup"
			fmt.Println("Authentication Successfull")
		}
	}else {
		fmt.Println("Authentication Failed")
	}

	//rsDetailobj.Roaming="True"
	//rsDetailobj.Action="Authentication"
	//rsDetailobj.TransType="Setup"

	////////////////////////////////////////////
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

//Update voice and data rates
func (t *SimpleChaincode) updateRates(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	var sp string
	err = json.Unmarshal(bytes, &rsDetailobj)
	if rsDetailobj.Roaming == "True" {
		sp = rsDetailobj.RP
		if sp == "VODAFONE" {
			rsDetailobj.RateType = "Roaming"
		}
	}
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	rsDetailobj.Action = "Register"
	rsDetailobj.TransType = "Setup"
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

//Call Out
func (t *SimpleChaincode) CallOut(stub shim.ChaincodeStubInterface, key string, destmsisdn string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.Destination = destmsisdn
	rsDetailobj.Action = "Call Initialization"
	rsDetailobj.TransType = "Call Out"
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if rsDetailobj.PublicKey == "rs5"
	{
		rsDetailobj.Flag= "OVERAGE"
	}
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

func (t *SimpleChaincode) Overage(stub shim.ChaincodeStubInterface, key string, destmsisdn string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.Destination = destmsisdn
	rsDetailobj.Action = "OverageCheck"
	rsDetailobj.TransType = "Call Out"
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	rsDetailobj.Flag= "OVERAGE"
	

	return nil, nil
}

//Call In
func (t *SimpleChaincode) CallIn(stub shim.ChaincodeStubInterface, key string, destmsisdn string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.Destination = destmsisdn
	rsDetailobj.Action = "Call Recieved"
	rsDetailobj.TransType = "Call In"
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

//Call End
func (t *SimpleChaincode) CallEnd(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.Action = "Call End"
	rsDetailobj.TransType = "Call Out"
	currentDateStr := time.Now().Format(time.RFC822)
	duration := time.Since(rsDetailobj.Time)
	dur := strconv.FormatFloat(duration.Minutes(), 'E', -1, 64)
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	rsDetailobj.Duration = string(dur)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

//Call Pay
func (t *SimpleChaincode) CallPay(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", key)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", key)
	}

	var rsDetailobj rsDetailBlock
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.Action = "Pay Charge"
	rsDetailobj.TransType = "Call Out"
	currentDateStr := time.Now().Format(time.RFC822)
	rsDetailobj.Charges = "10"
	rsDetailobj.Time, _ = time.Parse(time.RFC822, currentDateStr)
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(rsDetailobj.PublicKey, bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
	} else {
		fmt.Println("Success, updated record")
	}

	return nil, nil
}

//MAIN FUNCTION
func main() {
	err := shim.Start(new(SimpleChaincode))

	fmt.Printf("IN MAIN of TelcoChaincode")
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
