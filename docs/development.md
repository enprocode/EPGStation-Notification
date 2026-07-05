# 開発ガイド

## テスト

```shell
go vet ./...
go test -race ./...
```

`-race` は CI と同じ条件で競合検出を行うためのオプションです。

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

バージョンはリポジトリ直下の `VERSION` で管理し、[セマンティックバージョニング](https://semver.org/lang/ja/)に従います（バグ修正のみは patch、後方互換の機能追加は minor）。

1. `VERSION` を更新し、ブランチで PR を作成してマージします（例: patch リリース `1.0.1`）。

   ```shell
   git switch -c chore/bump-1.0.1
   echo "1.0.1" > VERSION
   git commit -am "chore: bump version to 1.0.1"
   git push -u origin chore/bump-1.0.1
   gh pr create --fill
   # CI 通過後にマージ
   ```

2. `main` を最新化し、`VERSION` と一致するタグを push します。

   ```shell
   git switch main && git pull
   git tag -a v1.0.1 -m "Release v1.0.1"
   git push origin v1.0.1
   ```

- `v*` タグ push で GitHub Actions が Release を作成
- `VERSION` とタグ名（`v` 除く）が一致しないと失敗
- 成果物: `epgst-notify-v{version}-linux-amd64.zip`, `checksums.txt`
- リリースノートはコミット/PR から自動生成されます

## CI / 自動化

| ワークフロー | 内容 |
|-------------|------|
| CI | vet / test / build |
| Release | タグ push 時に Release 作成 |
| CodeQL | セキュリティスキャン |
| Dependency Review | 依存関係の脆弱性チェック |

Dependabot の patch / minor は Mergify で自動マージ。major は手動レビュー。

**メンテナ向け:** [Dependency graph](https://github.com/enprocode/EPGStation-Notification/settings/security_analysis) を有効にしてください。
