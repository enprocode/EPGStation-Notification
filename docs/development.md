# 開発ガイド

## テスト

```shell
go vet ./...
go test ./...
```

## ビルド

```shell
GOOS=linux GOARCH=amd64 go build -o bin/epgst-notify main.go
```

## Docker

Go 1.26 環境でビルドする場合:

```shell
cd docker
docker compose up -d --build
docker compose exec go go vet ./...
docker compose exec go go test ./...
./docker.sh
```

`docker.sh` は `/opt/src/main.go` から Linux amd64 向け `bin/epgst-notify` を生成します。

## リリース

バージョンはリポジトリ直下の `VERSION` で管理します。

```shell
echo "1.1.0" > VERSION
git add VERSION
git commit -m "chore: bump version to 1.1.0"
git tag v1.1.0
git push origin main
git push origin v1.1.0
```

- `v*` タグ push で GitHub Actions が Release を作成
- `VERSION` とタグ名（`v` 除く）が一致しないと失敗
- 成果物: `epgst-notify-v{version}-linux-amd64.zip`, `checksums.txt`

## CI / 自動化

| ワークフロー | 内容 |
|-------------|------|
| CI | vet / test / build |
| Release | タグ push 時に Release 作成 |
| CodeQL | セキュリティスキャン |
| Dependency Review | 依存関係の脆弱性チェック |

Dependabot の patch / minor は Mergify で自動マージ。major は手動レビュー。

**メンテナ向け:** [Dependency graph](https://github.com/enprocode/EPGStation-Notification/settings/security_analysis) を有効にしてください。
