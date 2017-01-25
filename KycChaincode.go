package main

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Region Chaincode implementation
type KycChaincode struct {
}

type KycData struct {
	EnrollId    string `json:"EnrollId"`
	UserName    string `json:"UserName"`
	LastUpdated string `json:"LastUpdated"`
	BankName    string `json:"BankName"`
	ExpiryDate  string `json:"ExpiryDate"`
	Source      string `json:"Source"`
	KycStatus   string `json:"KycStatus"`
}

func (t *KycChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err := stub.CreateTable("KYC", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "enrollId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "userName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "lastUpdated", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "bankName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "expiryDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "source", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating KYC table.")
	}
	return nil, nil
}

// Add user KYC data in Blockchain
func (t *KycChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "write" {
		return t.AddKyc(stub, args)
	}
	if function == "update" {
		return t.UpdateKyc(stub, args)
	}
	return nil, nil
}

func (t *KycChaincode) AddKyc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Need 6 arguments")
	}

	enrollId := args[0]
	userName := args[1]
	lastUpdated := args[2]
	bankName := args[3]
	expiryDate := args[4]
	source := args[5]
	ok, err := stub.InsertRow("KYC", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: enrollId}},
			&shim.Column{Value: &shim.Column_String_{String_: userName}},
			&shim.Column{Value: &shim.Column_String_{String_: lastUpdated}},
			&shim.Column{Value: &shim.Column_String_{String_: bankName}},
			&shim.Column{Value: &shim.Column_String_{String_: expiryDate}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding KYC record.")
	}
	return nil, nil
}

func (t *KycChaincode) UpdateKyc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Need 6 arguments")
	}

	enrollId := args[0]
	userName := args[1]
	lastUpdated := args[2]
	bankName := args[3]
	expiryDate := args[4]
	source := args[5]
	ok, err := stub.ReplaceRow("KYC", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: enrollId}},
			&shim.Column{Value: &shim.Column_String_{String_: userName}},
			&shim.Column{Value: &shim.Column_String_{String_: lastUpdated}},
			&shim.Column{Value: &shim.Column_String_{String_: bankName}},
			&shim.Column{Value: &shim.Column_String_{String_: expiryDate}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding KYC record.")
	}
	return nil, nil
}

func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var KycDataObj KycData
	
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting enrollId to query")
	}

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)

	row, err := stub.GetRow("KYC", columns)
	if err != nil {
		return nil, errors.New("Failed to query")
	}
	
	KycDataObj.EnrollId=row.Columns[0].GetString_()
	KycDataObj.UserName = row.Columns[1].GetString_()
	KycDataObj.LastUpdated = row.Columns[2].GetString_()
	KycDataObj.BankName= row.Columns[3].GetString_()
	KycDataObj.ExpiryDate= row.Columns[4].GetString_()
	

	lastDate, _ := time.Parse("2006-01-02", row.Columns[2].GetString_())
	if lastDate.After(time.Now()) == true {
		KycDataObj.Source= ""
		KycDataObj.KycStatus= "Expired"
	}
	if lastDate.After(time.Now()) == false  {
		KycDataObj.Source= row.Columns[5].GetString_()
		KycDataObj.KycStatus= "OK"
	}
	jsonAsBytes, _ := json.Marshal(KycDataObj)

	return jsonAsBytes, nil
}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
