package graphapi

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"go.infratographer.com/x/gidx"
)

type Node struct {
	ID        gidx.PrefixedID `json:"id"`
	GraphType *graphql.Object
}

func (r *Resolver) GetNode(id gidx.PrefixedID) *Node {
	fmt.Println("Inside GetNode")
	fmt.Printf("Get Node was called with ID: %s prefix is: %s\n", id, id.Prefix())
	for key, obj := range r.prefixMap {
		fmt.Printf("Prefix has %s => %s\n", key, obj.Name())
	}
	if resType, ok := r.prefixMap[id.Prefix()]; ok {
		fmt.Printf("Get Node found prefix: %s is type: %s\n", id.Prefix(), resType.Name())
		return &Node{
			ID:        id,
			GraphType: resType,
		}
	}

	return &Node{
		ID:        id,
		GraphType: r.prefixMap[id.Prefix()],
	}

	// fmt.Printf("Get Node couldn't find prefix in map for %s\n", id.Prefix())
	// return nil
}
