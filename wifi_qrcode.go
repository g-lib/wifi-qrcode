package wifiqr

import (
	_ "embed"
	"errors"
	"fmt"
	"image"
	"strings"

	"bytes"
	"image/color"
	"image/png"

	qr "github.com/skip2/go-qrcode"
)

//go:embed wifi.png
var WIFIIMG []byte

func GenWIFICode(ssid, authType string, hidden bool, pwd ...string) (string, error) {
	isHidden := "false"
	authType = strings.ToUpper(authType)
	if hidden {
		isHidden = "true"
	}
	switch authType {
	case "NOPASS":
		if len(pwd) != 0 && pwd[0] != "" {
			return "", errors.New("for nopass, password should be None")
		}
		return fmt.Sprintf("WIFI:T:nopass;S:%s;H:%s;;", ssid, isHidden), nil
	case "WAP", "WEP":
		if len(pwd) == 0 || pwd[0] == "" {
			return "", errors.New("for WPA and WEP, password should not be None")
		}
		return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;H:%s;;",
			authType, ssid, pwd[0], isHidden), nil
	default:
		return "", errors.New("Unknown authentication_type:" + authType +
			".WPA 、WEP、nopass only")
	}
}

func GenWIFIQRCode(text string, size int, logo ...image.Image) (*bytes.Buffer, error) {
	var img image.Image
	if len(logo) == 0 || logo[0] == nil {
		img, _, _ = image.Decode(bytes.NewReader(WIFIIMG))
	}
	return Encode(text, img, size)
}

// Encoder defines settings for QR/Overlay encoder.
type Encoder struct {
	AlphaThreshold int
	GreyThreshold  int
	QRLevel        qr.RecoveryLevel
}

// DefaultEncoder is the encoder with default settings.
var DefaultEncoder = Encoder{
	AlphaThreshold: 2000,       // FIXME: don't remember where this came from
	GreyThreshold:  30,         // in percent
	QRLevel:        qr.Highest, // recommended, as logo steals some redundant space
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	return DefaultEncoder.Encode(str, logo, size)
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func (e Encoder) Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	code, err := qr.New(str, e.QRLevel)
	if err != nil {
		return nil, err
	}

	img := code.Image(size)
	e.overlayLogo(img, logo)

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// overlayLogo blends logo to the center of the QR code,
// changing all colors to black.
func (e Encoder) overlayLogo(dst, src image.Image) {
	grey := uint32(^uint16(0)) * uint32(e.GreyThreshold) / 100
	alphaOffset := uint32(e.AlphaThreshold)
	offset := dst.Bounds().Max.X/2 - src.Bounds().Max.X/2
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			if r, g, b, alpha := src.At(x, y).RGBA(); alpha > alphaOffset {
				col := color.Black
				if r > grey && g > grey && b > grey {
					col = color.White
				}
				dst.(*image.Paletted).Set(x+offset, y+offset, col)
			}
		}
	}
}
