package utils

import (
	"github.com/bwmarrin/snowflake"
)

func GenerateId(init int64) (int64, error) {
	// create a snowflake Node
	node, err := snowflake.NewNode(init)
	if err != nil {
		return 0, err
	}

	// generate a snowflake id
	return node.Generate().Int64(), nil
}
