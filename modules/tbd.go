/*
* @Author: mustafa
* @Date:   2016-03-29 18:46:24
* @Last Modified by:   mstg
* @Last Modified time: 2016-03-30 05:23:33
*/

package tbd

import (
  "bytes"
  "fmt"
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

  for _, v := range list.Archs {
    buffer.WriteString(fmt.Sprintf("  - archs:            [ %s ]\n", v.Name))

    if len(v.ReExports) > 0 {
      buffer.WriteString("    re-exports:       [ ")
      amount := 0
      for a, b := range v.ReExports {
        amount++

        if amount >= 2 {
          buffer.WriteString(fmt.Sprintf("                        %s", b))
          amount = 0
        } else {
          buffer.WriteString(b)
        }

        if len(v.ReExports)-1 != a {
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

    if len(v.Weak) > 0 {
      buffer.WriteString("    weak-def-symbols:  [ ")
      amount := 0
      for a, b := range v.Weak {
        amount++

        if amount >= 2 {
          buffer.WriteString(fmt.Sprintf("                        %s", b))
          amount = 0
        } else {
          buffer.WriteString(b)
        }

        if len(v.Weak)-1 != a {
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

    if len(v.Symbols) > 0 {
      buffer.WriteString("    symbols:          [ ")
      amount := 0
      for a, b := range v.Symbols {
        amount++

        if amount >= 2 {
          buffer.WriteString(fmt.Sprintf("                        %s", b))
          amount = 0
        } else {
          buffer.WriteString(b)
        }

        if len(v.Symbols)-1 != a {
          if amount >= 1 {
            buffer.WriteString(",\n")
          } else {
            buffer.WriteString(", ")
          }
        } else {
          buffer.WriteString(" ]\n")
        }
      }
    }

    if len(v.Classes) > 0 {
      buffer.WriteString("    objc-classes:     [ ")
      amount := 0
      for a, b := range v.Classes {
        amount++

        if amount >= 2 {
          buffer.WriteString(fmt.Sprintf("                        %s", b))
          amount = 0
        } else {
          buffer.WriteString(b)
        }

        if len(v.Classes)-1 != a {
          if amount >= 1 {
            buffer.WriteString(",\n")
          } else {
            buffer.WriteString(", ")
          }
        } else {
          buffer.WriteString(" ]\n")
        }
      }
    }

    if len(v.Ivars) > 0 {
      buffer.WriteString("    objc-ivars:       [ ")
      amount := 0
      for a, b := range v.Ivars {
        amount++

        if amount >= 2 {
          buffer.WriteString(fmt.Sprintf("                        %s", b))
          amount = 0
        } else {
          buffer.WriteString(b)
        }

        if len(v.Ivars)-1 != a {
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

  buffer.WriteString("...")

  return buffer
}
