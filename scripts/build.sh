#!/usr/bin/sh
# Copyright 2021 The X Programming Language.
# Use of this source code is governed by a BSD 3-Clause
# license that can be found in the LICENSE file.

if [ -f cmd/xxc/main.go ]; then
  MAIN_FILE="cmd/xxc/main.go"
else
  MAIN_FILE="../cmd/xxc/main.go"
fi

go build -o xxc.out -v $MAIN_FILE

if [ $? -eq 0 ]; then
  echo "Compile is successful!"
else
  echo "-----------------------------------------------------------------------"
  echo "An unexpected error occurred while compiling X. Check errors above."
fi

