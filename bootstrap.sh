#!/data/data/com.termux/files/usr/bin/bash

set -e

echo "================================"
echo "     OpenNet Clean Bootstrap"
echo "================================"
echo

if [ ! -d ".git" ]; then
    echo "ERROR: Run this from the OpenNet repository root."
    exit 1
fi

echo "This bootstrap will rebuild the project."
echo "README.md, LICENSE and .git/ will be preserved."
echo
