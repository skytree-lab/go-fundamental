package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"testing"
)

func Test_rsa_gen(t *testing.T) {
	priStr, pubStr, priPem, pubPem, err := GenRsaKey(1024)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("private string:")
	fmt.Println(priStr)
	fmt.Println("public string:")
	fmt.Println(pubStr)
	fmt.Println("private pem string:")
	fmt.Println(priPem)
	fmt.Println("public pem string:")
	fmt.Println(pubPem)
	fmt.Println("trim pem string:")
	base64ed := RemovePrefixAndBase64(pubPem)
	unbase64edbyte, err := base64.StdEncoding.DecodeString(base64ed)
	if err != nil {
		fmt.Println(err)
		return
	}
	unbase64ed := string(unbase64edbyte)
	fmt.Println("-----BEGIN RSA PUBLIC KEY-----" + unbase64ed + "-----END RSA PUBLIC KEY-----\n")
}

func Test_rsa_encrypto_decrypto(t *testing.T) {
	priStr, _, _, pubPem, err := GenRsaKey(1024)
	if err != nil {
		fmt.Println(err.Error())
	}

	var pubByte []byte
	p, rest := pem.Decode([]byte(pubPem))
	if p != nil && p.Bytes != nil {
		pubByte = p.Bytes
	} else {
		pubByte = rest
	}
	pubStr := hex.EncodeToString(pubByte)

	plainText := "Hello World"
	encrypto, err := RsaEncrypt([]byte(plainText), pubStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(encrypto))
	fmt.Println("decrypto text")
	plainBytes, err := RsaDecrypt(encrypto, priStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(plainBytes))
}

