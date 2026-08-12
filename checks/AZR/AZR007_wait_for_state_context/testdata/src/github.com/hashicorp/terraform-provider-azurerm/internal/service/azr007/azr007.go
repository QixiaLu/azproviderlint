package azr007

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Should be flagged: pointer StateChangeConf composite literal.
func invalidStateChangeConf(ctx context.Context) {
	stateConf := &pluginsdk.StateChangeConf{ // want `StateChangeConf`
		Pending: []string{"Creating"},
		Target:  []string{"Created"},
		Timeout: 10 * time.Minute,
	}
	_, _ = stateConf.WaitForStateContext(ctx)
}

// Should be flagged: value StateChangeConf composite literal.
func invalidStateChangeConfNoPointer(ctx context.Context) {
	stateConf := pluginsdk.StateChangeConf{ // want `StateChangeConf`
		Pending: []string{"Deleting"},
		Target:  []string{"Deleted"},
		Timeout: 5 * time.Minute,
	}
	_, _ = stateConf.WaitForStateContext(ctx)
}

// Should be flagged: empty StateChangeConf composite literal.
func invalidEmptyStateChangeConf(ctx context.Context) {
	conf := &pluginsdk.StateChangeConf{} // want `StateChangeConf`
	_, _ = conf.WaitForStateContext(ctx)
}

// Should NOT be flagged: a custom poller with a different type.
type MyCustomPoller struct{}

func (p *MyCustomPoller) PollUntilDone(ctx context.Context) error {
	return nil
}

func validCustomPoller() {
	p := &MyCustomPoller{}
	_ = p.PollUntilDone(context.Background())
}
