# 導入手順

## 必要環境

- Linux amd64（EPGStation 実行環境）
- ビルド時: Go 1.26 以上

## 1. 取得

### GitHub Release（推奨）

[Releases](https://github.com/enprocode/EPGStation-Notification/releases) から `epgst-notify-v{version}-linux-amd64.zip` をダウンロードし、任意のディレクトリに展開してください。

### ソースからビルド

```shell
git clone https://github.com/enprocode/EPGStation-Notification.git
cd EPGStation-Notification
GOOS=linux GOARCH=amd64 go build -o bin/epgst-notify main.go
```

展開後の構成:

```
/path/to/bin/
  epgst-notify
  config.yml
```

## 2. Slack 設定

1. [Slack API](https://api.slack.com) でアプリを作成
2. Bot Token と通知先チャンネル ID を取得
3. 参考: [Slack API 日本語ドキュメント](https://api.slack.com/lang/ja-jp)

## 3. Discord 設定（Discord 利用時）

1. 通知先チャンネル → 連携サービス → Webhook を作成
2. Webhook URL から ID とトークンを取得

```
https://discord.com/api/webhooks/{webhook-id}/{webhook-token}
```

## 4. config.yml

`config.yml` は **`epgst-notify` と同じディレクトリ** に置きます。

```yaml
slack-config:
  slack-token: "xoxb-..."
  channel: "C0123456789"

discord-config:
  discord-webhook-token: "your-webhook-token"
  discord-webhook: "1234567890123456789"
```

```shell
chmod 600 /path/to/bin/config.yml
```

## 5. EPGStation 設定

`/path/to/epgstation/config/config.yml` に外部コマンドを追加します。

**Slack**

```yaml
recordingStartCommand: "/path/to/bin/epgst-notify slackRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify slackRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify slackRecError"
```

**Discord**

```yaml
recordingStartCommand: "/path/to/bin/epgst-notify discordRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify discordRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify discordRecError"
```

## 6. 再起動

```shell
sudo pm2 restart epgstation
```

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `slackRecStart` | 録画開始（Slack） |
| `slackRecEnd` | 録画終了（Slack） |
| `slackRecError` | 録画エラー（Slack） |
| `discordRecStart` | 録画開始（Discord） |
| `discordRecEnd` | 録画終了（Discord） |
| `discordRecError` | 録画エラー（Discord） |

## 通知内容

- 番組名、放送局、開始/終了時刻、番組概要
- 録画パス（終了・エラー時）
- 録画エラー時: `ERROR_CNT` / `DROP_CNT` / `LOGPATH`
