Crypto Telegram Bot

A Go bot that sends live gold, silver, and cryptocurrency prices to a Telegram chat every 4 hours.

What it does

    Fetches current prices for Gold (XAU/USD), Silver (XAG/USD), Bitcoin (BTC/USD), Ethereum (ETH/USD), and Solana (SOL/USD)

    Reads prices from Chainlink oracle feeds on the Ethereum network

    Fetches all assets concurrently for speed

    Sends a formatted message with all prices to a configured Telegram chat

    Repeats automatically every 4 hours

    Runs as a background service that restarts on failure or reboot
