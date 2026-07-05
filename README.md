# EPGStation-Notification

EPGStation の録画開始・終了・エラー時に Slack / Discord へ通知する CLI **`epgst-notify`** です。

[![CI](https://github.com/enprocode/EPGStation-Notification/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/enprocode/EPGStation-Notification/actions/workflows/ci.yml)
[![CodeQL](https://github.com/enprocode/EPGStation-Notification/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/enprocode/EPGStation-Notification/actions/workflows/codeql.yml)
[![Release](https://img.shields.io/github/v/release/enprocode/EPGStation-Notification?sort=semver)](https://github.com/enprocode/EPGStation-Notification/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/enprocode/EPGStation-Notification)](go.mod)
[![License](https://img.shields.io/github/license/enprocode/EPGStation-Notification)](LICENSE)

## クイックスタート

1. [Releases](https://github.com/enprocode/EPGStation-Notification/releases) から zip を取得して展開し、`epgst-notify` を配置
2. `cp config.example.yml config.yml` でテンプレートをコピーし、`config.yml` を編集して `chmod 600 config.yml`
3. EPGStation の `config.yml` に外部コマンドを追加して再起動

```yaml
recordingStartCommand: "/path/to/bin/epgst-notify slackRecStart"
recordingFinishCommand: "/path/to/bin/epgst-notify slackRecEnd"
recordingFailedCommand: "/path/to/bin/epgst-notify slackRecError"
```

## ドキュメント

| ドキュメント | 内容 |
|-------------|------|
| [導入手順](docs/setup.md) | 取得方法、Slack / Discord 設定、EPGStation 連携 |
| [開発ガイド](docs/development.md) | ビルド、テスト、Docker、リリース、CI |

## Security

脆弱性を発見した場合の報告方法は [SECURITY.md](SECURITY.md) を参照してください。

## License

MIT License — 詳細は [LICENSE](LICENSE) を参照してください。
