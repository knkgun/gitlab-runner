//go:build boringcrypto
// +build boringcrypto

package boring

import "fmt"

func CheckBoring() {
	fmt.Println("BoringSSL enabled")
}
