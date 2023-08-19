# Replit Scraper
## What is this
Tools to automate Replit searches  
![Preview](https://img001.prntscr.com/file/img001/QM9k2ABpRF6F2NkAvcEAlg.png)

## How to use
Clone this repository and run the executable built with `go build`  
Config can be specified via argument (default is Config in the same directory)

Please put the `connect.sid` cookie in `sessions.txt`  
At least four accounts are required for stable and fast searches

To use the `massive` option, you need to have `crosis_api` running  
Directory lookups can take a while (because of websocket use)

## Setting up CrosisApi
After installing the package with `npm install`, simply run start.bat

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
Please run this tool at your own risk!