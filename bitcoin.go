// Copyright (c) 2009 Satoshi Nakamoto
// Distributed under the MIT/X11 software license, see the accompanying
// file license.txt or http://www.opensource.org/licenses/mit-license.php.
//
// Go implementation of the original Bitcoin commit by Satoshi Nakamoto

package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Why base-58 instead of standard base-64 encoding?
// - Don't want 0OIl characters that look the same in some fonts and
//   could be used to create visually identical looking account numbers.
// - A string with non-alphanumeric characters is not as easily accepted as an account number.
// - E-mail usually won't line-break if there's no punctuation to break at.
// - Doubleclicking selects the whole number as one word if it's all alphanumeric.

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
const addressVersion = 0

var (
	bigRadix   = big.NewInt(58)
	bigZero    = big.NewInt(0)
	ErrInvalidBase58 = errors.New("invalid base58 character")
)

// Hash160 represents a 160-bit hash (20 bytes)
type Hash160 [20]byte

// Hash256 represents a 256-bit hash (32 bytes) 
type Hash256 [32]byte

// EncodeBase58 encodes a byte slice to base58 string
func EncodeBase58(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	// Convert big endian data to little endian
	// Extra zero at the end make sure bignum will interpret as a positive number
	inputReversed := make([]byte, len(input)+1)
	for i, b := range input {
		inputReversed[len(input)-1-i] = b
	}
	inputReversed[len(input)] = 0

	// Convert little endian data to bignum
	bn := new(big.Int).SetBytes(reverse(inputReversed))

	// Convert bignum to string
	var result strings.Builder
	result.Grow((len(input)*138)/100 + 1) // Reserve space

	for bn.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		bn.DivMod(bn, bigRadix, mod)
		result.WriteByte(base58Alphabet[mod.Int64()])
	}

	// Leading zeroes encoded as base58 zeros
	for _, b := range input {
		if b != 0 {
			break
		}
		result.WriteByte(base58Alphabet[0])
	}

	// Convert little endian string to big endian
	return reverseString(result.String())
}

// DecodeBase58 decodes a base58 string to byte slice
func DecodeBase58(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, nil
	}

	// Skip leading whitespace
	s = strings.TrimLeft(s, " \t\n\r")
	if len(s) == 0 {
		return nil, nil
	}

	bn := new(big.Int)
	bnChar := new(big.Int)

	// Convert big endian string to bignum
	for _, c := range s {
		idx := strings.IndexRune(base58Alphabet, c)
		if idx == -1 {
			// Check if remaining characters are whitespace
			remaining := strings.TrimLeft(string(c), " \t\n\r")
			if len(remaining) > 0 {
				return nil, ErrInvalidBase58
			}
			break
		}
		bnChar.SetInt64(int64(idx))
		bn.Mul(bn, bigRadix)
		bn.Add(bn, bnChar)
	}

	// Get bignum as little endian data
	tmpBytes := bn.Bytes()
	
	// Trim off sign byte if present
	if len(tmpBytes) >= 2 && tmpBytes[0] == 0 && tmpBytes[1] >= 0x80 {
		tmpBytes = tmpBytes[1:]
	}

	// Restore leading zeros
	leadingZeros := 0
	for _, c := range s {
		if c == rune(base58Alphabet[0]) {
			leadingZeros++
		} else {
			break
		}
	}

	result := make([]byte, leadingZeros+len(tmpBytes))
	
	// Convert little endian data to big endian
	copy(result[leadingZeros:], reverse(tmpBytes))
	
	return result, nil
}

// Hash performs double SHA256 hash
func Hash(data []byte) Hash256 {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return Hash256(second)
}

// Hash160 performs SHA256 followed by RIPEMD160 (simplified to just SHA256 for this implementation)
func HashRIPEMD160(data []byte) Hash160 {
	// In the original Bitcoin implementation, this would be SHA256 followed by RIPEMD160
	// For simplicity in this Go implementation, we'll use SHA256 and take first 20 bytes
	hash := sha256.Sum256(data)
	var result Hash160
	copy(result[:], hash[:20])
	return result
}

