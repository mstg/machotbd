/*
* @Author: mustafa
* @Date:   2016-03-29 17:31:09
* @Last Modified by:   mstg
* @Last Modified time: 2016-11-29 11:47:52
 */

package main

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"github.com/mstg/machotbd/modules"
	"log"
	"os"
	"strings"
)

const (
	arm64 macho.Cpu = 0x100000C

	// NType From <mach-o/nlist.h>
	NType uint8 = 0x0e

	// NSect From <mach-o/nlist.h>
	NSect uint8 = 0xe

	// NExt From <mach-o/nlist.h>
	NExt uint8 = 0x01

	// NWeakRef From <mach-o/nlist.h>
	NWeakRef uint16 = 0x0040

	// LoadDylibIDCmd Load command in dylib
	LoadDylibIDCmd = 0xd

	// fileHeaderSize32 File header size of 32-bit binaries
	fileHeaderSize32 = 7 * 4

	// fileHeaderSize64 File header size of 64-bit binaries
	fileHeaderSize64 = 8 * 4

	// ReExportDylibCmd Reexports command in dylib
	ReExportDylibCmd = (0x1f | 0x80000000)
)

// DylibIDCmd Header information
type DylibIDCmd struct {
	Cmd            macho.LoadCmd
	Len            uint32
	Name           uint32
	Time           uint32
	CurrentVersion uint32
	CompatVersion  uint32
}

func cstring(b []byte) string {
	var i int
	for i = 0; i < len(b) && b[i] != 0; i++ {
	}
	return string(b[0:i])
}

func ver(rawVer uint32) string {
	return fmt.Sprintf("%d.%d.%d", rawVer>>16, (rawVer>>8)&0xff, rawVer&0xff)
}

func magicType(magic uint32) uint32 {
	if magic == macho.Magic32 {
		return 32
	} else if magic == macho.Magic64 {
		return 64
	} else if magic == macho.MagicFat {
		return 1
	}

	return 0
}

func cpuType(f *macho.File) string {
	if f.Cpu == macho.Cpu386 {
		return "i386"
	} else if f.Cpu == macho.CpuAmd64 {
		return "x86_64"
	} else if f.Cpu == macho.CpuArm && f.SubCpu == 6 {
		return "armv6"
	} else if f.Cpu == macho.CpuArm && f.SubCpu == 9 {
		return "armv7"
	} else if f.Cpu == macho.CpuArm && f.SubCpu == 11 {
		return "armv7s"
	} else if f.Cpu == arm64 {
		return "arm64"
	}

	return "uns"
}

func cmdScan(f *macho.File) (string, string, string, []string) {
	version := "0.0.0"
	compatibilityVersion := "0.0.0"
	path := ""
	realReexports := []string{}

	bo := f.ByteOrder
	for _, v := range f.Loads {
		dat := v.Raw()
		cmd, siz := uint32(bo.Uint32(dat[0:4])), bo.Uint32(dat[4:8])
		var cmddat []byte
		cmddat, _ = dat[0:siz], dat[siz:]

		switch cmd {
		case LoadDylibIDCmd:
			var hdr DylibIDCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				break
			}
			path = cstring(cmddat[hdr.Name:])
			version = ver(hdr.CurrentVersion)
			compatibilityVersion = ver(hdr.CompatVersion)
			break
		case ReExportDylibCmd:
			var hdr DylibIDCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				break
			}
			path := cstring(cmddat[hdr.Name:])
			realReexports = append(realReexports, path)
			break
		}
	}

	return version, compatibilityVersion, path, realReexports
}

func appropriatePlacement(v macho.Symbol, cput string) ([]string, []string, []string, []string) {
	realSymbols := []string{}
	realClasses := []string{}
	realIvars := []string{}
	realWeak := []string{}

	if strings.Contains(v.Name, "$ld$") {
		realName := fmt.Sprintf("'%s'", v.Name)
		if strings.Contains(v.Name, "_OBJC_CLASS") {
			realClasses = append(realClasses, realName)
		} else if strings.Contains(v.Name, "_OBJC_IVAR") {
			realIvars = append(realIvars, realName)
		} else if strings.Contains(v.Name, "_OBJC_METACLASS") {
		} else if cput == "i386" && strings.Contains(v.Name, ".objc_class_name") {
			realClasses = append(realClasses, realName)
		} else if v.Desc&NWeakRef == NWeakRef {
			realWeak = append(realWeak, realName)
		} else {
			realSymbols = append(realSymbols, realName)
		}
	} else if strings.Contains(v.Name, "_OBJC_CLASS") {
		realName := strings.Replace(v.Name, "_OBJC_CLASS_$", "", -1)
		realClasses = append(realClasses, realName)
	} else if strings.Contains(v.Name, "_OBJC_IVAR") {
		realName := strings.Replace(v.Name, "_OBJC_IVAR_$", "", -1)
		realIvars = append(realIvars, realName)
	} else if strings.Contains(v.Name, "_OBJC_METACLASS") {
	} else if cput == "i386" && strings.Contains(v.Name, ".objc_class_name") {
		realName := strings.Replace(v.Name, ".objc_class_name", "", -1)
		realClasses = append(realClasses, realName)
	} else if v.Desc&NWeakRef == NWeakRef {
		realWeak = append(realWeak, v.Name)
	} else {
		realSymbols = append(realSymbols, v.Name)
	}

	return realSymbols, realClasses, realIvars, realWeak
}

