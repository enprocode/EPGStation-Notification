# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

セキュリティ上の問題を見つけた場合は、GitHub Issues ではなく非公開で報告してください。

- リポジトリの **Security** タブ → **Report a vulnerability** から非公開で報告してください（GitHub Private Vulnerability Reporting）
- 直接リンク: https://github.com/enprocode/EPGStation-Notification/security/advisories/new

**初回応答の目安は 72 時間以内**です。内容を確認のうえ、修正版のリリース方針についても合わせてお知らせします。公開は修正リリース後に調整します。

## Security Recommendations

- `config.yml` は **`epgst-notify` と同じディレクトリ** に配置してください
- `config.yml` には API トークンが含まれます。配置後は `chmod 600 config.yml` を設定してください
- **実トークンを記入した `config.yml` を絶対にコミットしないでください。** リポジトリで追跡されるのはプレースホルダ入りの `bin/config.example.yml` のみで、実ファイル `bin/config.yml` は `.gitignore` 済みです
- トークンや Webhook URL を Issue やログに公開しないでください
- [GitHub Release](https://github.com/enprocode/EPGStation-Notification/releases) 成果物にはテンプレート `config.example.yml` が含まれます。`config.yml` にコピーし、本番用の値に必ず差し替えてください
