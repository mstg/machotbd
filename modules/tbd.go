/*
* @Author: mustafa
* @Date:   2016-03-29 18:46:24
* @Last Modified by:   mstg
* @Last Modified time: 2016-03-30 00:09:18
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
}

type Tbd_list struct {
  Install_name string
  Archs []Arch
}

func Tbd_form(list Tbd_list) (bytes.Buffer) {
  var buffer bytes.Buffer
  buffer.WriteString("---\n")
  buffer.WriteString("archs: [ ")

  for i, v := range list.Archs {
    buffer.WriteString(v.Name)
    if len(list.Archs)-1 != i {
      buffer.WriteString(",")
    } else {
      buffer.WriteString(" ]\n")
    }
  }

  buffer.WriteString("platform: ios\n")
  buffer.WriteString(fmt.Sprintf("install-name: /System/Library/Frameworks/%s.framework/Versions/A/%s\n", list.Install_name, list.Install_name))
  buffer.WriteString("current-version: 275.0\n")
  buffer.WriteString("exports:\n")

  for _, v := range list.Archs {
    buffer.WriteString(fmt.Sprintf("  - archs: [ %s ]\n", v.Name))

    if len(v.Symbols) > 0 {
      buffer.WriteString("    symbols: [ ")
      for a, b := range v.Symbols {
        buffer.WriteString(b)
        if len(v.Symbols)-1 != a {
          buffer.WriteString(",")
        } else {
          buffer.WriteString(" ]\n")
        }
      }
    }

    if len(v.Classes) > 0 {
      buffer.WriteString("    objc-classes: [ ")
      for a, b := range v.Classes {
        buffer.WriteString(b)
        if len(v.Classes)-1 != a {
          buffer.WriteString(",")
        } else {
          buffer.WriteString(" ]\n")
        }
      }
    }

    if len(v.Ivars) > 0 {
      buffer.WriteString("    objc-ivars: [ ")
      for a, b := range v.Ivars {
        buffer.WriteString(b)
        if len(v.Ivars)-1 != a {
          buffer.WriteString(",")
        } else {
          buffer.WriteString(" ]\n")
        }
      }
    }
  }

  buffer.WriteString("...")

  return buffer
}