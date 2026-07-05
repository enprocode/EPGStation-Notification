# EPGStation-Notification
EPGStationの録画開始・終了・エラー取得した内容をSlack・Discord通知スクリプトです。  
Slack and Discord notification script for EPGStation recording start, end, and error acquisition.

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/Junch25/EPGStation-Notification/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/Junch25/EPGStation-Notification/tree/main)

## 導入手順
### Slackアプリを作成
URL: https://api.slack.com

参考: https://api.slack.com/lang/ja-jp

### スクリプト配置
```shell script
$ git clone https://github.com/Junch25/EPGStation-Notification.git
```

### Slackの設定
```shell script
# 編集
$ cd epgstation_notification
$ vim bin/config.yml
slack-config:
  slack-token: "SLACK_API_TOKEN"
  channel: "SLACK_CHANNEL_ID"

discord-config:
  discord-webhook-token: "DISCORD_API_TOKEN"
  discord-webhook: "DISCORD_API_WEBHOOK"
```

### EPGStationへ設定追加
```shell script
$ vim /path/to/config/config.yml
# Slack
---
recordingStartCommand: "/path/to/bin/epgstation-notification slackRecStart"
recordingFinishCommand: "/path/to/bin/epgstation-notification slackRecEnd"
recordingFailedCommand: "/path/to/bin/epgstation-notification slackRecError"

# Discord
---
recordingStartCommand: "/path/to/bin/epgstation-notification discordRecStart"
recordingFinishCommand: "/path/to/bin/epgstation-notification discordRecEnd"
recordingFailedCommand: "/path/to/bin/epgstation-notification discordRecError"
```
### EPGStation再起動
```shell script
$ sudo pm2 restart epgstation
```

### Build

`EPGStation-Notification`というバイナリファイルできればOK
```shell script
$ cd epgstation_notification
$ GOOS=linux GOARCH=amd64 go build -o "bin/epgstation-notification" main.go
$ ls bin
  epgstation-notification
```

## リリース

バージョンはリポジトリ直下の `VERSION` ファイルで管理します。

```shell script
# 1. VERSION を更新
$ echo "1.1.0" > VERSION
$ git add VERSION
$ git commit -m "chore: bump version to 1.1.0"

# 2. タグを作成して push（GitHub Actions / CircleCI が Release を作成）
$ git tag v1.1.0
$ git push origin main
$ git push origin v1.1.0
```

GitHub Actions の Release ワークフローは `v*` タグの push で起動します。  
`VERSION` の内容とタグ名（`v` 除く）が一致しない場合は失敗します。

## License

MIT License — 詳細は [LICENSE](LICENSE) を参照してください。

