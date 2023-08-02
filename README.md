# Replit Scraper
## What is this
Replitの検索を自動化するツール  
![Preview](https://img001.prntscr.com/file/img001/QM9k2ABpRF6F2NkAvcEAlg.png)

## How to use
このレポジトリをクローンして、`go build`でビルドした実行ファイルを実行してください  
Configは引数から指定できます（デフォルトは同じディレクトリ内のConfig） 

`sessions.txt`には`connect.sid`クッキーを入れてください  
安定して高速に検索を行うには、少なくとも四つ以上のアカウントが必要です

Massiveオプションを使うには、`crosis_api`を実行させておく必要があります  
ディレクトリの検索には時間が掛かります(websocketを使用しているため)

## Setting up CrosisApi
`npm install`でパッケージをインストールした後、start.batを起動するだけです

## Config example
```json
{
  "scraper": {
    "sessions": "sessions.txt",
    "proxies": "proxies.txt",
    "page_limit": 21,
    "search_delay": 1000,

    "parallel": false,
    "finder": true,
    "massive": true,

    "finders": {
      "discord_token": {
        "active": true,
        "bot": true
      },
      "proxy": {
        "active": true
      },
      "email": {
        "active": true
      },
      "password": {
        "active": true
      },
      "phone": {
        "active": true
      },
      "captcha_service": {
        "active": true,
        "min_balance": 1
      },
      "openai_key": {
        "active": true
      },
      "google_api_key": {
        "active": true
      },
      "telegram_token": {
        "active": true
      }
    }
  },
  "search": {
    "query": "paypay",
    "sort": "RecentlyModified",
    "exact": true
  }
}
```

## Notes
このツールは自己責任で実行してください