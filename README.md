# EPGStation-Notification
EPGStationの録画開始・終了・エラー取得した内容をSlack・Discord通知スクリプトです。  
Slack and Discord notification script for EPGStation recording start, end, and error acquisition.

[![CI](https://github.com/enprocode/EPGStation-Notification/actions/workflows/ci.yml/badge.svg)](https://github.com/enprocode/EPGStation-Notification/actions/workflows/ci.yml)

## 導入手順
### Slackアプリを作成
URL: https://api.slack.com

参考: https://api.slack.com/lang/ja-jp

### スクリプト配置
```shell script
$ git clone https://github.com/enprocode/EPGStation-Notification.git
```

### Slackの設定
```shell script
# 編集
$ cd EPGStation-Notification
$ vim bin/config.yml
slack-config:
  slack-token: "SLACK_API_TOKEN"
  channel: "SLACK_CHANNEL_ID"

discord-config:
  discord-webhook-token: "DISCORD_WEBHOOK_TOKEN"
  discord-webhook: "1234567890123456789"
```

`config.yml` には API トークンが含まれるため、権限を制限してください。

```shell script
$ chmod 600 bin/config.yml
```

### EPGStationへ設定追加
```shell script
$ vim /path/to/config/config.yml
# Slack
---
recordingStartCommand: "/path/to/bin/epgst-notify slackRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify slackRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify slackRecError"

# Discord
---
recordingStartCommand: "/path/to/bin/epgst-notify discordRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify discordRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify discordRecError"
```
### EPGStation再起動
```shell script
$ sudo pm2 restart epgstation
```

### Build

`epgst-notify` というバイナリファイルができればOK
```shell script
$ cd EPGStation-Notification
$ GOOS=linux GOARCH=amd64 go build -o "bin/epgst-notify" main.go
$ ls bin
  epgst-notify
```

## リリース

バージョンはリポジトリ直下の `VERSION` ファイルで管理します。

```shell script
# 1. VERSION を更新
$ echo "1.1.0" > VERSION
$ git add VERSION
$ git commit -m "chore: bump version to 1.1.0"

# 2. タグを作成して push（GitHub Actions が Release を作成）
$ git tag v1.1.0
$ git push origin main
$ git push origin v1.1.0
```

GitHub Actions の Release ワークフローは `v*` タグの push で起動します。  
`VERSION` の内容とタグ名（`v` 除く）が一致しない場合は失敗します。

## License

MIT License — 詳細は [LICENSE](LICENSE) を参照してください。

