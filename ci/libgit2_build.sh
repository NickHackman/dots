#!/usr/bin/env bash
git clone --depth=1 -b maint/v1.0 https://github.com/libgit2/libgit2.git
mkdir libgit2/build
cd libgit2/build || exit 1
cmake .. -DCMAKE_INSTALL_PREFIX=../_install -DBUILD_CLAR=OFF
cmake --build . --target install
