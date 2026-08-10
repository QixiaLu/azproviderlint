// Package pluginsdk is a minimal stand-in for the terraform-plugin-sdk helper/resource
// package used only by the AZR007 analysistest fixtures.
package pluginsdk

import (
	"context"
	"time"
)

// StateChangeConf is a minimal stand-in for pluginsdk.StateChangeConf.
type StateChangeConf struct {
	Pending []string
	Target  []string
	Timeout time.Duration
}

// WaitForStateContext is a stub matching the real helper's signature.
func (c *StateChangeConf) WaitForStateContext(ctx context.Context) (interface{}, error) {
	return nil, nil
}
