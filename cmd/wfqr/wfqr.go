package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	wifiqr "github.com/g-lib/wifi-qrcode"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	ssid   = flag.String("ssid", "", "wifi ssid")
	auth   = flag.String("auth", "nopass", "auth-type(WAP、WEP、nopass)")
	pwd    = flag.String("pwd", "", "wifi password")
	img    = flag.String("img", "", "wifi icon-image path")
	out    = flag.String("out", "wifi-qr-code.png", "out path")
	hidden = flag.Bool("hiden", false, "is hidden?")
	size   = flag.Int("size", 512, "image size in pixels")
	code   = flag.Bool("code", false, "WIFI code only")
	yes    = flag.Bool("yes", false, "input wifi password directly")
)

func usage() {
	fmt.Fprintf(os.Stderr, `
Generate a qr code for your wifi network

Usage: 
  wfqr <-ssid> [-auth-type] [-pwd] [-img] [-out] 
               [-hidden] [-code] [-size] [-yes]  

Options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
`)
	os.Exit(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
	}
	passwd := *pwd
	if !*yes {
		fmt.Print("Please input wifi password:")
		password, _ := terminal.ReadPassword(0)
		passwd = string(password)
	}
	codeText, err := wifiqr.GenWIFICode(*ssid, *auth, *hidden, passwd)
	if err != nil {
		log.Fatalln(err)
	}
	if *code {
		fmt.Println("wifi-code:", codeText)
		os.Exit(0)
	}
	var logo image.Image
	if *img != "" {
		file, err := os.Open(*img)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		logo, _, err = image.Decode(file)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		logo = nil
	}

	result, err := wifiqr.GenWIFIQRCode(codeText, *size, logo)
	if err != nil {
		log.Fatalln(err)
	}
	outF, err := os.Create(*out)
	if err != nil {
		log.Fatalln(err)
	}
	outF.Write(result.Bytes())
	outF.Close()
	fmt.Println("Done! Written QR image to", *out)
}
