# MancoStats

MancoStats is a small Discord bot that provides statistics for players on PUBG.

It is written in Go and uses the [discordgo](https://github.com/bwmarrin/discordgo) library.

## Setup

Create a `.env` file with the following variables:

- `DISCORD_TOKEN`: Your Discord bot token
- `PUBG_API_TOKEN`: Your PUBG API key
- `PUBG_API_URL`: The URL of the PUBG API, if not provided it will default to the official PUBG API
- `BOT_PREFIX`: The prefix you must use to invoke the bot on Discord, if not specified it will default to `!manco`

## Build

Run the following command to build the bot for all platforms and archs:

```bash
make build
```

## Usage

Run the binary to start the bot:

```bash
./mancostats
```

The bot will connect to Discord and listen for commands. Use `!manco help` to get a list of available commands.

## Roadmap
- More commands for stats
- Auto-detect PUBG season
- I18n system for commands and responses
