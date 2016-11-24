package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Provider structure
type Provider struct {
	ProviderName          string `json:"providername"`
	ProviderAddr          string `json:"provideraddr"`
	ProviderCityZip       string `json:"providercityzip"`
	ProviderPhone         string `json:"providerphone"`
	ProviderFax           string `json:"providerfax"`
	ProviderContactPerson string `json:"providercontactperson"`
}

// Member structure
type Member struct {
	MemName string `json:"memname"`
	MemID   string `json:"memid"`
	MemDOB  string `json:"memdob"`
	MemDOR  string `json:"memdor"`
}

// Service structure
type Service struct {
	SrvRequested        string `json:"srvrequested"`
	SrvDOS              string `json:"srvdos"`
	SrvDiagnosis        string `json:"srvdiagnosis"`
	SrvCPTCode          string `json:"srvcptcode"`
	SrvICDCode          string `json:"srvicdcode"`
	SrvProviderFacility string `json:"srvproviderfacility"`
	SrvPhone            string `json:"srvphone"`
	SrvAddr             string `json:"srvaddr"`
	SrvCityZip          string `json:"srvcityzip"`
	SrvProcedure        string `json:"srvprocedure"`
	SrvProcOtherTxt     string `json:"srvprocothertxt"`
	SrvClinicalInfo     string `json:"srvclinicalinfo"`
}

// Payer structure
type Payer struct {
	PayerLOS         string `json:"payerLOS"`
	PayerProvTIN     string `json:"payerProvTIN"`
	PayerDOS         string `json:"payerDOS"`
	PayerBillTIN     string `json:"payerBillTIN"`
	PayerAmtAuth     string `json:"payerAmtAuth"`
	PayerDiag        string `json:"payerDiag"`
	PayerAllowedProc string `json:"payerAllowedProc"`
	PayerComment     string `json:"payerComment"`
	PayerAddDocReqd  string `json:"payerAddDocReqd"`
}

