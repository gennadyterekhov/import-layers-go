package layers

import "strings"

const ViolationMessage = "cannot import package from lower layer"

type Layer struct {
	Name string
	rank int
}

func (l Layer) Forbids(target Layer) bool {
	return target.rank > 0 && l.rank > target.rank
}

type Policy struct {
	qualified []Layer
	segments  []Layer
}

func New(names []string) Policy {
	var policy Policy
	for index, name := range names {
		layer := Layer{Name: name, rank: len(names) - index}
		if strings.Contains(name, "/") {
			policy.qualified = append(policy.qualified, layer)
		} else {
			policy.segments = append(policy.segments, layer)
		}
	}
	return policy
}

// Match gives qualified selectors priority over single segments. When a
// directory is unavailable (for example during go vet), compare qualified
// selectors against complete segments of the package import path.
func (p Policy) Match(packagePath, directory string) Layer {
	for _, layer := range p.qualified {
		if within(packagePath, layer.Name) || within(directory, layer.Name) ||
			directory == "" && containsSegments(packagePath, layer.Name) {
			return layer
		}
	}
	for _, layer := range p.segments {
		if containsSegments(packagePath, layer.Name) {
			return layer
		}
	}
	return Layer{}
}

func within(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func containsSegments(path, selector string) bool {
	return strings.Contains("/"+path+"/", "/"+selector+"/")
}
