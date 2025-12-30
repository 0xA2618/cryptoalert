# CryptoAlert

CryptoAlert is a powerful cryptocurrency price monitoring and alerting tool built with Go. it monitors Binance Futures market data and sends real-time notifications to Telegram based on technical indicators like MACD and RSI.

## Features

- **Real-time Monitoring**: Automatically fetches all available symbols from Binance Futures.
- **Multiple Timeframes**: Supports multiple monitoring cycles (e.g., 5m, 30m, 1h, 4h).
- **Technical Analysis**:
  - **MACD Crosses**: Detects Gold and Death crosses.
  - **RSI Filtering**: Filters signals based on RSI overbought/oversold levels.
  - **Trend Reversal**: Identifies potential momentum exhaustion in the last 3 K-lines.
- **Detailed Alerts**:
  - Current price and latest K-line price change percentage.
  - MACD cross type.
  - RSI value.
  - Current funding rate.
- **Telegram Integration**: Categorized alerts sent to different Telegram topics based on timeframe and importance.

## Project Structure

- `base/`: Core data fetching logic (klines, symbols).
- `calculate/`: Technical analysis engines (MACD, RSI, Funding Rate).
- `config/`: Configuration management using JSON and Viper.
- `main/`: Application entry point and lifecycle management.
- `utils/`: Utility functions (Telegram notifications).

## Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/0xA2618/cryptoalert.git
   cd cryptoalert
   ```

2. **Configure the application**:
   Create a `config/config.json` file based on the required structure (refer to `config/config.go`).
   > [!IMPORTANT]
   > `config/config.json` is ignored by git for security. Do not commit your API tokens or credentials.

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Run the application**:
   ```bash
   go run main/main.go
   ```

## Requirements

- Go 1.25.4 or higher.
- A Telegram Bot token and Chat ID.
- Access to Binance Futures API (no API key required for public data).

## License

This project is licensed under the MIT License.