// EncodeBase58Check encodes with 4-byte checksum
func EncodeBase58Check(input []byte) string {
	// Add 4-byte hash check to the end
	payload := make([]byte, len(input))
	copy(payload, input)
	
	hash := Hash(payload)
	payload = append(payload, hash[:4]...)
	
	return EncodeBase58(payload)
}

// DecodeBase58Check decodes and verifies 4-byte checksum
func DecodeBase58Check(s string) ([]byte, error) {
	decoded, err := DecodeBase58(s)
	if err != nil {
		return nil, err
	}
	
	if len(decoded) < 4 {
		return nil, errors.New("decoded data too short")
	}
	
	// Verify checksum
	payload := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]
	
	hash := Hash(payload)
	if !bytesEqual(hash[:4], checksum) {
		return nil, errors.New("checksum mismatch")
	}
	
	return payload, nil
}

// Hash160ToAddress converts a 160-bit hash to a Bitcoin address
func Hash160ToAddress(hash160 Hash160) string {
	// Add 1-byte version number to the front
	payload := make([]byte, 1+len(hash160))
	payload[0] = addressVersion
	copy(payload[1:], hash160[:])
	
	return EncodeBase58Check(payload)
}

// AddressToHash160 converts a Bitcoin address to a 160-bit hash
func AddressToHash160(address string) (Hash160, error) {
	var hash160 Hash160
	
	decoded, err := DecodeBase58Check(address)
	if err != nil {
		return hash160, err
	}
	
	if len(decoded) == 0 {
		return hash160, errors.New("empty decoded data")
	}
	
	version := decoded[0]
	if len(decoded) != len(hash160)+1 {
		return hash160, errors.New("invalid address length")
	}
	
	if version > addressVersion {
		return hash160, errors.New("invalid address version")
	}
	
	copy(hash160[:], decoded[1:])
	return hash160, nil
}

// IsValidBitcoinAddress checks if a string is a valid Bitcoin address
func IsValidBitcoinAddress(address string) bool {
	_, err := AddressToHash160(address)
	return err == nil
}

// PubKeyToAddress converts a public key to a Bitcoin address
func PubKeyToAddress(pubKey []byte) string {
	hash160 := HashRIPEMD160(pubKey)
	return Hash160ToAddress(hash160)
}

// Helper functions

// reverse reverses a byte slice
func reverse(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// bytesEqual compares two byte slices for equality
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Example usage and tests
func main() {
	// Test Base58 encoding/decoding
	testData := []byte("Hello, Bitcoin!")
	encoded := EncodeBase58(testData)
	fmt.Printf("Original: %s\n", testData)
	fmt.Printf("Base58 Encoded: %s\n", encoded)
	
	decoded, err := DecodeBase58(encoded)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return
	}
	fmt.Printf("Decoded: %s\n", decoded)
	
	// Test Base58Check encoding/decoding
	encodedCheck := EncodeBase58Check(testData)
	fmt.Printf("Base58Check Encoded: %s\n", encodedCheck)
	
	decodedCheck, err := DecodeBase58Check(encodedCheck)
	if err != nil {
		fmt.Printf("DecodeCheck error: %v\n", err)
		return
	}
	fmt.Printf("Decoded Check: %s\n", decodedCheck)
	
	// Test address generation
	samplePubKey := []byte("sample public key data for testing")
	address := PubKeyToAddress(samplePubKey)
	fmt.Printf("Generated Address: %s\n", address)
	fmt.Printf("Address is valid: %t\n", IsValidBitcoinAddress(address))
	
	// Test address to hash160 conversion
	hash160, err := AddressToHash160(address)
	if err != nil {
		fmt.Printf("Address to Hash160 error: %v\n", err)
		return
	}
	fmt.Printf("Hash160 from address: %x\n", hash160)
}
