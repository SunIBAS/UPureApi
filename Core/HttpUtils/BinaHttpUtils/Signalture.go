package BinaHttpUtils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Sign accepts a message and a private key in PEM format,
// signs the message using RSASSA-PKCS1-V1_5 with SHA-256, and returns the base64-encoded signature.
// str = "symbol=BTCUSDT&side=BUY&type=LIMIT&quantity=1&price=9000&timeInForce=GTC&recvWindow=5000&timestamp=1591702613943"
// key = "2b5eb11e18796d12d88f13dc27dbbd02c2cc51ff7059765ed9821957d82bb4d9"
// return "3c661234138461fcc7a7d8746c6558c9842d4e10870d2ecbedf7777cad694af9"
func Sign(str, key string) string {
	// Convert the key and message to byte slices
	keyBytes := []byte(key)
	messageBytes := []byte(str)

	// Create a new HMAC by defining the hash type (SHA-256) and the key
	h := hmac.New(sha256.New, keyBytes)

	// Write the message to the HMAC object
	h.Write(messageBytes)

	// Get the final HMAC result as a byte slice
	hash := h.Sum(nil)

	// Convert the byte slice to a hex string and return it
	return hex.EncodeToString(hash)
}

// https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/general-info#%E7%A4%BA%E4%BE%8B-1-%E6%89%80%E6%9C%89%E5%8F%82%E6%95%B0%E9%80%9A%E8%BF%87-query-string-%E5%8F%91%E9%80%81