func Test_rsa_encrypto_decrypto1(t *testing.T) {
	pubStr := "30819f300d06092a864886f70d010101050003818d0030818902818100b044128d458c728599f9d70d13c0068ce70feeb57c3a37589720664d87dfa4b3c9ca89e8c321e94503294b0c14eaad852922965a015cc96f43ec41a65893b403556d01bfe6de18e6cae3b029009092ce38ff62515c986f1e8a81dec98bbd92c621ba60a61c8c8e4115938f470359a723c5af3cf7c72c15fce325c181eb7584a90203010001"
	plainText := "0x5ff0f6f2b918c86814967d7efd01929fcbff25d44c054244300069005d5d834f"
	encrypto, err := RsaEncrypt([]byte(plainText), pubStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	encryptoStr := hex.EncodeToString(encrypto)
	fmt.Println(encryptoStr)

	fmt.Println("decrypto text")
	priStr := "3082025d02010002818100b044128d458c728599f9d70d13c0068ce70feeb57c3a37589720664d87dfa4b3c9ca89e8c321e94503294b0c14eaad852922965a015cc96f43ec41a65893b403556d01bfe6de18e6cae3b029009092ce38ff62515c986f1e8a81dec98bbd92c621ba60a61c8c8e4115938f470359a723c5af3cf7c72c15fce325c181eb7584a902030100010281810084d6685ccb233353785d0f32adc5b3aa10a0b33756add0f414c0b81889e3838e72ef0ecbb9f31e87820066ac6e3f06122a940457445b96fb2167357a959a8ab8100baca37c1900822311da9c48e9f26282d1508f586fb78cfd527a23282ab71cf293c1d5ec7343bf068fef33c6a4b747c6c69cbb55221fcebb555015390a6c19024100dc88c739978906f4ff8b6205ebddbdb0eeee461e54a7995fdae1193277aa6ce6266b9bf7ba0c36bf36c3805d11bbf5256f353ad6df1a289358c147daa9e8a757024100cc9cccd215c1208633f1538c3060372084764a465ea3f2678d243b5aeb92db3255a330d0c4d4de7785e4ffef9030797bbaf7971da1a49e52702ce9f8a5d3b3ff024100c84368eef29ddb7475eea3c80ec560f1a0373df3631a831bd98e99ac0ba0f69d14fc99389f7961e9c81846a3bd6bfa94d0e4fc968d289afa1b1a015f1ef607a702402fa56fb8982241cd9e78dac8b157265f271958906c6767022006c8df922dbf674833d9213444918d699b7ad1b154e8651c939d17e4552e1cea4c3b2b9089ecc7024079989be293ed322539cfbc5fb78a010044788bc1cc91c7399e0eb5cd4a0537ea6060d1a08873b82a28dc6d0969d7454a18ada3452b88c929b4762fd11ecdb9a3"
	plainBytes, err := RsaDecrypt(encrypto, priStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(plainBytes))
	fmt.Println("######################################")
	//#####################################
	plainText = "0;0xdd36720d5086ffc3d2174d84e0780b06c8d5958a8a77bf5c6b488060814f4a83"
	encrypto, err = RsaEncrypt([]byte(plainText), pubStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	encryptoStr = hex.EncodeToString(encrypto)
	fmt.Println(encryptoStr)
	fmt.Println("decrypto text")
	priStr = "3082025d02010002818100b044128d458c728599f9d70d13c0068ce70feeb57c3a37589720664d87dfa4b3c9ca89e8c321e94503294b0c14eaad852922965a015cc96f43ec41a65893b403556d01bfe6de18e6cae3b029009092ce38ff62515c986f1e8a81dec98bbd92c621ba60a61c8c8e4115938f470359a723c5af3cf7c72c15fce325c181eb7584a902030100010281810084d6685ccb233353785d0f32adc5b3aa10a0b33756add0f414c0b81889e3838e72ef0ecbb9f31e87820066ac6e3f06122a940457445b96fb2167357a959a8ab8100baca37c1900822311da9c48e9f26282d1508f586fb78cfd527a23282ab71cf293c1d5ec7343bf068fef33c6a4b747c6c69cbb55221fcebb555015390a6c19024100dc88c739978906f4ff8b6205ebddbdb0eeee461e54a7995fdae1193277aa6ce6266b9bf7ba0c36bf36c3805d11bbf5256f353ad6df1a289358c147daa9e8a757024100cc9cccd215c1208633f1538c3060372084764a465ea3f2678d243b5aeb92db3255a330d0c4d4de7785e4ffef9030797bbaf7971da1a49e52702ce9f8a5d3b3ff024100c84368eef29ddb7475eea3c80ec560f1a0373df3631a831bd98e99ac0ba0f69d14fc99389f7961e9c81846a3bd6bfa94d0e4fc968d289afa1b1a015f1ef607a702402fa56fb8982241cd9e78dac8b157265f271958906c6767022006c8df922dbf674833d9213444918d699b7ad1b154e8651c939d17e4552e1cea4c3b2b9089ecc7024079989be293ed322539cfbc5fb78a010044788bc1cc91c7399e0eb5cd4a0537ea6060d1a08873b82a28dc6d0969d7454a18ada3452b88c929b4762fd11ecdb9a3"
	plainBytes, err = RsaDecrypt(encrypto, priStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(plainBytes))
}

func Test_unbase64(t *testing.T) {
	pubStr := "30819f300d06092a864886f70d010101050003818d0030818902818100b044128d458c728599f9d70d13c0068ce70feeb57c3a37589720664d87dfa4b3c9ca89e8c321e94503294b0c14eaad852922965a015cc96f43ec41a65893b403556d01bfe6de18e6cae3b029009092ce38ff62515c986f1e8a81dec98bbd92c621ba60a61c8c8e4115938f470359a723c5af3cf7c72c15fce325c181eb7584a90203010001"
	pubBytes, err := hex.DecodeString(pubStr)
	if err != nil {
		log.Fatalln(err)
	}
	block := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	}
	pub := pem.EncodeToMemory(&block)
	fmt.Println(string(pub))

	priStr := "3082025d02010002818100b044128d458c728599f9d70d13c0068ce70feeb57c3a37589720664d87dfa4b3c9ca89e8c321e94503294b0c14eaad852922965a015cc96f43ec41a65893b403556d01bfe6de18e6cae3b029009092ce38ff62515c986f1e8a81dec98bbd92c621ba60a61c8c8e4115938f470359a723c5af3cf7c72c15fce325c181eb7584a902030100010281810084d6685ccb233353785d0f32adc5b3aa10a0b33756add0f414c0b81889e3838e72ef0ecbb9f31e87820066ac6e3f06122a940457445b96fb2167357a959a8ab8100baca37c1900822311da9c48e9f26282d1508f586fb78cfd527a23282ab71cf293c1d5ec7343bf068fef33c6a4b747c6c69cbb55221fcebb555015390a6c19024100dc88c739978906f4ff8b6205ebddbdb0eeee461e54a7995fdae1193277aa6ce6266b9bf7ba0c36bf36c3805d11bbf5256f353ad6df1a289358c147daa9e8a757024100cc9cccd215c1208633f1538c3060372084764a465ea3f2678d243b5aeb92db3255a330d0c4d4de7785e4ffef9030797bbaf7971da1a49e52702ce9f8a5d3b3ff024100c84368eef29ddb7475eea3c80ec560f1a0373df3631a831bd98e99ac0ba0f69d14fc99389f7961e9c81846a3bd6bfa94d0e4fc968d289afa1b1a015f1ef607a702402fa56fb8982241cd9e78dac8b157265f271958906c6767022006c8df922dbf674833d9213444918d699b7ad1b154e8651c939d17e4552e1cea4c3b2b9089ecc7024079989be293ed322539cfbc5fb78a010044788bc1cc91c7399e0eb5cd4a0537ea6060d1a08873b82a28dc6d0969d7454a18ada3452b88c929b4762fd11ecdb9a3"
	priBytes, err := hex.DecodeString(priStr)
	if err != nil {
		log.Fatalln(err)
	}
	block = pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: priBytes,
	}

	pri := pem.EncodeToMemory(&block)
	fmt.Println(string(pri))

	base64ed := `ioyjaGFL7tmPTb1JJBjWvQxqdHFT+wdkLI96ARHoKvuKi29pXlY3w8gpQet74vYFXaG+m6TpzqHnDV82Zl05Yo6lD/c0ga5LE0JNiFBRat04eJqqFEXvoH+2tMTqDbqt+3UeF6TCh3mzG/l/0sXDHv6ICGtBr7m748XX65F/2Yc=`
	unbased, err := base64.StdEncoding.DecodeString(base64ed)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(hex.EncodeToString(unbased))

	plainBytes, err := RsaDecrypt(unbased, priStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(plainBytes))
}

