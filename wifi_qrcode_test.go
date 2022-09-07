package wifiqr

import "testing"

func TestWIFICode(t *testing.T) {
	t.Log(GenWIFICode("tacey", "nopass", false))
	t.Log(GenWIFICode("tacey", "nopass", false, ""))
	t.Log(GenWIFICode("tacey", "nopass", false, "123456"))
	t.Log(GenWIFICode("tacey", "WAP", false, "123456"))
	t.Log(GenWIFICode("tacey", "WEP", false))
	t.Log(GenWIFICode("tacey", "WxP", false))

}

func TestWIFIQRCode(t *testing.T) {

}
