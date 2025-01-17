: Copyright 2021 The Jule Programming Language.
: Use of this source code is governed by a BSD 3-Clause
: license that can be found in the LICENSE file.

@echo off

if exist .\julec.exe ( del /f julec.exe )

if exist cmd\julec\main.go (
  go build -o julec.exe -v cmd\julec\main.go
) else (
  go build -o julec.exe -v ..\cmd\julec\main.go
)

if exist .\julec.exe (
  echo Compile is successful!
) else (
  echo -----------------------------------------------------------------------
  echo An unexpected error occurred while compiling JuleC. Check errors above.
)
