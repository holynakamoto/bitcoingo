# Bitcoin Original Commit - Go Implementation

A Go implementation of Satoshi Nakamoto's original Bitcoin commit, featuring Base58 encoding/decoding and Bitcoin address generation functions.

## Overview

This project recreates the functionality from [Bitcoin's first commit](https://github.com/bitcoin/bitcoin/commit/4405b78d6059e536c36974088a8ed4d9f0f29898) in Go, providing the fundamental building blocks for Bitcoin address handling and Base58 encoding that were present in the original Bitcoin codebase.

## Features

- **Base58 Encoding/Decoding**: Convert binary data to/from Base58 format
- **Base58Check**: Base58 encoding with 4-byte checksum for error detection
- **Bitcoin Address Generation**: Create Bitcoin addresses from public keys
- **Address Validation**: Verify Bitcoin address format and checksums
- **Hash Functions**: SHA256 and simplified RIPEMD160 implementation

## Installation

1. Ensure you have Go installed (version 1.16 or later recommended)
2. Clone or download the source code
3. Run the program:

```bash
go run bitcoin.go
```

## Usage

### Basic Base58 Operations

```go
// Encode binary data to Base58
data := []byte("Hello, Bitcoin!")
encoded := EncodeBase58(data)
fmt.Printf("Base58: %s\n", encoded)

// Decode Base58 string back to binary
decoded, err := DecodeBase58(encoded)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Decoded: %s\n", decoded)
```

### Base58Check (with Checksum)

```go
// Encode with checksum protection
data := []byte("Important data")
encodedCheck := EncodeBase58Check(data)
fmt.Printf("Base58Check: %s\n", encodedCheck)

// Decode and verify checksum
decodedCheck, err := DecodeBase58Check(encodedCheck)
if err != nil {
    log.Fatal("Invalid checksum:", err)
}
fmt.Printf("Verified data: %s\n", decodedCheck)
```

### Bitcoin Address Generation

```go
// Generate a Bitcoin address from a public key
pubKey := []byte("your public key bytes here")
address := PubKeyToAddress(pubKey)
fmt.Printf("Bitcoin Address: %s\n", address)

// Validate an address
isValid := IsValidBitcoinAddress(address)
fmt.Printf("Address is valid: %t\n", isValid)
```

### Address to Hash160 Conversion

```go
address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
hash160, err := AddressToHash160(address)
if err != nil {
    log.Fatal("Invalid address:", err)
}
fmt.Printf("Hash160: %x\n", hash160)

// Convert back to address
reconstructed := Hash160ToAddress(hash160)
fmt.Printf("Reconstructed: %s\n", reconstructed)
```

## API Reference

### Core Functions

#### `EncodeBase58(input []byte) string`
Encodes a byte slice to a Base58 string using Bitcoin's alphabet.

#### `DecodeBase58(s string) ([]byte, error)`
Decodes a Base58 string back to bytes. Returns error for invalid characters.

#### `EncodeBase58Check(input []byte) string`
Encodes with a 4-byte SHA256 checksum appended for error detection.

#### `DecodeBase58Check(s string) ([]byte, error)`
Decodes and verifies the checksum. Returns error if checksum doesn't match.

### Bitcoin Address Functions

#### `PubKeyToAddress(pubKey []byte) string`
Generates a Bitcoin address from a public key by hashing and encoding.

#### `Hash160ToAddress(hash160 Hash160) string`
Converts a 160-bit hash directly to a Bitcoin address.

#### `AddressToHash160(address string) (Hash160, error)`
Extracts the 160-bit hash from a Bitcoin address.

#### `IsValidBitcoinAddress(address string) bool`
Validates whether a string is a properly formatted Bitcoin address.

### Utility Functions

#### `Hash(data []byte) Hash256`
Performs double SHA256 hashing (Bitcoin's standard hash function).

#### `HashRIPEMD160(data []byte) Hash160`
Simplified hash function (uses SHA256, truncated to 20 bytes).

## Types

```go
type Hash160 [20]byte  // 160-bit hash (20 bytes)
type Hash256 [32]byte  // 256-bit hash (32 bytes)
```

## Examples

### Complete Address Generation Workflow

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Step 1: Start with some "public key" data
    pubKey := []byte("sample public key for demonstration")
    
    // Step 2: Generate Bitcoin address
    address := PubKeyToAddress(pubKey)
    fmt.Printf("Generated Address: %s\n", address)
    
    // Step 3: Validate the address
    if IsValidBitcoinAddress(address) {
        fmt.Println("✓ Address is valid")
    } else {
        fmt.Println("✗ Address is invalid")
    }
    
    // Step 4: Extract hash160 from address
    hash160, err := AddressToHash160(address)
    if err != nil {
        log.Fatal("Error extracting hash:", err)
    }
    fmt.Printf("Hash160: %x\n", hash160)
    
    // Step 5: Reconstruct address from hash160
    reconstructed := Hash160ToAddress(hash160)
    fmt.Printf("Reconstructed: %s\n", reconstructed)
    
    // Verify they match
    if address == reconstructed {
        fmt.Println("✓ Round-trip successful")
    }
}
```

### Error Handling Example

```go
// Always check for errors when decoding
encoded := "invalid base58 string with 0 and O"
decoded, err := DecodeBase58(encoded)
if err != nil {
    fmt.Printf("Decode failed: %v\n", err)
    return
}

// Checksum verification
tamperedData := "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2" // hypothetical tampered address
_, err = DecodeBase58Check(tamperedData)
if err != nil {
    fmt.Printf("Checksum verification failed: %v\n", err)
}
```

## Technical Notes

### Base58 Alphabet
This implementation uses Bitcoin's Base58 alphabet: `123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz`

Notable exclusions:
- `0` (zero) - looks like `O`
- `O` (capital o) - looks like `0` 
- `I` (capital i) - looks like `l`
- `l` (lowercase L) - looks like `I`

### Address Version
- Current implementation uses version byte `0` (mainnet addresses starting with `1`)
- Testnet and other networks use different version bytes

### Limitations
- RIPEMD160 is simplified to use SHA256 (first 20 bytes)
- No secp256k1 public key validation
- Intended for educational/demonstration purposes

## Security Considerations

⚠️ **Important**: This implementation is for educational purposes and demonstrates the concepts from Bitcoin's original commit. For production use:

1. Use proper RIPEMD160 implementation
2. Validate public keys using secp256k1
3. Use established, audited cryptographic libraries
4. Implement proper random number generation for key creation

## License

This implementation follows the same MIT/X11 license as specified in Satoshi's original commit.

## References

- [Original Bitcoin Commit](https://github.com/bitcoin/bitcoin/commit/4405b78d6059e536c36974088a8ed4d9f0f29898)
- [Base58 Specification](https://en.bitcoin.it/wiki/Base58Check_encoding)
- [Bitcoin Address Format](https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses)
