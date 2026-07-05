# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

セキュリティ上の問題を見つけた場合は、GitHub Issues ではなく非公開で報告してください。

- リポジトリ: https://github.com/enprocode/EPGStation-Notification
- 連絡先: GitHub の Security Advisories から報告するか、リポジトリオーナー (`enprocode`) に直接連絡してください

報告後、内容を確認のうえ可能な限り早く返信します。修正版のリリース方針についても合わせてお知らせします。

## Security Recommendations

- `bin/config.yml` には API トークンが含まれます。配置後は `chmod 600 bin/config.yml` を設定してください
- トークンや Webhook URL を Issue やログに公開しないでください
- Release 成果物に含まれる `config.yml` はテンプレートです。本番用の値に必ず差し替えてください