func parseMacho(f *macho.File, stdout *log.Logger, stderr *log.Logger) (tbd.Arch, []string, error) {
	mt := magicType(f.Magic)
	cput := cpuType(f)

	var _syms tbd.Arch

	if cput == "uns" {
		return _syms, []string{}, errors.New("Unsupported arch")
	}

	stdout.Println(mt, "bit", cput, "slice")

	symtab := f.Symtab
	realSymbols := []string{}
	realClasses := []string{}
	realIvars := []string{}
	realWeak := []string{}
	for _, v := range symtab.Syms {
		if v.Type&NType == NSect && v.Type&NExt == NExt {
			if v.Name != "" {
				realS, realC, realI, realW := appropriatePlacement(v, cput)
				realSymbols = append(realSymbols, realS...)
				realClasses = append(realClasses, realC...)
				realIvars = append(realIvars, realI...)
				realWeak = append(realWeak, realW...)
			}
		}
	}

	version, compatibilityVersion, path, realReexports := cmdScan(f)

	_syms = tbd.Arch{Name: cput, Symbols: realSymbols, Classes: realClasses, Ivars: realIvars, Weak: realWeak, ReExports: realReexports}
	return _syms, []string{version, path, compatibilityVersion}, nil
}

func parseFat(f *macho.FatFile, stdout *log.Logger, stderr *log.Logger) tbd.List {
	stdout.Println("Universal Mach-O")

	retSym := tbd.List{}
	for _, v := range f.Arches {
		retMachoSym, info, err := parseMacho(v.File, stdout, stderr)
		if err == nil {
			retSym.Archs = append(retSym.Archs, retMachoSym)
			retSym.InstallName = info[1]
			retSym.Version = info[0]
			retSym.CompVersion = info[2]
		}
	}

	return retSym
}

func printOrWrite(stdout *log.Logger, stderr *log.Logger, _list tbd.List) {
	_buf := tbd.Form(_list)
	printit := 0
	if *print == true && *out == "" {
		stdout.Println(_buf.String())
	} else if *out != "" {
		_, err := os.Stat(*out)

		if os.IsNotExist(err) {
			var file, err = os.Create(*out)
			if err != nil {
				printit = 1
			} else {
				defer file.Close()
			}
		}

		file, err := os.OpenFile(*out, os.O_RDWR, 0644)

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
		stderr.Println("An error occurred during I/O, printing to stdout")
		stdout.Println(_buf.String())
	} else if *out != "" {
		stdout.Println("Wrote to", *out)
	}
}

var out = flag.String("out", "", "path to export tbd to")
var print = flag.Bool("print", true, "print tbd to stdout")
var plt = flag.String("platform", "ios", "platform to define in the output tbd")

func machoTbd(args []string) {
	stderr := log.New(os.Stderr, "[?] ", 0)
	stdout := log.New(os.Stdout, "[+] ", 0)
	file := ""
	if len(args) > 0 {
		file = args[0]
	} else {
		stderr.Println("No Mach-O file provided")
		os.Exit(1)
	}

	if *out != "" {
		*print = false
	}

	if *plt != "ios" && *plt != "macosx" && *plt != "watchos" && *plt != "tvos" {
		stderr.Println("Unsupported platform, only ios, macosx, watchos and tvos is supported")
		os.Exit(1)
	}

	machoFile, err := macho.Open(file)
	var machoFatFile *macho.FatFile
	universal := false

	if err != nil {
		machoFatFile, err = macho.OpenFat(file)
	}

	if err != nil {
		stderr.Println("Malformed or invalid Mach-O provided, err:", err)
		os.Exit(1)
	}

	if machoFatFile != nil {
		universal = true
	}

	var _list tbd.List
	if universal {
		_list = parseFat(machoFatFile, stdout, stderr)

		stdout.Println("Arch count:", len(_list.Archs))
	} else {
		_unpreplist, info, err := parseMacho(machoFile, stdout, stderr)
		if err == nil {
			archArr := []tbd.Arch{_unpreplist}
			_list = tbd.List{Archs: archArr}
			_list.InstallName = info[1]
			_list.Version = info[0]
			_list.CompVersion = info[2]
		}
	}

	_list.Platform = *plt

	printOrWrite(stdout, stderr, _list)
}

func main() {
	flag.Parse()
	machoTbd(flag.Args())
}
