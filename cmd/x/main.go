// Copyright 2021 The X Authors.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/the-xlang/x/parser"
	"github.com/the-xlang/x/pkg/io"
	"github.com/the-xlang/x/pkg/x"
	"github.com/the-xlang/x/pkg/x/xset"
)

func help(cmd string) {
	if cmd != "" {
		println("This module can only be used as single!")
		return
	}
	helpContent := [][]string{
		{"help", "Show help."},
		{"version", "Show version."},
		{"init", "Initialize new project here."},
	}
	maxlen := len(helpContent[0][0])
	for _, part := range helpContent {
		length := len(part[0])
		if length > maxlen {
			maxlen = length
		}
	}
	var sb strings.Builder
	const space = 5 // Space of between command name and description.
	for _, part := range helpContent {
		sb.WriteString(part[0])
		sb.WriteString(strings.Repeat(" ", (maxlen-len(part[0]))+space))
		sb.WriteString(part[1])
		sb.WriteByte('\n')
	}
	println(sb.String()[:sb.Len()-1])
}

func version(cmd string) {
	if cmd != "" {
		println("This module can only be used as single!")
		return
	}
	println("The X Programming Language\n" + x.Version)
}

func initProject(cmd string) {
	if cmd != "" {
		println("This module can only be used as single!")
		return
	}
	err := io.WriteFileTruncate(x.SettingsFile, []byte(`{
  "cxx_out_dir": "./dist/",
  "cxx_out_name": "x.cxx",
  "out_name": "main"
}`))
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	println("Initialized project.")
}

func processCommand(namespace, cmd string) bool {
	switch namespace {
	case "help":
		help(cmd)
	case "version":
		version(cmd)
	case "init":
		initProject(cmd)
	default:
		return false
	}
	return true
}

func init() {
	x.ExecutablePath = filepath.Dir(os.Args[0])
	// Not started with arguments.
	// Here is "2" but "os.Args" always have one element for store working directory.
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	var sb strings.Builder
	for _, arg := range os.Args[1:] {
		sb.WriteString(" " + arg)
	}
	os.Args[0] = sb.String()[1:]
	arg := os.Args[0]
	index := strings.Index(arg, " ")
	if index == -1 {
		index = len(arg)
	}
	if processCommand(arg[:index], arg[index:]) {
		os.Exit(0)
	}
}

func loadXSet() {
	// File check.
	info, err := os.Stat(x.SettingsFile)
	if err != nil || info.IsDir() {
		println(`X settings file ("` + x.SettingsFile + `") is not found!`)
		os.Exit(0)
	}
	bytes, err := os.ReadFile(x.SettingsFile)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	x.XSet, err = xset.Load(bytes)
	if err != nil {
		println("X settings has errors;")
		println(err.Error())
		os.Exit(0)
	}
}

func printErrors(errors []string) {
	defer os.Exit(0)
	for _, message := range errors {
		fmt.Println(message)
	}
}

