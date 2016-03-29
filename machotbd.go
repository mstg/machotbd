/*
* @Author: mustafa
* @Date:   2016-03-29 17:31:09
* @Last Modified by:   mstg
* @Last Modified time: 2016-03-30 00:07:24
*/

package main

import (
  "github.com/codegangsta/cli"
  "os"
  "log"
  "debug/macho"
  "github.com/mstg/machotbd/modules"
  "errors"
  "strings"
)

const (
  arm64 macho.Cpu = 16777228

  // From <mach-o/nlist.h>
  N_TYPE uint8 = 0x0e
  N_SECT uint8 = 0xe
)

func magic_type(magic uint32) (uint32) {
  if magic == macho.Magic32 {
    return 32
  } else if magic == macho.Magic64 {
    return 64
  } else if magic == macho.MagicFat {
    return 1
  }

  return 0
}

func cpu_type(cpu macho.Cpu) (string) {
  if cpu == macho.CpuAmd64 {
    return "x86_64"
  } else if cpu == macho.CpuArm {
    return "armv7"
  } else if cpu == arm64 {
    return "arm64"
  }

  return "uns"
}

func parse_macho(f *macho.File, stdout *log.Logger, stderr *log.Logger) (tbd.Arch, error) {
  mt := magic_type(f.Magic)
  cput := cpu_type(f.Cpu)

  var _syms tbd.Arch

  if cput == "uns" {
    return _syms, errors.New("Unsupported arch")
  }

  stdout.Println(mt, "bit", cput, "slice")

  symtab := f.Symtab
  real_symbols := []string{}
  real_classes := []string{}
  real_ivars := []string{}
  for _, v := range symtab.Syms {
    if v.Type & N_TYPE == N_SECT {
      if v.Name != "" {
        if strings.Contains(v.Name, "_OBJC_CLASS") {
          real_classes = append(real_classes, v.Name)
        } else if strings.Contains(v.Name, "_OBJC_IVAR") {
          real_ivars = append(real_ivars, v.Name)
        } else {
          real_symbols = append(real_symbols, v.Name)
        }
      }
    }
  }

  _syms = tbd.Arch{Name: cput, Symbols: real_symbols, Classes: real_classes, Ivars: real_ivars}
  return _syms, nil
}

func parse_fat(f *macho.FatFile, stdout *log.Logger, stderr *log.Logger) (tbd.Tbd_list) {
  stdout.Println("Universal Mach-O")

  _ret_sym := tbd.Tbd_list{}
  for _, v := range f.Arches {
    _ret_macho_sym, err := parse_macho(v.File, stdout, stderr)
    if err == nil {
      _ret_sym.Archs = append(_ret_sym.Archs, _ret_macho_sym)
    }
  }

  return _ret_sym
}

func macho_tbd(c *cli.Context) {
  stderr := log.New(os.Stderr, "[?] ", 0)
  stdout := log.New(os.Stdout, "[+] ", 0)
  file := ""
  if c.NArg() > 0 {
    file = c.Args()[0]
  } else {
    stderr.Println("No Mach-O file provided")
    os.Exit(1)
  }

  macho_file, err := macho.Open(file)
  var macho_fat_file *macho.FatFile
  universal := false

  if err != nil {
    macho_fat_file, err = macho.OpenFat(file)
  }

  if err != nil {
    stderr.Println("Malformed or invalid Mach-O provided, err:", err)
    os.Exit(1)
  }

  if macho_fat_file != nil {
    universal = true
  }

  var _list tbd.Tbd_list
  if universal {
    _list = parse_fat(macho_fat_file, stdout, stderr)

    stdout.Println(len(_list.Archs))
  } else {
    _unpreplist, err := parse_macho(macho_file, stdout, stderr)
    if err == nil {
      arch_arr := []tbd.Arch{_unpreplist}
      _list = tbd.Tbd_list{Archs: arch_arr}
    }
  }

  _list.Install_name = file

  armv7s := false
  for i, v := range _list.Archs {
    if v.Name == "armv7" && !armv7s {
      _list.Archs[i].Name = "armv7s"
      armv7s = true
    }
  }

  _buf := tbd.Tbd_form(_list)

  printit := 0
  if c.Int("print") == 1 && c.String("o") == ""{
    println(_buf.String())
  } else if c.String("o") != "" {
    _, err := os.Stat(c.String("o"))

    if os.IsNotExist(err) {
      var file, err = os.Create(c.String("o"))
      if err != nil {
        printit = 1
      } else {
        defer file.Close()
      }
    }

    file, err := os.OpenFile(c.String("o"), os.O_RDWR, 0644)

    if err != nil {
      printit = 1
    } else {
      defer file.Close()
    }

    _, err = file.WriteString(_buf.String())

    if err != nil {
      printit = 1
    } else {
      err = file.Sync()
      if err != nil {
        printit = 1
      }
    }
  }

  if printit == 1 {
    stderr.Println("An error occured during I/O, printing to stdout")
    println(_buf.String())
  } else {
    stdout.Println("Wrote to", c.String("o"))
  }
}

func main() {
  app := cli.NewApp()
  app.Name = "machotbd"
  app.Usage = "dump mach-o symbols to a tbd file"
  app.Flags = []cli.Flag {
    cli.IntFlag{
      Name: "print",
      Value: 1,
      Usage: "print symbols to stdout",
    },
    cli.StringFlag{
      Name: "o",
      Value: "",
      Usage: "path to the file should be exported to",
    },
  }
  app.Action = macho_tbd

  app.Run(os.Args)
}