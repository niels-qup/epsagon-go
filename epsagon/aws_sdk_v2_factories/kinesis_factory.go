package epsagonawsv2factories

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/epsagon/epsagon-go/protocol"
	"github.com/epsagon/epsagon-go/tracer"
	"reflect"
)

// KinesisEventDataFactory creates epsagon Resource from aws.Request to Kinesis
func KinesisEventDataFactory(r *aws.Request, res *protocol.Resource, metadataOnly bool) {
	inputValue := reflect.ValueOf(r.Params).Elem()
	streamName, ok := getFieldStringPtr(inputValue, "StreamName")
	if !ok {
		tracer.AddExceptionTypeAndMessage("aws-sdk-go",
			"kinesisEventDataFactory: couldn't find StreamName")
	}
	res.Name = streamName
	updateMetadataFromValue(inputValue, "PartitionKey", "partition_key", res.Metadata)
	if !metadataOnly {
		dataField := inputValue.FieldByName("Data")
		if dataField != (reflect.Value{}) {
			res.Metadata["data"] = string(dataField.Bytes())
		}
	}
	handleSpecificOperation(r, res, metadataOnly,
		map[string]specificOperationHandler{
			"PutRecord": handleKinesisPutRecord,
		}, nil,
	)
}

func handleKinesisPutRecord(r *aws.Request, res *protocol.Resource, metadataOnly bool) {
	outputValue := reflect.ValueOf(r.Data).Elem()
	updateMetadataFromValue(outputValue, "ShardId", "shared_id", res.Metadata)
	updateMetadataFromValue(outputValue, "SequenceNumber", "sequence_number", res.Metadata)
}