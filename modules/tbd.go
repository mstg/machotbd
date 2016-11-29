/*
* @Author: mustafa
* @Date:   2016-03-29 18:46:24
* @Last Modified by:   mstg
* @Last Modified time: 2016-11-29 11:29:52
 */

package tbd

import (
	"bytes"
	"fmt"
	"sort"
)

// Arch is the symbol container
type Arch struct {
	Name      string
	Symbols   []string
	Classes   []string
	Ivars     []string
	Weak      []string
	ReExports []string
}

// List is the list information container
type List struct {
	InstallName string
	Version     string
	CompVersion string
	Platform    string
	Archs       []Arch
}

// tbdSection is the section arch container
type tbdSection struct {
	archN []string
	arch  Arch
}

type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

type sectionSorter struct {
	sections []tbdSection
	by       func(p1, p2 *tbdSection) bool
}

type byFunc func(p1, p2 *tbdSection) bool

func (by byFunc) Sort(sections []tbdSection) {
	ss := &sectionSorter{sections: sections, by: by}
	sort.Sort(ss)
}

func (s *sectionSorter) Len() int {
	return len(s.sections)
}

func (s *sectionSorter) Swap(i, j int) {
	s.sections[i], s.sections[j] = s.sections[j], s.sections[i]
}

func (s *sectionSorter) Less(i, j int) bool {
	return s.by(&s.sections[i], &s.sections[j])
}

func acontains(s []string, a string) (bool, int) {
	for i, v := range s {
		if v == a {
			return true, i
		}
	}

	return false, 0
}

func writeSection(buffer *bytes.Buffer, sect []string, sectDef string) {
	if len(sect) > 0 {
		buffer.WriteString(sectDef)
		amount := 0
		for a, b := range sect {
			amount++

			if amount >= 2 {
				buffer.WriteString(fmt.Sprintf("                        %s", b))
				amount = 0
			} else {
				buffer.WriteString(b)
			}

			if len(sect)-1 != a {
				if amount == 1 {
					buffer.WriteString(",\n")
				} else {
					buffer.WriteString(", ")
				}
			} else {
				buffer.WriteString(" ]\n")
			}
		}
	}
}

func remDepSym(list List, i int, section *tbdSection) {
	for _, k := range list.Archs[i].Symbols {
		for b, l := range list.Archs {
			cont, _a := acontains(l.Symbols, k)
			if cont {
				_L, _ := acontains(section.arch.Symbols, k)
				if !_L {
					section.arch.Symbols = append(section.arch.Symbols, k)
				}

				list.Archs[b].Symbols = append(list.Archs[b].Symbols[:_a], list.Archs[b].Symbols[_a+1:]...)

				cont2, _ := acontains(section.archN, l.Name)
				if !cont2 {
					section.archN = append(section.archN, l.Name)
				}
			}
		}
	}
}

func remDepClass(list List, i int, section *tbdSection) {
	for _, k := range list.Archs[i].Classes {
		for b, l := range list.Archs {
			cont, _a := acontains(l.Classes, k)
			cont2, _ := acontains(section.archN, l.Name)
			if cont {
				_L, _ := acontains(section.arch.Classes, k)
				if !_L {
					section.arch.Classes = append(section.arch.Classes, k)
				}

				list.Archs[b].Classes = append(list.Archs[b].Classes[:_a], list.Archs[b].Classes[_a+1:]...)

				if !cont2 {
					section.archN = append(section.archN, l.Name)
				}
			}
		}
	}
}

func remDepRe(list List, i int, section *tbdSection) {
	for _, k := range list.Archs[i].ReExports {
		for b, l := range list.Archs {
			cont, _a := acontains(l.ReExports, k)
			cont2, _ := acontains(section.archN, l.Name)
			if cont {
				_L, _ := acontains(section.arch.ReExports, k)
				if !_L {
					section.arch.ReExports = append(section.arch.ReExports, k)
				}

				list.Archs[b].ReExports = append(list.Archs[b].ReExports[:_a], list.Archs[b].ReExports[_a+1:]...)

				if !cont2 {
					section.archN = append(section.archN, l.Name)
				}
			}
		}
	}
}

