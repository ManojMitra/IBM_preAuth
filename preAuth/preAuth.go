package main

import (
	"bytes"
	"encoding/binary"
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
	PreAuthProvider Provider `json:"Provider"`
	PreAuthMember   Member   `json:"Member"`
	PreAuthService  Service  `json:"Service"`
	PreAuthPayer    Payer    `json:"Payer"`
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
	prov := &Provider{"John Smith", "XYZ Capitol avenue NY", "22322", "112-223-22222", "112-223-22223", "Susan Smith"}
	jsonAsBytes, _ := json.Marshal(prov)
	err := stub.PutState("PR001", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	prov = &Provider{"Steven Foss", "ABC Capitol avenue NY", "22321", "112-223-33333", "112-223-33334", "Susan Smith"}
	jsonAsBytes, _ = json.Marshal(prov)
	err = stub.PutState("PR002", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	prov = &Provider{"Tad Harison", "ABC Capitol avenue NY", "22323", "112-223-33344", "112-223-33345", "Robert Smith"}
	jsonAsBytes, _ = json.Marshal(prov)
	err = stub.PutState("PR003", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	prov = &Provider{"Albert", "ABC Capitol avenue NY", "22323", "112-223-33355", "112-223-33356", "Robert Smith"}
	jsonAsBytes, _ = json.Marshal(prov)
	err = stub.PutState("PR004", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running MM" + function)
	fmt.Println(args)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "writeProvider" {
		return t.write(stub, "writeProvider", args)
	} else if function == "writeMember" {
		return t.write(stub, "writeMember", args)
	} else if function == "writeService" {
		return t.write(stub, "writeService", args)
	} else if function == "writePayer" {
		return t.write(stub, "writePayer", args)
	} else if function == "writePreAuth" {
		return t.write(stub, "writePreAuth", args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running MM " + function)

	// Handle different functions
	if function == "readProvider" {
		return t.read(stub, "readProvider", args[0])
	} else if function == "readMember" {
		return t.read(stub, "readMember", args[0])
	} else if function == "readService" {
		return t.read(stub, "readService", args[0])
	} else if function == "readPayer" {
		return t.read(stub, "readPayer", args[0])
	} else if function == "readPreAuth" {
		return t.read(stub, "readPreAuth", args[0])
	}

	return nil, errors.New("Received unknown function query: " + function)
}

// Add to chaincode method
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Write called for " + function)
	var key, jsonResp string
	var err error
	var rcvdKey, chaincodeObj []byte

	if function == "writeProvider" {
		if len(args) == 6 {
			key = "PRO" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Provider Key is: " + key)
			pt := &Provider{args[0], args[1], args[2], args[3], args[4], args[5]}
			jsonAsBytes, _ := json.Marshal(pt)
			err = stub.PutState(key, jsonAsBytes)
			fmt.Println("Provider added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 6")
		}
	} else if function == "writeMember" {
		if len(args) == 4 {
			key = "MEM" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Member Key is: " + key)

			mem := &Member{args[0], args[1], args[2], args[3]}

			jsonAsBytes, _ := json.Marshal(mem)
			err = stub.PutState(key, jsonAsBytes)
			fmt.Println("Member added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 4")
		}
	} else if function == "writeService" {
		if len(args) == 12 {
			key = "SRV" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Service Key is: " + key)

			srv := &Service{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11]}

			jsonAsBytes, _ := json.Marshal(srv)
			err = stub.PutState(key, jsonAsBytes)
			fmt.Println("Member added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 12")
		}
	} else if function == "writePayer" {
		if len(args) == 9 {
			key = "PYR" + strconv.Itoa(rand.Intn(1000))
			fmt.Println("Service Key is: " + key)

			pyr := &Payer{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]}

			jsonAsBytes, _ := json.Marshal(pyr)
			err = stub.PutState(key, jsonAsBytes)
			fmt.Println("Member added successfully with key: " + key)
		} else {
			return nil, errors.New("Incorrect number of arguments. Expecting 9")
		}
	} else if function == "writePreAuth" {
		key = "PA" + strconv.Itoa(rand.Intn(10000000))
		fmt.Println("PreAuth Key is: " + key)

		preAuth := &PreAuthForm{}
		preAuth.PreAuthID = key
		preAuth.PreAuthStatus = "Submitted"

		if len(args) == 6 {
			rcvdKey, err = t.write(stub, "writeProvider", args)
			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			readKey := fmt.Sprintf("%s", rcvdKey)
			chaincodeObj, err = t.read(stub, "readProvider", readKey)
			buf := bytes.NewReader(chaincodeObj)
			var p Provider
			err := binary.Read(buf, binary.LittleEndian, &p)

			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			preAuth.PreAuthProvider = p
		} else if len(args) == 10 {
			rcvdKey, err = t.write(stub, "writeProvider", args[5:9])
			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			readKey := fmt.Sprintf("%s", rcvdKey)
			chaincodeObj, err = t.read(stub, "readMember", readKey)
			buf := bytes.NewReader(chaincodeObj)
			var mem Member
			err := binary.Read(buf, binary.LittleEndian, &mem)

			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			preAuth.PreAuthMember = mem
		} else if len(args) == 22 {
			rcvdKey, err = t.write(stub, "writeService", args[10:21])
			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			readKey := fmt.Sprintf("%s", rcvdKey)
			chaincodeObj, err = t.read(stub, "readService", readKey)
			buf := bytes.NewReader(chaincodeObj)
			var srv Service
			err := binary.Read(buf, binary.LittleEndian, &srv)

			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			preAuth.PreAuthService = srv
		} else if len(args) == 31 {
			rcvdKey, err = t.write(stub, "writePayer", args[22:30])
			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			readKey := fmt.Sprintf("%s", rcvdKey)
			chaincodeObj, err = t.read(stub, "readPayer", readKey)
			buf := bytes.NewReader(chaincodeObj)
			var pyr Payer
			err := binary.Read(buf, binary.LittleEndian, &pyr)

			if err != nil {
				jsonResp = "{\"Error\":\"Failed to perform operation}"
				return nil, errors.New(jsonResp)
			}
			preAuth.PreAuthPayer = pyr
		} else {
			return nil, errors.New("Incorrect number of arguments received.")
		}
		jsonAsBytes, _ := json.Marshal(preAuth)
		err = stub.PutState(key, jsonAsBytes)
	} else {
		return nil, errors.New("Received unknown function invocation: " + function)
	}

	if err != nil {
		return nil, err
	}

	return []byte(key), nil
}

// Read from chaincode method
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, function string, args string) ([]byte, error) {
	fmt.Println("Read called for " + function)
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	key = args
	fmt.Println("Reading for key: " + key)
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get object for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	if function == "readProvider" {
		var p Provider
		err = json.Unmarshal(valAsbytes, &p)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Provider for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Provider fetched successfully with key: ")
	} else if function == "readMember" {
		var m Member
		err = json.Unmarshal(valAsbytes, &m)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Member for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Member fetched successfully with key: ")
	} else if function == "readService" {
		var s Service
		err = json.Unmarshal(valAsbytes, &s)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Service for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Service fetched successfully with key: ")
	} else if function == "readPayer" {
		var pyr Payer
		err = json.Unmarshal(valAsbytes, &pyr)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get Payer for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("Payer fetched successfully with key: ")
	} else if function == "readPreAuth" {
		if err != nil {
			var p PreAuthForm
			err = json.Unmarshal(valAsbytes, &p)
			jsonResp = "{\"Error\":\"Failed to get PreAuth for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println("PreAuth fetched successfully with key: ")
	} else {
		return nil, errors.New("Received unknown function query: " + function)
	}

	fmt.Print(key + " : " + string(valAsbytes))
	return valAsbytes, nil
}
