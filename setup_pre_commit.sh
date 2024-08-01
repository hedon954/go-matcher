#!/bin/bash

set -e

# 安装 pre-commit 工具
if ! command -v pre-commit &> /dev/null
then
    echo "pre-commit not found, installing..."
    brew install pre-commit
else
    echo "pre-commit is already installed"
fi

# 检查并 Go 版本
if ! command -v go &> /dev/null
then
    echo "Go not found, please install Go first."
    exit 1
else
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    MIN_GO_VERSION="1.22"

    if [ "$(printf '%s\n' "$MIN_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_GO_VERSION" ]; then
        echo "Go version $GO_VERSION is not supported. Please install Go version $MIN_GO_VERSION or higher."
        exit 1
    fi
fi

echo "Installing Go tools..."

# 安装 go-imports
if ! command -v goimports &> /dev/null
then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
else
    echo "goimports is already installed"
fi

# 安装 go-critic
if ! command -v gocritic &> /dev/null
then
    echo "Installing gocritic..."
    go install github.com/go-critic/go-critic/cmd/gocritic@latest
else
    echo "gocritic is already installed"
fi

# 安装 golangci-lint
if ! command -v golangci-lint &> /dev/null
then
    echo "Installing golangci-lint..."
     go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
else
    echo "golangci-lint is already installed"
fi


# 安装 typos
if ! command -v typos &> /dev/null
then
    echo "Installing typos..."
    brew install typos-cli
else
    echo "typos is already installed"
fi

# 安装 pre-commit hooks
echo "Installing pre-commit hooks..."
pre-commit install

echo "All done! pre-commit hooks are installed and configured."