func remDepWeak(list List, i int, section *tbdSection) {
	for _, k := range list.Archs[i].Weak {
		for b, l := range list.Archs {
			cont, _a := acontains(l.Weak, k)
			cont2, _ := acontains(section.archN, l.Name)
			if cont {
				_L, _ := acontains(section.arch.Weak, k)
				if !_L {
					section.arch.Weak = append(section.arch.Weak, k)
				}

				list.Archs[b].Weak = append(list.Archs[b].Weak[:_a], list.Archs[b].Weak[_a+1:]...)

				if !cont2 {
					section.archN = append(section.archN, l.Name)
				}
			}
		}
	}
}

func remDepIvar(list List, i int, section *tbdSection) {
	for _, k := range list.Archs[i].Ivars {
		for b, l := range list.Archs {
			cont, _a := acontains(l.Ivars, k)
			cont2, _ := acontains(section.archN, l.Name)
			if cont {
				_L, _ := acontains(section.arch.Ivars, k)
				if !_L {
					section.arch.Ivars = append(section.arch.Ivars, k)
				}

				list.Archs[b].Ivars = append(list.Archs[b].Ivars[:_a], list.Archs[b].Ivars[_a+1:]...)

				if !cont2 {
					section.archN = append(section.archN, l.Name)
				}
			}
		}
	}
}

// Form Generates the tbd output
func Form(list List) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("---\n")
	buffer.WriteString("archs:           [ ")

	var archs []string

	for _, v := range list.Archs {
		archs = append(archs, v.Name)
	}

	sort.Strings(archs)

	for i, v := range archs {
		buffer.WriteString(v)
		if len(list.Archs)-1 != i {
			buffer.WriteString(", ")
		} else {
			buffer.WriteString(" ]\n")
		}
	}

	buffer.WriteString(fmt.Sprintf("platform:        %s\n", list.Platform))
	buffer.WriteString(fmt.Sprintf("install-name:    %s\n", list.InstallName))

	buffer.WriteString(fmt.Sprintf("current-version: %s\n", list.Version))
	buffer.WriteString(fmt.Sprintf("compatibility-version: %s\n", list.CompVersion))

	buffer.WriteString("exports:\n")

	var tbdSections []tbdSection

	for len(list.Archs) > 0 {
		for i := 0; i < len(list.Archs); i++ {
			var section tbdSection
			section.arch = Arch{}

			remDepSym(list, i, &section)
			remDepClass(list, i, &section)
			remDepRe(list, i, &section)
			remDepWeak(list, i, &section)
			remDepIvar(list, i, &section)

			list.Archs = append(list.Archs[:i], list.Archs[i+1:]...)

			if len(section.arch.Symbols) > 0 || len(section.arch.Classes) > 0 ||
				len(section.arch.ReExports) > 0 || len(section.arch.Weak) > 0 ||
				len(section.arch.Ivars) > 0 {

				sort.Sort(byLength(section.arch.ReExports))
				sort.Strings(section.arch.Weak)
				sort.Strings(section.arch.Symbols)
				sort.Strings(section.arch.Classes)
				sort.Strings(section.arch.Ivars)
				sort.Strings(section.archN)
				tbdSections = append(tbdSections, section)
			}
		}
	}

	count := func(s1, s2 *tbdSection) bool {
		return len(s1.archN) < len(s2.archN)
	}

	byFunc(count).Sort(tbdSections)
	for _, v := range tbdSections {
		buffer.WriteString("  - archs:            [ ")

		for b, a := range v.archN {
			if len(v.archN)-1 != b {
				buffer.WriteString(fmt.Sprintf("%s, ", a))
			} else {
				buffer.WriteString(fmt.Sprintf("%s ]\n", a))
			}
		}

		writeSection(&buffer, v.arch.ReExports, "    re-exports:       [ ")
		writeSection(&buffer, v.arch.Weak, "    weak-def-symbols:  [ ")
		writeSection(&buffer, v.arch.Symbols, "    symbols:          [ ")
		writeSection(&buffer, v.arch.Classes, "    objc-classes:     [ ")
		writeSection(&buffer, v.arch.Ivars, "    objc-ivars:       [ ")
	}

	buffer.WriteString("...\n")

	return buffer
}
