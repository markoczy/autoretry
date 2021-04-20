package main

import (
	"crypto/aes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/markoczy/xtools/common/flags"
	"github.com/markoczy/xtools/common/helpers"
	"github.com/markoczy/xtools/common/logger"
	"golang.org/x/crypto/sha3"
)

type mode string

const (
	modeSha1     = mode("sha")
	modeSha2     = mode("sha2")
	modeSha3     = mode("sha3")
	modeShake128 = mode("shake128")
	modeShake256 = mode("shake256")
	modeAes      = mode("aes")
	modeBase64   = mode("base64")
	modeHex      = mode("hex")
	modeBinary   = mode("binary")

	undefined = "<undefined>"
	// modeSha224 = mode("sha224")
	// modeSha256 = mode("sha256")
	// modeSha384 = mode("sha384")
	// modeSha512 = mode("sha512")
)

func (m *mode) ValidCrypto() error {
	switch *m {
	case modeSha1, modeSha2, modeSha3:
		return nil
	}
	return fmt.Errorf("Invalid mode")
}

var (
	// validShaModes = types.NewEnum([]string{"1", "224", "256", "384", "512"})
	// modeSha1, _   = validShaModes.ValueOf("1")
	// modeSha224, _ = validShaModes.ValueOf("224")
	// modeSha256, _ = validShaModes.ValueOf("256")
	// modeSha384, _ = validShaModes.ValueOf("384")
	// modeSha512, _ = validShaModes.ValueOf("512")

	cryptoMode mode
	formatMode = modeBinary
	log        logger.Logger
	file       string
	key        string
	hashKey    bool
	hashSize   int

	// shaMode, _ = validShaModes.ValueOf("256")
	shaMode = "256"
	aesMode = "256"
)

func initFlags() (err error) {
	sha2 := flags.NewEnum([]string{"224", "256", "384", "512", "512_224", "512_256"}, "")
	sha3 := flags.NewEnum([]string{"224", "256", "384", "512"}, "")
	aes := flags.NewEnum([]string{"128", "192", "256"}, "")
	shake128 := flags.NewSwitchable("")
	shake256 := flags.NewSwitchable("")

	logFactory := logger.NewAutoFlagFactory()

	// flag.IntVar()
	sha1Ptr := flag.Bool("sha1", false, "Switch to sha1 hashing mode")
	flag.Var(shake128, "shake128", "Switch to sha3 shake128 hashing mode and input arbitrary hash length in bytes, `int` > 0")
	flag.Var(shake256, "shake256", "Switch to sha3 shake256 hashing mode and input arbitrary hash length in bytes, `int` > 0")
	flag.Var(sha2, "sha", "Switch to sha2 hashing mode and input size `string`, allowed inputs: '224', '256', '384', '512', '512_224', '512_256'")
	flag.Var(sha2, "sha2", "(Synonym for sha)``")
	flag.Var(sha3, "sha3", "Switch to sha3 hashing mode and input size `string`, allowed inputs: '224', '256', '384', '512'")
	flag.Var(aes, "aes", "Switch aes encryption mode and input size, allowed inputs: '128', '192', '256'")
	base64Ptr := flag.Bool("base64", false, "Format output as base64 instead of binary (cannot be combined with other formatters)")
	hexPtr := flag.Bool("hex", false, "Format output as hexadecimal instead of binary (cannot be combined with other formatters)")
	filePtr := flag.String("file", undefined, "Write output to file (if not defined, output will be written to standard output)")
	keyPtr := flag.String("key", undefined, "Input encryption key for password based encryption like aes")
	hashKeyPtr := flag.Bool("hashkey", false, "Hash the password key to match length to the algo's required key length (using shake256 algorithm)")

	logFactory.InitFlags()
	flag.Parse()

	log = logFactory.Create()

	cryptoModeCnt := 0
	if *sha1Ptr {
		cryptoMode = modeSha1
		cryptoModeCnt++
	}
	if shake128.Defined() {
		cryptoMode = modeShake128
		if hashSize, err = strconv.Atoi(shake128.String()); err != nil {
			return
		}
		cryptoModeCnt++
	}
	if shake256.Defined() {
		cryptoMode = modeShake256
		if hashSize, err = strconv.Atoi(shake256.String()); err != nil {
			return
		}
		cryptoModeCnt++
	}
	if sha2.String() != "" {
		shaMode = sha2.String()
		cryptoMode = modeSha2
		cryptoModeCnt++
	}
	if sha3.String() != "" {
		shaMode = sha3.String()
		cryptoMode = modeSha3
		cryptoModeCnt++
	}
	if aes.String() != "" {
		aesMode = aes.String()
		cryptoMode = modeAes
		cryptoModeCnt++
	}
	if cryptoModeCnt > 1 {
		err = fmt.Errorf("More than one crypto mode selected")
		return
	}

	formatModeCnt := 0
	if *base64Ptr {
		formatMode = modeBase64
		formatModeCnt++
	}
	if *hexPtr {
		formatMode = modeHex
		formatModeCnt++
	}
	if formatModeCnt > 1 {
		err = fmt.Errorf("More than one format mode selected")
		return
	}

	file = *filePtr
	key = *keyPtr
	hashKey = *hashKeyPtr
	log.Debug("Crypte mode:", cryptoMode)
	return
}

