package writer

import (
	"net/mail"
	"testing"
)

func TestAddressEncoding(t *testing.T) {
	addr := &mail.Address{
		Name:    "技研科-盧詩雲",
		Address: "shihyun2800@mail.cepp.gov.tw",
	}
	t.Logf("Encoded: %s", addr.String())
}
