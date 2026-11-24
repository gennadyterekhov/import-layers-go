package analyzer

import (
	"go/parser"
	"go/token"
	"maps"
	"path"
	"slices"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/filechecker"
	"github.com/gennadyterekhov/import-layers-go/internal/project"
	"golang.org/x/tools/go/packages"
)

type Request struct {
	Directory  string
	ConfigPath string
	Paths      []string
	OnFile     func(path string)
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Graph struct {
	Nodes []string `json:"nodes"`
	Edges []Edge   `json:"edges"`
}

type Violation struct {
	File      string
	Line      int
	Import    string
	FromLayer string
	ToLayer   string
}

type Result struct {
	Graph      Graph
	Errors     []error
	Violations []Violation
}

type sourceFile struct {
	filename    string
	packagePath string
	imports     map[string]*packages.Package
}

func CheckProject(request Request) (Result, error) {
	cfg, err := config.Read(request.Directory, request.ConfigPath)
	if err != nil {
		return Result{}, err
	}
	selection, err := cfg.Root.SelectPackages(request.Directory, cfg.RootDir, request.Paths)
	if err != nil {
		return Result{}, err
	}
	roots, err := packages.Load(&packages.Config{
		Dir:   string(cfg.Root),
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps,
		Tests: !cfg.IgnoreTests,
	}, selection.Patterns...)
	if err != nil {
		return Result{}, err
	}

	files, issues := collectFiles(roots, cfg, selection)
	nodes := slices.Sorted(maps.Keys(files))
	if nodes == nil {
		nodes = []string{}
	}
	result := Result{
		Graph:  Graph{Nodes: nodes, Edges: []Edge{}},
		Errors: issues,
	}
	fset := token.NewFileSet()
	for _, relative := range result.Graph.Nodes {
		if cfg.Debug && request.OnFile != nil {
			request.OnFile(relative)
		}
		source := files[relative]
		parsed, err := parser.ParseFile(fset, source.filename, nil, parser.SkipObjectResolution)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}

		imports, errs := filechecker.Imports(parsed)
		result.Errors = append(result.Errors, errs...)
		fromLayer := cfg.Layers.Match(source.packagePath, path.Dir(relative))
		targetFiles := make(map[string]bool)
		for _, ref := range imports {
			target := source.imports[ref.Path]
			if target == nil {
				continue
			}
			targetDirectory := ""
			for _, filename := range target.GoFiles {
				to, ok := cfg.Root.RelativePath(filename)
				if !ok {
					continue
				}
				targetDirectory = path.Dir(to)
				if _, included := files[to]; included {
					targetFiles[to] = true
				}
			}
			if fromLayer.Name == "" {
				continue
			}
			toLayer := cfg.Layers.Match(ref.Path, targetDirectory)
			if fromLayer.Forbids(toLayer) {
				position := fset.PositionFor(ref.Position, false)
				result.Violations = append(result.Violations, Violation{
					File: relative, Line: position.Line, Import: ref.Path,
					FromLayer: fromLayer.Name, ToLayer: toLayer.Name,
				})
			}
		}
		for _, to := range slices.Sorted(maps.Keys(targetFiles)) {
			result.Graph.Edges = append(result.Graph.Edges, Edge{From: relative, To: to})
		}
	}
	return result, nil
}

func collectFiles(roots []*packages.Package, cfg *config.Config, selection project.PackageSelection) (map[string]sourceFile, []error) {
	files := make(map[string]sourceFile)
	issues := make(map[string]error)
	packages.Visit(roots, func(pkg *packages.Package) bool {
		for _, issue := range pkg.Errors {
			issues[issue.Error()] = issue
		}
		for _, filename := range pkg.GoFiles {
			relative, ok := cfg.Root.RelativePath(filename)
			if !ok || !selection.IncludesFile(relative) || cfg.IgnoreTests && project.IsTestFile(filename) {
				continue
			}

			source, found := files[relative]
			if !found {
				source = sourceFile{
					filename:    filename,
					packagePath: pkg.PkgPath,
					imports:     make(map[string]*packages.Package, len(pkg.Imports)),
				}
				files[relative] = source
			}
			// Test variants may add imports to a file's package.
			maps.Copy(source.imports, pkg.Imports)
		}
		return true
	}, nil)

	var orderedIssues []error
	for _, message := range slices.Sorted(maps.Keys(issues)) {
		orderedIssues = append(orderedIssues, issues[message])
	}
	return files, orderedIssues
}