func main() {
	initFlags()

	d, err := helpers.ReadStdin()
	check(err)

	var enc []byte
	switch cryptoMode {
	case modeSha1:
		enc = encryptSha1([]byte(d))
	case modeSha2:
		enc = encryptSha2([]byte(d))
	case modeSha3:
		enc = encryptSha3([]byte(d))
	case modeShake128:
		enc = encryptShake128([]byte(d), hashSize)
	case modeShake256:
		enc = encryptShake256([]byte(d), hashSize)
	case modeAes:
		enc, err = encryptAes([]byte(d))
		check(err)
	}

	switch formatMode {
	case modeBase64:
		dst := make([]byte, base64.RawStdEncoding.EncodedLen(len(enc)))
		base64.RawStdEncoding.Encode(dst, enc)
		enc = dst
	case modeHex:
		dst := make([]byte, hex.EncodedLen(len(enc)))
		hex.Encode(dst, enc)
		enc = dst
	default:
		log.Debug("No format mode selected")
		// fmt.Print(string(enc))
	}

	if file == undefined {
		log.Info("Writing output to STDOUT")
		fmt.Print(string(enc))
		return
	}
	log.Info("Writing output to file", file)
	ioutil.WriteFile(file, enc, 644)
}

func encryptSha1(d []byte) (ret []byte) {
	x := sha1.Sum([]byte(d))
	return x[:]
}

func encryptSha2(d []byte) (ret []byte) {
	switch shaMode {
	case "224":
		x := sha256.Sum224([]byte(d))
		return x[:]
	case "256":
		x := sha256.Sum256([]byte(d))
		return x[:]
	case "384":
		x := sha512.Sum384([]byte(d))
		return x[:]
	case "512":
		x := sha512.Sum512([]byte(d))
		return x[:]
	case "512_224":
		x := sha512.Sum512_224([]byte(d))
		return x[:]
	case "512_256":
		x := sha512.Sum512_256([]byte(d))
		return x[:]
	default:
		log.Error("Invalid SHA2 Mode:", shaMode)
		return nil
	}
}

func encryptSha3(d []byte) (ret []byte) {
	switch shaMode {
	case "224":
		x := sha3.Sum224([]byte(d))
		return x[:]
	case "256":
		x := sha3.Sum256([]byte(d))
		return x[:]
	case "384":
		x := sha3.Sum384([]byte(d))
		return x[:]
	case "512":
		x := sha3.Sum512([]byte(d))
		return x[:]
	// case "shake128":
	// 	x := sha3.ShakeSum128([]byte(d))
	// 	return x[:]
	// case "shake256":
	// 	x := sha3.ShakeSum256([]byte(d))
	// 	return x[:]
	default:
		log.Error("Invalid SHA3 Mode:", shaMode)
		return nil
	}
}

func encryptShake128(d []byte, size int) (ret []byte) {
	dst := make([]byte, size)
	sha3.ShakeSum128(dst, d)
	return dst
}

func encryptShake256(d []byte, size int) (ret []byte) {
	dst := make([]byte, size)
	sha3.ShakeSum256(dst, d)
	return dst
}

func encryptAes(d []byte) ([]byte, error) {
	var cipherKey []byte
	if hashKey {
		switch aesMode {
		case "128":
			log.Info("AES Mode: 128")
			cipherKey = encryptShake256([]byte(key), 16)
		case "192":
			log.Info("AES Mode: 192")
			cipherKey = encryptShake256([]byte(key), 24)
		case "256":
			log.Info("AES Mode: 256")
			cipherKey = encryptShake256([]byte(key), 32)
		}
	} else {
		cipherKey = []byte(key)
	}
	log.Info("AES Key size: %d", len(cipherKey))
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(d))
	c.Encrypt(dst, d)

	// c.$
	// switch aesMode {
	// case "128":
	// 	x := sha1.Sum([]byte(d))
	// case "192":
	// 	x := sha256.Sum224([]byte(d))
	// case "256":
	// 	x := sha256.Sum256([]byte(d))
	// default:
	// 	log.Error("Invalid AES Mode:", aesMode)
	// }
	return dst, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
