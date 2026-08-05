package azr006

import "context"

type Client struct {
	StopContext context.Context
}

type ResourceData struct{}

func ForCreate(ctx context.Context, d *ResourceData) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// Should be flagged: ctx assigned directly from the meta object
func badRead(d *ResourceData, meta interface{}) error {
	ctx := meta.(*Client).StopContext // want `use a timeouts-wrapped StopContext \(timeouts\.ForCreate/ForCreateUpdate/ForRead/ForUpdate/ForDelete\) so Custom Timeouts are supported, instead of assigning ctx from meta directly`
	_ = ctx
	return nil
}

// Should NOT be flagged: wrapped StopContext with cancel
func goodCreate(d *ResourceData, meta interface{}) error {
	ctx, cancel := ForCreate(meta.(*Client).StopContext, d)
	defer cancel()
	_ = ctx
	return nil
}

// Should NOT be flagged: ctx from context package
func goodBackground() {
	ctx := context.Background()
	_ = ctx
}

// Should NOT be flagged: different variable name reading from meta
func goodOtherName(meta interface{}) {
	stopCtx := meta.(*Client).StopContext
	_ = stopCtx
}
