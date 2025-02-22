package golinters

import (
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	_ "golang.org/x/tools/go/analysis/passes/buildssa" // unused, internal analyzer
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	_ "golang.org/x/tools/go/analysis/passes/ctrlflow" // unused, internal analyzer
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	_ "golang.org/x/tools/go/analysis/passes/inspect" // unused internal analyzer
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	_ "golang.org/x/tools/go/analysis/passes/pkgfact" // unused, internal analyzer
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

var (
	allAnalyzers = []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
	}

	// https://github.com/golang/go/blob/b56645a87b28840a180d64077877cb46570b4176/src/cmd/vet/main.go#L49-L81
	defaultAnalyzers = []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}
)

func NewGovet(settings *config.GovetSettings) *goanalysis.Linter {
	var conf map[string]map[string]any
	if settings != nil {
		conf = settings.Settings
	}

	return goanalysis.NewLinter(
		"govet",
		"Vet examines Go source code and reports suspicious constructs. "+
			"It is roughly the same as 'go vet' and uses its passes.",
		analyzersFromConfig(settings),
		conf,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func analyzersFromConfig(settings *config.GovetSettings) []*analysis.Analyzer {
	if settings == nil {
		return defaultAnalyzers
	}

	if settings.CheckShadowing {
		// Keeping for backward compatibility.
		settings.Enable = append(settings.Enable, shadow.Analyzer.Name)
	}

	var enabledAnalyzers []*analysis.Analyzer
	for _, a := range allAnalyzers {
		if isAnalyzerEnabled(a.Name, settings, defaultAnalyzers) {
			enabledAnalyzers = append(enabledAnalyzers, a)
		}
	}

	return enabledAnalyzers
}

func isAnalyzerEnabled(name string, cfg *config.GovetSettings, defaultAnalyzers []*analysis.Analyzer) bool {
	// TODO(ldez) remove loopclosure when go1.23
	if name == loopclosure.Analyzer.Name && config.IsGoGreaterThanOrEqual(cfg.Go, "1.22") {
		return false
	}

	switch {
	case cfg.EnableAll:
		return !slices.Contains(cfg.Disable, name)

	case slices.Contains(cfg.Enable, name):
		return true

	case slices.Contains(cfg.Disable, name):
		return false

	case cfg.DisableAll:
		return false

	default:
		return slices.ContainsFunc(defaultAnalyzers, func(a *analysis.Analyzer) bool { return a.Name == name })
	}
}
