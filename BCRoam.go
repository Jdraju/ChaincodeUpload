
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

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// This is our structure for the broadcaster creating bulk inventory
type rsDetail struct {
	name       string  `json:"name"`
	address    string  `json:"address"`
	ho         string  `json:"ho"`
	rp         string  `json:"rp"`
	roaming    string  `json:"roaming"`
	location   string  `json:"location"`
	plan       string  `json:"plan"`
	voinceOutL string `json:"voinceOutL"`
	voinceInL  string `json:"voinceInL"`
	dataL      string  `json:"float64"`
	voiceOutR  string `json:"voiceOutR"`
	voiceInR   string `json:"voiceInR"`
	dataR      string `json:"dataR"`
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
	
	//Peers hard coded here
	ho := "ATT"
	rp := "Vodafone"

	//Create array for all adspots in ledger
	var AllPeersArray AllPeers

	t.putNetworkPeers(stub, AllPeersArray, ho)
	t.putNetworkPeers(stub, AllPeersArray, rp)

	fmt.Println("Init Function Complete")
	return nil, nil
}

//Invoke function

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function :%v",function)

	showArgs(args)
	var msisdn,val string

	// Handle different functions
	if function == "discoverRP" {
		fmt.Printf("Function is discoverRP")
		val=args[0]
		msisdn=args[1]
		return t.discoverRP(stub,msisdn,val)
	} else if function == "enterData" {
		fmt.Printf("Function is enterData")
		return t.enterData(stub,args)
	} else if function == "roamOnOff" {
		fmt.Printf("Function is roamOnOff")
		msisdn=args[0]
		return t.roamOnOff(stub, msisdn)
	}else if function == "updateRates" {
		fmt.Printf("Function is updateRates")
		val=args[0]
		msisdn=args[1]
		return t.updateRates(stub,msisdn,val)
	} 
        return nil, errors.New("Received unknown function invocation")
}

//QUERY FUNCTION
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("======== Query called, determining function")

	showArgs(args)

	if function == "queryPeers" {
		fmt.Printf("Function is queryPeers")
		return t.queryPeers(stub, args)
	} else {
		fmt.Printf("Invalid Function!")
	}

	return nil, nil
}
////////////////////////////////////////////////////

//Redirect FUNCTIONS
func (t *SimpleChaincode) queryPeers(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
	fmt.Println("queryPeers called")
	var user string
	user= args[0]
	fmt.Println("User name: %v",user)
	bytes,_:= stub.GetState(user)
	var peer string
	err := json.Unmarshal(bytes, &peer)
	if err != nil{
		fmt.Printf("Error in Unmarshalling")
	} else {
		fmt.Printf("Peer name: %v",peer)
	}
	return nil,nil
}

func (t *SimpleChaincode) enterData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var msisdn=args[0]
	var rsDetailObj rsDetail
		rsDetailObj.name = args[1]
		rsDetailObj.address = args[2]
		rsDetailObj.ho = args[3]
		rsDetailObj.rp = args[4]   
		rsDetailObj.roaming = args[5]
		rsDetailObj.location = args[6]
		rsDetailObj.plan = args[7]
		rsDetailObj.voinceOutL = args[8]
		rsDetailObj.voinceInL = args[9]
		rsDetailObj.dataL = args[10]     
		rsDetailObj.voiceOutR = args[11]
		rsDetailObj.voiceInR = args[12]
		rsDetailObj.dataR = args[13]
	
		bytes, _ := json.Marshal(rsDetailObj)
		err := stub.PutState(msisdn, bytes)
		if err != nil {
			fmt.Println("Error - could not Marshall in rsDetailObj")
		} else {
			fmt.Println("Success -  works")
		}
	
	showArgs(args)

	return nil, errors.New("Received unknown function invocation")
}


//putNetworkPeers: To put an array containing pointers to all blocks for a particular user(or peer) on the ledger
func (t *SimpleChaincode) putNetworkPeers(stub shim.ChaincodeStubInterface, allPeersObj AllPeers, userId string) ([]byte, error) {
	//marshalling
	fmt.Println("Launching putNetworkPeers helper function userid: ", userId)
	fmt.Printf("putNetworkPeers: %+v ", allPeersObj)
	fmt.Printf("\n")
	bytes, _ := json.Marshal(allPeersObj)
	err2 := stub.PutState(userId, bytes)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in putNetworkPeers")
		//return nil, err
	} else {
		fmt.Println("Success - Marshall in putNetworkPeers")
	}
	fmt.Println("putNetworkPeers Function Complete - userid: ", userId)
	return nil, nil
}

//Remote Partner Discovery
func (t *SimpleChaincode) discoverRP(stub shim.ChaincodeStubInterface, msisdn string,val string) ([]byte, error) {

	bytes, err := stub.GetState(msisdn)
	if err != nil {
		fmt.Println("Error - Could not get User details : %s", msisdn)
		//return nil, err
	} else {
		fmt.Println("Success - User details found %s", msisdn)
	}

	var rsDetailobj rsDetail
	err = json.Unmarshal(bytes, &rsDetailobj)
	rsDetailobj.rp=val
	bytes2, _ := json.Marshal(rsDetailobj)
	err2 := stub.PutState(msisdn,bytes2)
	if err2 != nil {
		fmt.Println("Error - could not Marshall in msisdn")
		//return nil, err
	} else {
		fmt.Println("Success, updated record")
	}
	
	return nil,nil
}

func (t *SimpleChaincode) roamOnOff(stub shim.ChaincodeStubInterface, msisdn string) ([]byte, error) {
	return nil,nil
}

func (t *SimpleChaincode) updateRates(stub shim.ChaincodeStubInterface, msisdn string, val string) ([]byte, error) {
	return nil,nil
}


	
//MAIN FUNCTION
func main() {
	err := shim.Start(new(SimpleChaincode))

	fmt.Printf("IN MAIN of TelcoChaincode")
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
		if(roaming == "true"){
			if(ho=="ATT" && rp=="Vodafone") 
			{ 
				rates = rates*2
			}
			if (ho=="Vodafone" && rp=="ATT")
			{
				rates = rates*3
			}
		}
	}
}
