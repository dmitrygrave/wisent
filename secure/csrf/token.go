package csrf

import (
	"encoding/base64"
	"net/http"

	"github.com/dmitrygrave/wisent/secure/helpers"
)

// getTokenFromRequest tries to find the token in one of three different places
// the header, the post form value, or the multipart form. If it can find it
// there it returns it, otherwise it returns an empty byte slice
func getTokenFromRequest(r *http.Request) ([]byte, error) {
	// Check http header
	token := r.Header.Get(HeaderName)

	if len(token) == 0 {
		token = r.PostFormValue(FormFieldName)
	}

	if len(token) == 0 && r.MultipartForm != nil {
		vals := r.MultipartForm.Value[FormFieldName]
		if len(vals) != 0 {
			token = vals[0]
		}
	}

	return base64.URLEncoding.DecodeString(token)
}

// maskToken generates a masked token of tokenLength * 2 which is the result of
// base64 encoding the xor of tokenLength random bytes and the generated token
func maskToken(realToken []byte, r *http.Request) string {
	opt, err := helpers.GenerateSecureBytes(tokenLength)

	if err != nil {
		return ""
	}

	// XOR the OTP with the real token to generate a masked token.
	return base64.URLEncoding.EncodeToString(append(opt, xorToken(opt, realToken)...))
}

// unmaskToken returns the xor value of the received token which is the result
// of maskToken
func unmaskToken(issued []byte) []byte {
	if len(issued) != tokenLength*2 {
		return nil
	}

	otp := issued[tokenLength:]
	masked := issued[:tokenLength]

	return xorToken(otp, masked)
}

// compareTokens securely compares two byte slices and returns if they are equal
func compareTokens(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	return helpers.SecureCompare(a, b)
}

// xorToken xor's two byte slices and returns the result
func xorToken(a, b []byte) []byte {
	n := len(a)

	if len(b) < n {
		n = len(b)
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = a[i] ^ b[i]
	}

	return res
}
