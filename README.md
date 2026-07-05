# EPGStation-Notification

EPGStation の録画開始・終了・エラー時に Slack / Discord へ通知する CLI ツール **`epgst-notify`** です。

[![CI](https://github.com/enprocode/EPGStation-Notification/actions/workflows/ci.yml/badge.svg)](https://github.com/enprocode/EPGStation-Notification/actions/workflows/ci.yml)

## 機能

- 録画開始・終了・失敗の通知（Slack / Discord）
- EPGStation から渡される環境変数（番組名、放送局、時刻、録画パスなど）を通知に反映
- 録画失敗時は `ERROR_CNT` / `DROP_CNT` / `LOGPATH` も表示

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `slackRecStart` | 録画開始（Slack） |
| `slackRecEnd` | 録画終了（Slack） |
| `slackRecError` | 録画エラー（Slack） |
| `discordRecStart` | 録画開始（Discord） |
| `discordRecEnd` | 録画終了（Discord） |
| `discordRecError` | 録画エラー（Discord） |

## 必要環境

- Linux amd64（EPGStation 実行環境向け）
- ビルド時: Go 1.25 以上

## 導入手順

### 1. 取得

**GitHub Release から取得（推奨）**

[Releases](https://github.com/enprocode/EPGStation-Notification/releases) から `epgst-notify-v{version}-linux-amd64.zip` をダウンロードし、任意のディレクトリに展開してください。

**ソースからビルド**

```shell
git clone https://github.com/enprocode/EPGStation-Notification.git
cd EPGStation-Notification
GOOS=linux GOARCH=amd64 go build -o bin/epgst-notify main.go
```

展開後の構成例:

```
/path/to/bin/
  epgst-notify    # 実行ファイル
  config.yml      # 設定ファイル
```

### 2. Slack アプリを作成

- URL: https://api.slack.com
- 参考: https://api.slack.com/lang/ja-jp
- Bot Token と通知先チャンネル ID を取得

### 3. Discord Webhook を作成（Discord 利用時）

1. 通知先チャンネルの設定 → 連携サービス → Webhook を作成
2. Webhook URL から ID とトークンを取得

   ```
   https://discord.com/api/webhooks/{webhook-id}/{webhook-token}
   ```

### 4. 設定ファイルを編集

`config.yml` は **`epgst-notify` と同じディレクトリ** に配置してください。

```shell
vim /path/to/bin/config.yml
```

```yaml
slack-config:
  slack-token: "xoxb-..."
  channel: "C0123456789"

discord-config:
  discord-webhook-token: "your-webhook-token"
  discord-webhook: "1234567890123456789"
```

トークンを含むため、権限を制限してください。

```shell
chmod 600 /path/to/bin/config.yml
```

### 5. EPGStation へ設定追加

`/path/to/epgstation/config/config.yml` に外部コマンドを追加します。

**Slack のみ使う場合**

```yaml
recordingStartCommand: "/path/to/bin/epgst-notify slackRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify slackRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify slackRecError"
```

**Discord のみ使う場合**

```yaml
recordingStartCommand: "/path/to/bin/epgst-notify discordRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify discordRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify discordRecError"
```

### 6. EPGStation を再起動

```shell
sudo pm2 restart epgstation
```

## 開発

### テスト

```shell
go vet ./...
go test ./...
```

### Docker でビルド

`go.mod` と同じ Go 1.25 環境でビルドする場合:

```shell
cd docker
docker compose up -d --build
docker compose exec go go vet ./...
docker compose exec go go test ./...
./docker.sh
```

`docker.sh` は Linux amd64 向けに `bin/epgst-notify` をビルドします。

### リリース

バージョンはリポジトリ直下の `VERSION` ファイルで管理します。

```shell
echo "1.1.0" > VERSION
git add VERSION
git commit -m "chore: bump version to 1.1.0"
git tag v1.1.0
git push origin main
git push origin v1.1.0
```

`v*` タグを push すると GitHub Actions が Release を作成します。  
`VERSION` の内容とタグ名（`v` 除く）が一致しない場合は失敗します。

Release 成果物:

- `epgst-notify-v{version}-linux-amd64.zip`
- `checksums.txt`

### CI / 自動化

| ワークフロー | 内容 |
|-------------|------|
| CI | vet / test / build |
| Release | タグ push 時にビルド・Release 作成 |
| CodeQL | セキュリティスキャン |
| Dependency Review | 依存関係の脆弱性チェック |

Dependabot による patch / minor 更新は Mergify で自動マージされます。major 更新は手動レビューが必要です。

**メンテナ向け:** Dependency Review を有効にするには、リポジトリ設定で [Dependency graph](https://github.com/enprocode/EPGStation-Notification/settings/security_analysis) を有効にしてください。

## セキュリティ

脆弱性の報告方法は [SECURITY.md](SECURITY.md) を参照してください。

## License

MIT License — 詳細は [LICENSE](LICENSE) を参照してください。