func Test_unbase642(t *testing.T) {
	priStr := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQC2bJ242Nbcr/NAdLzhiYO1ZwIMmtPfuQ5Ud122iTd8lNhWLze+
kAMZ8HfD9YuwEMkhREgFw6AP4I2hKFu1xmLHYq41d4qmMJGw5U8WW89jxwHmCYjn
pyFhKwptlKaELlvntk26w6hNwTm9zOtdrPlTdF5pY+83gxwUuXsKs0+9awIDAQAB
AoGATytQP9qLVPdfYEfJeTAO4vkzYHkJHnzXhTzQhNmEJSMLyEtxq9ggIG/F6lt0
JhMC1tclUqsVmg5gKOK0+myAjT9VtcE/K3SGi2Ha5Tmdf4yvMV1n3gn0ciTyAn+a
ysMB4U7AdZyMELXydrpm98NzzPDY4CqVpsm+TtY0z2NlVpECQQDpvFqsn0jlH3cr
FjhgNztMb1+fOtOW6WGQFiq0wCFU6gEFCUgVdWDJfonKClFnXLnIBXe6Qm3u6Diy
zYnc2/THAkEAx80IHKf117I9ZWoS2Ubjfsct3DnHEoTBk9L7O19H/i4okopz5UAl
PYx3DyissqtEnC8G+0J1OqcsEwetOkhGPQJBAM2isruQX5eE+Vjk8bBkpzz191Cl
dqYcqYMUnXnsVubYARGnOs32OLLNKvif30Jqp+A8PzuVWhaII+nbPuEaet0CQF7g
tRZuh+aqhE8xRMwaxpMAYH2wjiD0G7xnttel0lkYcDnON2TIxiZpGmx/Yyrg+46z
IjyUfXihw+MxhQglI/0CQFsfksqKrTtcaW6W0jSUqAGG5OB6FcJPqkmDgAZkgbCX
rxZM/Ozm2mGIDqLo0QCFmCbpOn77WMkKTLjWfkB47Sc=
-----END RSA PRIVATE KEY-----`

	cipherBase64 := `nm47m9TRK7wnjyxDgMSwW+2wD96XJgJUk6dmbKfpvcNgbswQzOUFCXfPap2EkRLCIi8bkhHcMO+mQTVPeMtnzUY6ouHN/mpOg9Vae5t7dce12h50VjC+udtRu51pXriZfL5YUXtCqs30rVi8lp5xKWC0uXlFP+9gD8dqz1Du6v8=`
	copher, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		log.Fatalln(err)
	}

	var priBytes []byte
	p, rest := pem.Decode([]byte(priStr))
	if p != nil && p.Bytes != nil {
		priBytes = p.Bytes
	} else {
		priBytes = rest
	}

	plainBytes, err := RsaDecrypt(copher, hex.EncodeToString(priBytes))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(plainBytes))
}
