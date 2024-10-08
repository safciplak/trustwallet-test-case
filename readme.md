# Trust Wallet Ethereum Parser

![Check](https://github.com/trustwallet/assets/workflows/Check/badge.svg)

## Overview

This project aims to develop an Ethereum blockchain parser that enables querying transactions for subscribed addresses, enhancing the Trust Wallet user experience.

## Problem Statement

Users are currently unable to receive push notifications for incoming/outgoing transactions. By implementing a Parser interface, we can connect this functionality to our notification service, allowing us to alert users about any incoming/outgoing transactions in real-time.

## Project Goal

Create a robust Ethereum blockchain parser that efficiently queries transactions for subscribed addresses, enabling seamless integration with our notification system.

## Getting Started

To set up your development environment:

```sh
cp .env.example .env
```

## Available Scripts

We provide several scripts for maintainers:

- `make dev` - Set up and run the development environment
- `make test` - Execute automated tests

## Additional Information

1. Storage options are extensible via the `.env` config file. The `storage/storage.go` interface can be implemented for new storage solutions such as Redis, MongoDB, MySQL, etc.
2. For push notifications, transactions are currently logged as no dedicated service is provided yet.

## License

This project is released under the [MIT License](LICENSE)