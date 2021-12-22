package task

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"go.dagger.io/dagger/compiler"
	"go.dagger.io/dagger/plancontext"
	"go.dagger.io/dagger/solver"
)

func init() {
	Register("Merge", func() Task { return &mergeTask{} })
}

type mergeTask struct {
}

type copyInfo struct {
	Source struct {
		root llb.State
		path string
	}
	Dest string
}

func (t *mergeTask) Run(ctx context.Context, pctx *plancontext.Context, s solver.Solver, v *compiler.Value) (*compiler.Value, error) {
	var err error

	// input, err := pctx.FS.FromValue(v.Lookup("input"))
	// if err != nil {
	// 	return nil, err
	// }

	// inputState, err := input.State()
	// if err != nil {
	// 	return nil, err
	// }

	cueLayers, err := v.Lookup("layers").List()
	if err != nil {
		return nil, err
	}

	var copyInfos []copyInfo

	for _, layer := range cueLayers {
		copy := copyInfo{}
		sourceRoot, err := pctx.FS.FromValue(v.Lookup("source.root"))
		if err != nil {
			return nil, err
		}
		copy.Source.root, err = sourceRoot.State()
		if err != nil {
			return nil, err
		}
		copy.Source.path, err = layer.Lookup("source.path").String()
		if err != nil {
			return nil, err
		}
		copy.Dest, err = layer.Lookup("dest").String()
		if err != nil {
			return nil, err
		}
		copyInfos = append(copyInfos, copy)
	}

	outputState := llb.Merge(copyInfos)

	outputState := inputState.File(
		llb.Copy(
			sourceState,
			sourcePath,
			destPath,
			// FIXME: allow more configurable llb options
			// For now we define the following convenience presets:
			&llb.CopyInfo{
				CopyDirContentsOnly: true,
				CreateDestPath:      true,
				AllowWildcard:       true,
			},
		),
		withCustomName(v, "Copy %s %s", sourcePath, destPath),
	)

	result, err := s.Solve(ctx, outputState, pctx.Platform.Get())
	if err != nil {
		return nil, err
	}

	fs := pctx.FS.New(result)

	return compiler.NewValue().FillFields(map[string]interface{}{
		"output": fs.MarshalCUE(),
	})
}

func extractCopyInfo(ctx context.Context, pctx *plancontext.Context, v *compiler.Value) (*copyInfo, error) {
	var err error

	var source struct {
		root llb.State
		path string
	}
	if err = v.Decode(&source); err != nil {
		return nil, err
	}

	dest, err := v.Lookup("dest").String()
	if err != nil {
		return nil, err
	}

	return &copyInfo{
		Source: source,
		Dest:   dest,
	}, nil
}
