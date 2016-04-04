/*
* @Author: mustafa
* @Date:   2016-03-29 18:46:24
* @Last Modified by:   Mustafa
* @Last Modified time: 2016-04-04 05:01:50
*/

package tbd

import (
  "bytes"
  "fmt"
  "sort"
)

type Arch struct {
  Name string
  Symbols []string
  Classes []string
  Ivars []string
  Weak []string
  ReExports []string
}

type Tbd_list struct {
  Install_name string
  Version string
  CompVersion string
  Platform string
  Archs []Arch
}

type tbd_section struct {
  arch_n []string
  arch Arch
}

type ByLength []string

func (s ByLength) Len() int {
  return len(s)
}
func (s ByLength) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
  return len(s[i]) > len(s[j])
}

type sectionSorter struct {
  sections []tbd_section
  by func(p1, p2 *tbd_section) bool
}

type By func(p1, p2 *tbd_section) bool

func (by By) Sort(sections []tbd_section) {
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

func write_section (buffer *bytes.Buffer, sect []string, sect_def string) {
  if len(sect) > 0 {
    buffer.WriteString(sect_def)
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

func Tbd_form(list Tbd_list) (bytes.Buffer) {
  var buffer bytes.Buffer
  buffer.WriteString("---\n")
  buffer.WriteString("archs:           [ ")

  for i, v := range list.Archs {
    buffer.WriteString(v.Name)
    if len(list.Archs)-1 != i {
      buffer.WriteString(", ")
    } else {
      buffer.WriteString(" ]\n")
    }
  }

  buffer.WriteString(fmt.Sprintf("platform:        %s\n", list.Platform))
  buffer.WriteString(fmt.Sprintf("install-name:    %s\n", list.Install_name))

  buffer.WriteString(fmt.Sprintf("current-version: %s\n", list.Version))
  if list.CompVersion != "" {
    buffer.WriteString(fmt.Sprintf("compatibility-version: %s\n", list.CompVersion))
  }

  buffer.WriteString("exports:\n")

  var tbd_sections []tbd_section

  for len(list.Archs) > 0 {
    for i := 0; i < len(list.Archs); i++ {
      var section tbd_section
      section.arch = Arch{}

      for _, k := range list.Archs[i].Symbols {
        for b, l := range list.Archs {
          cont, _a := acontains(l.Symbols, k)
          if cont {
            __l, _ := acontains(section.arch.Symbols, k)
            if !__l {
              section.arch.Symbols = append(section.arch.Symbols, k)
            }

            list.Archs[b].Symbols = append(list.Archs[b].Symbols[:_a], list.Archs[b].Symbols[_a+1:]...)

            cont2, _ := acontains(section.arch_n, l.Name)
            if !cont2 {
              section.arch_n = append(section.arch_n, l.Name)
            }
          }
        }
      }

      for _, k := range list.Archs[i].Classes {
        for b, l := range list.Archs {
          cont, _a := acontains(l.Classes, k)
          cont2, _ := acontains(section.arch_n, l.Name)
          if cont {
            __l, _ := acontains(section.arch.Classes, k)
            if !__l {
              section.arch.Classes = append(section.arch.Classes, k)
            }

            list.Archs[b].Classes = append(list.Archs[b].Classes[:_a], list.Archs[b].Classes[_a+1:]...)

            if !cont2 {
              section.arch_n = append(section.arch_n, l.Name)
            }
          }
        }
      }

      for _, k := range list.Archs[i].ReExports {
        for b, l := range list.Archs {
          cont, _a := acontains(l.ReExports, k)
          cont2, _ := acontains(section.arch_n, l.Name)
          if cont {
            __l, _ := acontains(section.arch.ReExports, k)
            if !__l {
              section.arch.ReExports = append(section.arch.ReExports, k)
            }

            list.Archs[b].ReExports = append(list.Archs[b].ReExports[:_a], list.Archs[b].ReExports[_a+1:]...)

            if !cont2 {
              section.arch_n = append(section.arch_n, l.Name)
            }
          }
        }

        for _, k := range list.Archs[i].Weak {
          for b, l := range list.Archs {
            cont, _a := acontains(l.Weak, k)
            cont2, _ := acontains(section.arch_n, l.Name)
            if cont {
              __l, _ := acontains(section.arch.Weak, k)
              if !__l {
                section.arch.Weak = append(section.arch.Weak, k)
              }

              list.Archs[b].Weak = append(list.Archs[b].Weak[:_a], list.Archs[b].Weak[_a+1:]...)

              if !cont2 {
                section.arch_n = append(section.arch_n, l.Name)
              }
            }
          }
        }

        for _, k := range list.Archs[i].Ivars {
          for b, l := range list.Archs {
            cont, _a := acontains(l.Ivars, k)
            cont2, _ := acontains(section.arch_n, l.Name)
            if cont {
              __l, _ := acontains(section.arch.Ivars, k)
              if !__l {
                section.arch.Ivars = append(section.arch.Ivars, k)
              }

              list.Archs[b].Ivars = append(list.Archs[b].Ivars[:_a], list.Archs[b].Ivars[_a+1:]...)

              if !cont2 {
                section.arch_n = append(section.arch_n, l.Name)
              }
            }
          }
        }
      }

      list.Archs = append(list.Archs[:i], list.Archs[i+1:]...)

      if len(section.arch.Symbols) > 0 || len(section.arch.Classes) > 0 ||
        len(section.arch.ReExports) > 0 || len(section.arch.Weak) > 0 ||
        len(section.arch.Ivars) > 0 {

        sort.Sort(ByLength(section.arch.ReExports))
        sort.Strings(section.arch.Weak)
        sort.Strings(section.arch.Symbols)
        sort.Strings(section.arch.Classes)
        sort.Strings(section.arch.Ivars)
        sort.Strings(section.arch_n)
        tbd_sections = append(tbd_sections, section)
      }
    }
  }

  count := func(s1, s2 *tbd_section) (bool) {
    return len(s1.arch_n) < len(s2.arch_n)
  }

  By(count).Sort(tbd_sections)
  for _, v := range tbd_sections {
    buffer.WriteString("  - archs:            [ ")

    for b, a := range v.arch_n {
      if len(v.arch_n)-1 != b {
        buffer.WriteString(fmt.Sprintf("%s, ", a))
      } else {
        buffer.WriteString(fmt.Sprintf("%s ]\n", a))
      }
    }

    write_section(&buffer, v.arch.ReExports, "    re-exports:       [ ")
    write_section(&buffer, v.arch.Weak, "    weak-def-symbols:  [ ")
    write_section(&buffer, v.arch.Symbols, "    symbols:          [ ")
    write_section(&buffer, v.arch.Classes, "    objc-classes:     [ ")
    write_section(&buffer, v.arch.Ivars, "    objc-ivars:       [ ")
  }

  buffer.WriteString("...\n")

  return buffer
}
