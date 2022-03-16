package types

import "fmt"

// Pruning option string constants
const (
	PruningOptionDefault    = "default"
	PruningOptionEverything = "everything"
	PruningOptionNothing    = "nothing"
	PruningOptionCustom     = "custom"
)

type Type int

const (
	// Default defines a pruning strategy where the last 100,000 heights are
	// kept where to-be pruned heights are pruned at every 10th height.
	// The last 100000 heights are kept(approximately 1 week worth of state) assuming the typical
	// block time is 6s. If these values
	// do not match the applications' requirements, use the "custom" option.
	Default Type = iota
	// Everything defines a pruning strategy where all committed heights are
	// deleted, storing only the current height and last 10 states. To-be pruned heights are
	// pruned at every 10th height.
	Everything
	// Nothing defines a pruning strategy where all heights are kept on disk.
	// This is the only stretegy where KeepEvery=1 is allowed with state-sync snapshots disabled.
	Nothing
	// Custom defines a pruning strategy where the user specifies the pruning. The `pruning-keep-every`
	// for this strategy must be equal to `snapshot-interval`
	Custom
	// Undefined defines an undefined pruning strategy. It is to be returned by stores that do not support pruning.
	Undefined
)

// PruningOptions defines the pruning strategy used when determining which
// heights are removed from disk when committing state.
type PruningOptions struct {
	// KeepRecent defines how many recent heights to keep on disk.
	KeepRecent uint64

	// KeepEvery defines how many offset heights are kept on disk past KeepRecent.
	KeepEvery uint64

	// Interval defines when the pruned heights are removed from disk.
	Interval uint64

	Type Type
}

func NewPruningOptions(pruningType Type) *PruningOptions {
	switch pruningType {
		case Default:
			return &PruningOptions{
				KeepRecent: 100_000,
				KeepEvery:  0,
				Interval:   100,
				Type:       Default,
			}
		case Everything:
			return &PruningOptions{
				KeepRecent: 10,
				KeepEvery:  0,
				Interval:   10,
				Type:       Everything,
			}
		case Nothing:
			return &PruningOptions{
				KeepRecent: 0,
				KeepEvery:  1,
				Interval:   0,
				Type:       Nothing,
			}
		case Custom:
			return &PruningOptions{
				Type:       Custom,
			}
		default:
			return &PruningOptions{
				Type:       Undefined,
			}
	}
}

func NewCustomPruningOptions(keepRecent, keepEvery, interval uint64) *PruningOptions {
	return &PruningOptions{
		KeepRecent: keepRecent,
		KeepEvery:  keepEvery,
		Interval:   interval,
		Type:       Custom,
	}
}

func (po PruningOptions) GetType() Type {
	return po.Type
}

func (po PruningOptions) Validate() error {
	if po.KeepEvery == 0 && po.Interval == 0 {
		return fmt.Errorf("invalid 'Interval': %d, when pruning everything", po.Interval)
	}
	if po.KeepEvery == 1 && po.Interval != 0 { // prune nothing
		return fmt.Errorf("invalid 'Interval': %d when pruning nothing", po.Interval)
	}
	if po.KeepEvery > 1 && po.Interval == 0 {
		return fmt.Errorf("invalid 'Interval': %d when `KeepEvery`: %d", po.Interval, po.KeepEvery)
	}
	if po.KeepRecent > 0 && po.Interval == 0 {
		return fmt.Errorf("invalid 'Interval': %d when 'KeepRecent': %d", po.Interval, po.KeepRecent)
	}

	return nil
}

func NewPruningOptionsFromString(strategy string) *PruningOptions {
	switch strategy {
	case PruningOptionEverything:
		return NewPruningOptions(Everything)

	case PruningOptionNothing:
		return NewPruningOptions(Nothing)

	case PruningOptionDefault:
		return  NewPruningOptions(Default)

	default:
		return NewPruningOptions(Default)
	}
}