func appendStandards(code *string) {
	year, month, day := time.Now().Date()
	hour, min, _ := time.Now().Clock()
	timeString := fmt.Sprintf("%d/%d/%d %d.%d (DD/MM/YYYY) (HH.MM)",
		day, month, year, hour, min)
	*code = `// Auto generated by X compiler.
// X compiler version: ` + x.Version + `
// Date:               ` + timeString + `

#pragma region X_STANDARD_IMPORTS
#include <iostream>
#include <string>
#include <functional>
#include <vector>
#include <locale.h>
#pragma endregion X_STANDARD_IMPORTS

#pragma region X_RUNTIME_FUNCTIONS
static inline void panic(const std::wstring message) {
  std::wcout << message << std::endl;
  std::exit(1);
}

template <typename ET, typename IT, typename ETET>
static inline void foreach(ET enumerable,
	                         std::function<void(IT index, ETET element)> block) {
  IT index = 0;
  for (auto element : enumerable)
  { block(index++, element); }
}

template <typename ET, typename IT>
static inline void foreach(ET enumerable, std::function<void(IT index)> block) {
  IT index = 0;
  for (auto element : enumerable)
  { block(index++); }
}
#pragma endregion X_RUNTIME_FUNCTIONS

#pragma region X_BUILTIN_TYPES
typedef size_t   size;
typedef int8_t   i8;
typedef int16_t  i16;
typedef int32_t  i32;
typedef int64_t  i64;
typedef uint8_t  u8;
typedef uint16_t u16;
typedef uint32_t u32;
typedef uint64_t u64;
typedef float    f32;
typedef double   f64;
typedef wchar_t  rune;

class str {
public:
#pragma region FIELDS
  std::wstring string;
#pragma endregion FIELDS

#pragma region CONSTRUCTORS
  str(const std::wstring& string) { this->string = string; }
  str(const rune* string)         { this->string = string; }
#pragma endregion CONSTRUCTORS

#pragma region DESTRUCTOR
  ~str() { this->string.clear(); }
#pragma endregion DESTRUCTOR

#pragma region FOREACH_SUPPORT
  typedef rune       *iterator;
  typedef const rune *const_iterator;
  iterator begin()             { return &this->string[0]; }
  const_iterator begin() const { return &this->string[0]; }
  iterator end()               { return &this->string[this->string.size()]; }
  const_iterator end() const   { return &this->string[this->string.size()]; }
#pragma endregion FOREACH_SUPPORT

#pragma region OPERATOR_OVERFLOWS
  bool operator==(const str& string) { return this->string == string.string; }
  bool operator!=(const str& string) { return !(this->string == string.string); }
  str operator+(const str& string)   { return str(this->string + string.string); }
  void operator+=(const str& string) { this->string += string.string; }

  rune& operator[](const int index) {
    const u32 length = this->string.length();
    if (index < 0) {
      panic(L"stackoverflow exception:\n index is less than zero");
    } else if (index >= length) {
      panic(L"stackoverflow exception:\nindex overflow " +
        std::to_wstring(index) + L":" + std::to_wstring(length));
    }
    return this->string[index];
  }

  friend std::wostream& operator<<(std::wostream &os, const str& string)
  { os << string.string; return os; }
#pragma endregion OPERATOR_OVERFLOWS
};
#pragma endregion X_BUILTIN_TYPES

#pragma region X_BUILTIN_VALUES
#define nil nullptr
#pragma endregion X_BUILTIN_VALUES

#pragma region X_STRUCTURES
template<typename T>
class array {
public:
#pragma region FIELDS
  std::vector<T> vector;
#pragma endregion FIELDS

#pragma region CONSTRUCTORS
  array()                                                       { this->vector = { }; }
  array(const std::vector<T>& vector)                           { this->vector = vector; }
  array(std::nullptr_t ) : array()                              { }
  array(const array<T>& arr): array(std::vector<T>(arr.vector)) { }
#pragma endregion CONSTRUCTORS

#pragma region DESTRUCTOR
  ~array() { this->vector.clear(); }
#pragma endregion DESTRUCTOR

#pragma region FOREACH_SUPPORT
  typedef T       *iterator;
  typedef const T *const_iterator;
  iterator begin()             { return &this->vector[0]; }
  const_iterator begin() const { return &this->vector[0]; }
  iterator end()               { return &this->vector[this->vector.size()]; }
  const_iterator end() const   { return &this->vector[this->vector.size()]; }
#pragma endregion FOREACH_SUPPORT

#pragma region OPERATOR_OVERFLOWS
  bool operator==(const array& array) {
    const u32 vector_length = this->vector.size();
    const u32 array_vector_length = array.vector.size();
    if (vector_length != array_vector_length) { return false; }
    for (int index = 0; index < vector_length; ++index)
    { if (this->vector[index] != array.vector[index]) { return false; } }
    return true;
  }

  bool operator==(std::nullptr_t)     { return this->vector.empty(); }
  bool operator!=(const array& array) { return !(*this == array); }
  bool operator!=(std::nullptr_t)     { return !this->vector.empty(); }

  T& operator[](const int index) {
    const u32 length = this->vector.size();
         if (index < 0) { panic(L"stackoverflow exception:\n index is less than zero"); }
    else if (index >= length) {
      panic(L"stackoverflow exception:\nindex overflow " +
        std::to_wstring(index) + L":" + std::to_wstring(length));
    }
    return this->vector[index];
  }

  friend std::wostream& operator<<(std::wostream &os, const array<T>& array) {
    os << L"[";
    const u32 size = array.vector.size();
    for (int index = 0; index < size;) {
      os << array.vector[index++];
      if (index < size) { os << L", "; }
    }
    os << L"]";
    return os;
  }
#pragma endregion OPERATOR_OVERFLOWS
};
#pragma endregion X_STRUCTURES

#pragma region X_BUILTIN_FUNCTIONS
#define _out(v) std::wcout << v
#define _outln(v) _out(v); std::wcout << std::endl
#pragma endregion X_BUILTIN_FUNCTIONS

#pragma region TRANSPILED_X_CODE
` + *code + `
#pragma endregion TRANSPILED_X_CODE

#pragma region X_ENTRY_POINT
int main() {
#pragma region X_ENTRY_POINT_STANDARD_CODES
  setlocale(0x0, "");
#pragma endregion X_ENTRY_POINT_STANDARD_CODES
  _main();

#pragma region X_ENTRY_POINT_END_STANDARD_CODES
  return EXIT_SUCCESS;
#pragma endregion X_ENTRY_POINT_END_STANDARD_CODES
}
#pragma endregion X_ENTRY_POINT`
}

func writeCxxOutput(info *parser.ParseFileInfo) {
	path := filepath.Join(x.XSet.CxxOutDir, x.XSet.CxxOutName)
	err := os.MkdirAll(x.XSet.CxxOutDir, 0511)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	err = io.WriteFileTruncate(path, []byte(info.X_CXX))
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
}

var routines *sync.WaitGroup

func main() {
	f, err := io.GetX(os.Args[0])
	if err != nil {
		println(err.Error())
		return
	}
	loadXSet()
	routines = new(sync.WaitGroup)
	info := new(parser.ParseFileInfo)
	info.File = f
	info.Routines = routines
	routines.Add(1)
	go parser.ParseFile(info)
	routines.Wait()
	if info.Errors != nil {
		printErrors(info.Errors)
	}
	appendStandards(&info.X_CXX)
	writeCxxOutput(info)
}
