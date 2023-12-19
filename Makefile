# ターゲットを指定せず、単に `make` のみを実行した場合のデフォルトターゲット
.DEFAULT_GOAL := test

fmt:
# goimports は go fmt の機能強化版
# アルファベット順ソートや使われていないものの削除、不足を補うなど
# -l: フォーマットが正しくないファイルをコンソールに表示
# -w: ファイルを直接書き換える
	goimports -l -w .
# .PHONY is used to specify that the target (this time, fmt) is not a file
.PHONY: fmt

# 単に make すると lint は実行されない（vet: fmt として lint を省略しているため）
# そのため、make lint と明示的に指定する必要がある
lint: fmt
# staticcheck はスタイルガイドとの整合をかなり詳細に確認する
# https://staticcheck.io/docs/checks
	-staticcheck ./...
.PHONY: lint

vet: lint
# go vet はコンパイルエラーに留まらない潜在的なバグ（と疑われるコード）を見つける
	-go vet ./...
.PHONY: vet

build: vet
# 必要なサードパーティライブラリのダウンロードと、不要になったファイルの削除
	go mod tidy
# コマンドに利用していないパッケージも含め、全てのパッケージをビルドする
# ビルドエラーが出ないことの確認のためであり、本来の動作上は不必要
	go build ./...
# 各コマンドをビルドして実行バイナリを生成する
# 生成するコマンドが増えた場合、ここに行を追加する必要がある
	go build -o bin/match ./cmd/match
.PHONY: build

test: build
	go test ./...
.PHONY: test

# デフォルトの `make` では実行されない、オプションとしてのターゲット
test_v: build
# -v: テストの詳細を表示する
	go test -v ./...
.PHONY: test_v