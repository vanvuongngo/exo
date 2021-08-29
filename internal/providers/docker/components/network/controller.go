package network

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (n *Network) InitResource() error {
	if err := n.UnmarshalSpec(&n.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(n.ComponentState, &n.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (n *Network) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(n.State)
}
