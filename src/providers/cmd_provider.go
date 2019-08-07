package provider

import (
	"bytes"
	"encoding/csv"
	"io"
	"os/exec"

	models "../models"
)

// CmdProvider command line provider
type CmdProvider struct{}

// Call calls a node and return the status of all nodes
func (ch CmdProvider) Call(address string) ([]models.NodeStatus, error) {
	cmd := exec.Command("cockroach", "node", "status", "--format=csv", "--host="+address, "--insecure")
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	response := output.Bytes()

	// Parse cmd results
	var status []models.NodeStatus
	csvReader := csv.NewReader(bytes.NewReader(response))
	record, err := csvReader.Read() // skip the header
	for {
		record, err = csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		status = append(status, models.NodeStatus{
			Id:          record[0],
			Address:     record[1],
			Build:       record[2],
			StartedAt:   record[3],
			UpdatedAt:   record[4],
			IsAvailable: record[5] == "true",
			IsLive:      record[6] == "true",
		})
	}
	return status, nil
}
