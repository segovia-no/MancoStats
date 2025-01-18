# MancoStats

MancoStats is a small Discord bot that provides statistics for players on PUBG.

It is written in Go and uses the [discordgo](https://github.com/bwmarrin/discordgo) library.

## Setup

Create a `.env` file with the following variables:

- `DISCORD_TOKEN`: Your Discord bot token
- `PUBG_API_TOKEN`: Your PUBG API key

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
