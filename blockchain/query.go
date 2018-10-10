package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryRange(from string, to string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "queryRange")
	args = append(args, from)
	args = append(args, to)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return response.Payload, nil
}

func (setup *FabricSetup) Query(key string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, key)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return response.Payload, nil
}