// PreAuthForm structure
type PreAuthForm struct {
	PreAuthID       string   `json:"preauthid"`
	PreAuthStatus   string   `json:"preauthstatus"`
	PreAuthProvider Provider `json:"provider"`
	PreAuthMember   Member   `json:"member"`
	PreAuthService  Service  `json:"service"`
	PreAuthPayer    Payer    `json:"payer"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	return t.writeDummyProvider(stub)
}

func (t *SimpleChaincode) writeDummyProvider(stub shim.ChaincodeStubInterface) ([]byte, error) {
	a := []string{"John Smith", "XYZ Capitol avenue NY", "22322", "112-223-22222", "112-223-22223", "Susan Smith"}
	t.write(stub, a)

	/*b := []string{"Steven Foss", "ABC Capitol avenue NY", "22321", "112-223-33333", "112-223-33334", "Susan Smith"}
	t.write(stub, b)

	c := []string{"Tad Harison", "ABC Capitol avenue NY", "22323", "112-223-33344", "112-223-33345", "Robert Smith"}
	t.write(stub, c)

	d := []string{"Albert", "ABC Capitol avenue NY", "22323", "112-223-33355", "112-223-33356", "Robert Smith"}
	t.write(stub, d)*/

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke called for " + function)

	// Handle different functions
	if function == "write" {
		return t.write(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running MM " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, "read", args)
	}

	return nil, errors.New("Received unknown function query: " + function)
}

// Add to chaincode
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//fmt.Println("Write called with args " + string(args))
	var key, jsonResp string
	var err error
	var rcvdKey, chaincodeObj []byte

	key = "PA" + strconv.Itoa(rand.Intn(10000000))
	fmt.Println("PreAuth Key is: " + key)

	preAuth := PreAuthForm{}
	preAuth.PreAuthID = key
	preAuth.PreAuthStatus = "Submitted"

	if len(args) == 6 {
		rcvdKey, err1 := t.writeStruct(stub, "writeProvider", args)
		if err1 != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("rcvdKeya: ", rcvdKey)
		readKey := fmt.Sprintf("%s", rcvdKey)
		fmt.Println("readKey: ", readKey)
		chaincodeObj, err = t.readStruct(stub, "readProvider", readKey)

		var p Provider
		err := json.Unmarshal(chaincodeObj, &p)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		preAuth.PreAuthProvider = p
	} else if len(args) == 10 {
		rcvdKey, err = t.writeStruct(stub, "writeProvider", args[5:9])
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		readKey := fmt.Sprintf("%s", rcvdKey)
		chaincodeObj, err = t.readStruct(stub, "readMember", readKey)

		var mem Member
		err := json.Unmarshal(chaincodeObj, &mem)

		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		preAuth.PreAuthMember = mem
	} else if len(args) == 22 {
		rcvdKey, err = t.writeStruct(stub, "writeService", args[10:21])
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		readKey := fmt.Sprintf("%s", rcvdKey)
		chaincodeObj, err = t.readStruct(stub, "readService", readKey)

		var srv Service
		err := json.Unmarshal(chaincodeObj, &srv)

		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		preAuth.PreAuthService = srv
	} else if len(args) == 31 {
		rcvdKey, err = t.writeStruct(stub, "writePayer", args[22:30])
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		readKey := fmt.Sprintf("%s", rcvdKey)
		chaincodeObj, err = t.readStruct(stub, "readPayer", readKey)

		var pyr Payer
		err := json.Unmarshal(chaincodeObj, &pyr)

		if err != nil {
			jsonResp = "{\"Error\":\"Failed to perform operation}"
			return nil, errors.New(jsonResp)
		}
		preAuth.PreAuthPayer = pyr
	} else {
		return nil, errors.New("Incorrect number of arguments received.")
	}
	jsonAsBytes, err1 := json.Marshal(preAuth)
	if err1 != nil {
		jsonResp := "{\"Error\":\"Failed to perform operation}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("PreAuthForm marshaled successfully with key: " + key)

	err = stub.PutState(key, jsonAsBytes)

	if err != nil {
		return nil, errors.New("Received unknown function invocation: ")
	}
	fmt.Println("PreAuthForm added successfully with key: " + key)
	return []byte(key), nil
}

// Add structure
func (t *SimpleChaincode) writeStruct(stub shim.ChaincodeStubInterface, structName string, args []string) ([]byte, error) {
	fmt.Println("Write called for " + structName)
	var key string

	if structName == "writeProvider" {
		if len(args) == 6 {
			key = "PRO" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Provider Key is: " + key)

			pt := &Provider{args[0], args[1], args[2], args[3], args[4], args[5]}

			jsonAsBytes, _ := json.Marshal(pt)
			err := stub.PutState(key, jsonAsBytes)
			if err != nil {
				jsonResp := "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			fmt.Println("Provider added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 6")
		}
	} else if structName == "writeMember" {
		if len(args) == 4 {
			key = "MEM" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Member Key is: " + key)

			mem := &Member{args[0], args[1], args[2], args[3]}

			jsonAsBytes, _ := json.Marshal(mem)
			err := stub.PutState(key, jsonAsBytes)
			if err != nil {
				jsonResp := "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			fmt.Println("Member added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 4")
		}
	} else if structName == "writeService" {
		if len(args) == 12 {
			key = "SRV" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Service Key is: " + key)

			srv := &Service{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11]}

			jsonAsBytes, _ := json.Marshal(srv)
			err := stub.PutState(key, jsonAsBytes)
			if err != nil {
				jsonResp := "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			fmt.Println("Member added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 12")
		}
	} else if structName == "writePayer" {
		if len(args) == 9 {
			key = "PYR" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Service Key is: " + key)

			pyr := &Payer{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]}

			jsonAsBytes, _ := json.Marshal(pyr)
			err := stub.PutState(key, jsonAsBytes)
			if err != nil {
				jsonResp := "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			fmt.Println("Member added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 9")
		}
	}

	retVal := []byte(key)
	return retVal, nil
}

// Read from chaincode method
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Read called for " + function)
	var key, jsonResp string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	if len(args) == 1 {
		key = args[0]
		fmt.Println("Reading for key: " + key)
		valAsbytes, err := stub.GetState(key)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get object for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}

		var p PreAuthForm
		err = json.Unmarshal(valAsbytes, &p)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get PreAuth for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("PreAuth fetched successfully with key: " + key + " : " + string(valAsbytes))
		return valAsbytes, nil
	} else if len(args) == 2 {
		structName := args[0]
		key = args[1]
		return t.readStruct(stub, structName, key)
	} else {
		return nil, errors.New("Received unknown function query: " + function)
	}
}

// Read from chaincode structure
func (t *SimpleChaincode) readStruct(stub shim.ChaincodeStubInterface, structName string, arg string) ([]byte, error) {
	fmt.Println("Read called for " + structName)
	var jsonResp string
	var err error

	key := arg
	fmt.Println("Reading for key: " + key)
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get object for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	if structName == "readProvider" {
		var p Provider
		err = json.Unmarshal(valAsbytes, &p)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Provider for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Provider fetched successfully with key: " + key)
	} else if structName == "readMember" {
		var m Member
		err = json.Unmarshal(valAsbytes, &m)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Member for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Member fetched successfully with key: ")
	} else if structName == "readService" {
		var s Service
		err = json.Unmarshal(valAsbytes, &s)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Service for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Service fetched successfully with key: " + key)
	} else if structName == "readPayer" {
		var pyr Payer
		err = json.Unmarshal(valAsbytes, &pyr)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Payer for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Payer fetched successfully with key: " + key)
	} else {
		return nil, errors.New("Received unknown structure or key in query: " + structName + ":" + key)
	}

	return valAsbytes, nil
}
