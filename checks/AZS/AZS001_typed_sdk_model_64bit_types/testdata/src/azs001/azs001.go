package azs001

// Should be flagged: non-64-bit numeric fields with tfschema tags
type BadModel struct {
	Count  int     `tfschema:"count"`  // want `property Count in model BadModel should be type int64, got int`
	Small  int16   `tfschema:"small"`  // want `property Small in model BadModel should be type int64, got int16`
	Medium int32   `tfschema:"medium"` // want `property Medium in model BadModel should be type int64, got int32`
	Ratio  float32 `tfschema:"ratio"`  // want `property Ratio in model BadModel should be type float64, got float32`
}

// Should be flagged: slices, pointers, and maps of non-64-bit numerics
type BadCompositeModel struct {
	Counts   []int              `tfschema:"counts"`   // want `property Counts in model BadCompositeModel should be type \[\]int64, got \[\]int`
	Ratios   []float32          `tfschema:"ratios"`   // want `property Ratios in model BadCompositeModel should be type \[\]float64, got \[\]float32`
	Optional *int32             `tfschema:"optional"` // want `property Optional in model BadCompositeModel should be type \*int64, got \*int32`
	Weights  map[string]float32 `tfschema:"weights"`  // want `property Weights in model BadCompositeModel should be type map\[string\]float64, got map\[string\]float32`
	Nested   []*int             `tfschema:"nested"`   // want `property Nested in model BadCompositeModel should be type \[\]\*int64, got \[\]\*int`
}

// Should be flagged: multiple names in a single field declaration
type BadMultiNameModel struct {
	First, Second int `tfschema:"first"` // want `property First in model BadMultiNameModel should be type int64, got int` `property Second in model BadMultiNameModel should be type int64, got int`
}

// Should be flagged: named types and aliases resolve to their underlying type
type Capacity int
type Percentage = float32
type PortList []int32

type BadNamedModel struct {
	Capacity   Capacity   `tfschema:"capacity"`   // want `property Capacity in model BadNamedModel should be type int64, got Capacity \(underlying int\)`
	Percentage Percentage `tfschema:"percentage"` // want `property Percentage in model BadNamedModel should be type float64, got Percentage \(underlying float32\)`
	Ports      PortList   `tfschema:"ports"`      // want `property Ports in model BadNamedModel should be type \[\]int64, got PortList \(underlying \[\]int32\)`
	Capacities []Capacity `tfschema:"capacities"` // want `property Capacities in model BadNamedModel should be type \[\]int64, got \[\]Capacity \(underlying \[\]int\)`
}

// Should NOT be flagged: named types with 64-bit underlying types
type Size int64
type Weight = float64

type GoodNamedModel struct {
	Size   Size   `tfschema:"size"`
	Weight Weight `tfschema:"weight"`
}

// Should NOT be flagged: correct 64-bit types
type GoodModel struct {
	Count   int64              `tfschema:"count"`
	Ratio   float64            `tfschema:"ratio"`
	Name    string             `tfschema:"name"`
	Enabled bool               `tfschema:"enabled"`
	Counts  []int64            `tfschema:"counts"`
	Weights map[string]float64 `tfschema:"weights"`
	Nested  []NestedModel      `tfschema:"nested"`
}

type NestedModel struct {
	Value int64 `tfschema:"value"`
}

// Should NOT be flagged: no tfschema tags, not a typed SDK model
type NotAModel struct {
	Count int
	Ratio float32
	Sizes []int32
}

// Should NOT be flagged: tagged, but not tfschema
type JSONModel struct {
	Count int     `json:"count"`
	Ratio float32 `json:"ratio"`
}
