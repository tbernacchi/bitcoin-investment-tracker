# Bitcoin Investment Tracker

A real-time Bitcoin investment tracker that calculates your profit/loss based on your initial investment. It combines historical purchase data from [Mercado Bitcoin](https://www.mercadobitcoin.com.br/) with real-time price updates from Binance WebSocket API to provide accurate investment tracking.

The application uses your original purchase price (stored in `.env`) as a baseline and compares it with current market prices to calculate your investment returns. While real-time price data comes from Binance's WebSocket API for instant updates, the profit/loss calculation is based on your actual purchase price from Mercado Bitcoin.

## Features

- 🚀 Real-time price updates via Binance WebSocket;     
- 💰 Investment profit/loss calculation based on your actual purchase price;
- 📊 24-hour price change tracking;
- 🇧🇷 Brazilian Real (BRL) currency formatting;
- 📈 Visual indicators for price movements (↑↓=)

## Prerequisites

- Go 1.21 or higher
- Git
- Environment variables set up in `.env` file

## Installation

1. Clone the repository:
```
git clone git@github.com:tbernacchi/bitcoin-investment-tracker.git
cd bitcoin-investment-tracker
```

2. Install dependencies:
```
go mod tidy
```

3. Create a `.env` file in the root directory:
```
MERCADO_BITCOIN_API_URL=https://www.mercadobitcoin.net/api/BTC/ticker/  # Used for historical reference
BINANCE_USDT_API_URL=https://api.binance.com/api/v3/ticker/price?symbol=USDTBRL  # USDT/BRL conversion

# Update frequency for calculations
CHECK_INTERVAL=1s  # How often to update the investment calculations

# Your investment details
INVESTMENT_BRL=20000     # Your initial investment in BRL
BTC_AMOUNT=0.05912530    # Amount of BTC you own
```

## Usage

Run the application:
```
go run cmd/main.go
```

The application will:
- Connect to Binance WebSocket for real-time BTC/BRL price updates;
- Display formatted prices with 24h change percentage;
- Calculate and show your investment returns;
- Update the web interface with current data;

Example output (logs):
```
2025/04/28 13:05:04 BTC Price: R$ 534.186,00 (↓ -0.65%)
2025/04/28 13:05:05 BTC Price: R$ 534.185,00 (↓ -0.65%)
2025/04/28 13:05:07 BTC Price: R$ 534.229,00 (↑ 0.12%)
```

## Project Structure

```
bitcoin-investment-tracker/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── calculator/       # Investment calculations
│   ├── webserver/       # Web interface
│   └── websocket/       # Binance WebSocket client
├── .env                 # Environment variables
├── go.mod              # Go modules file
├── go.sum              # Go modules checksum
└── README.md           # This file
```

## Technical Details

- Uses Binance WebSocket API for real-time price data;
- Implements goroutines for concurrent processing;
- Maintains persistent WebSocket connection;
- Handles JSON parsing of ticker data;
- Formats currency in Brazilian Real (BRL);

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Binance WebSocket API](https://binance-docs.github.io/apidocs/spot/en/#websocket-market-streams)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

## Author

- **Tadeu Bernacchi** - [GitHub](https://github.com/tbernacchi)

## Future Improvements

- [ ] Add Prometheus metrics export
- [ ] Implement price alerts system
- [ ] Add support for multiple cryptocurrencies
- [ ] Create Docker deployment
- [ ] Add automated tests


## Environment Configuration

The application uses a `.env` file to store your investment details and API configurations. 

Here's what each variable means:

```
# API URLs for price fetching
MERCADO_BITCOIN_API_URL=https://www.mercadobitcoin.net/api/BTC/ticker/  # Used for historical reference
BINANCE_USDT_API_URL=https://api.binance.com/api/v3/ticker/price?symbol=USDTBRL  # USDT/BRL conversion

# Update frequency for calculations
CHECK_INTERVAL=1s  # How often to update the investment calculations

# Your investment details
INVESTMENT_BRL=20000     # Your initial investment in BRL
BTC_AMOUNT=0.05912530    # Amount of BTC you own
```

### Using With Different Exchanges

While this project was initially built to track investments made on Mercado Bitcoin, you can easily adapt it for investments made on any exchange:

1. Simply update the `INVESTMENT_BRL` with your initial investment amount in Brazilian Reais
2. Update `BTC_AMOUNT` with the amount of Bitcoin you purchased
3. The real-time price tracking will work regardless of where you bought your Bitcoin

For example, if you bought Bitcoin on Binance:
```
INVESTMENT_BRL=25000      # How much you invested in BRL
BTC_AMOUNT=0.06000000    # How much BTC you got
CHECK_INTERVAL=1s        # Keep the 1-second update interval
```

The application will:
- Use your provided investment amount and BTC quantity
- Track current prices via Binance WebSocket
- Calculate your profit/loss based on your actual purchase data
- Show real-time updates of your investment performance

```
go mod init go-bitcoin-price
go mod tidy 
```
