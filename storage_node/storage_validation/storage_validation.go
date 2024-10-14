package storage_validation

import (
	master "Awesome-DFS/master_connection"
	val "Awesome-DFS/protobuf/validation"
	"context"
)

func ValidateChunk(fileUuid string) {
	client := master.GetValidationClient()

	request := &val.ValidationRequest{FileUuid: fileUuid}

	client.Validate(context.Background(), request)
}
