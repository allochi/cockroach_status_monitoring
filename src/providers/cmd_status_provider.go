package providers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	models "../models"
)

// CmdStatusProvider command line provider
type CmdStatusProvider struct {
	Address string
}

// Call calls a node and return the status of all nodes
func (sp CmdStatusProvider) Call() ([]models.Node, error) {
	if sp.Address == "" {
		return []models.Node{}, fmt.Errorf("cmd host address is required")
	}

	cmd := exec.Command("cockroach", "node", "status", "--format=csv", "--host="+sp.Address, "--insecure")
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	response := output.Bytes()
	nodes, err := sp.parseNodes(response)
	if err != nil {
		return []models.Node{}, err
	}

	return nodes, nil
}

func (sp CmdStatusProvider) parseNodes(content []byte) ([]models.Node, error) {
	var status []models.Node
	csvReader := csv.NewReader(bytes.NewReader(content))
	record, err := csvReader.Read() // skip the header
	for {
		record, err = csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []models.Node{}, err
		}
		id, _ := strconv.ParseInt(record[0], 10, 64)
		status = append(status, models.Node{
			ID:          id,
			Address:     record[1],
			StartedAt:   record[3],
			UpdatedAt:   record[4],
			IsAvailable: record[5] == "true",
			IsLive:      record[6] == "true",
		})
	}
	return status, nil
}
