package main

import (
	"crypto/rand"
	"math/big"
	"strings"
	"syscall/js"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
)

func main() {
	// SM2 非对称加密
	JsSm2 := js.FuncOf(func(this js.Value, args []js.Value) any {
		return this
	})
	JsSm2.Set("prototype", map[string]interface{}{
		"C1C3C2": js.ValueOf(0),
		"C1C2C3": js.ValueOf(1),
		"setSm2Mode": js.FuncOf(func(this js.Value, args []js.Value) any {
			this.Set("sm2Mode", args[0])
			return this
		}),
		"setSm2SignUid": js.FuncOf(func(this js.Value, args []js.Value) any {
			this.Set("sm2SignUid", args[0])
			return this
		}),
		"encrypt": js.FuncOf(func(this js.Value, args []js.Value) any {
			// C1C3C2 C1C2C3
			var mode = this.Get("sm2Mode").Int()
			// 公钥 16进制字符串
			var khex = args[0].String()
			// 从JS读取 Uint8Array
			var in = make([]byte, args[1].Length())
			js.CopyBytesToGo(in, args[1])
			// 读取公钥
			pub, err := x509.ReadPublicKeyFromHex(khex)
			if err != nil {
				panic(err)
			}
			// 加密
			data, err := sm2.Encrypt(pub, in, rand.Reader, mode)
			if err != nil {
				panic(err)
			}
			// 拷贝到JS Uint8Array
			out := js.Global().Get("Uint8Array").New(len(data))
			js.CopyBytesToJS(out, data)
			return out
		}),
		"decrypt": js.FuncOf(func(this js.Value, args []js.Value) any {
			// C1C3C2 C1C2C3
			var mode = this.Get("sm2Mode").Int()
			// 私钥 16进制字符串
			var khex = args[0].String()
			// 从JS读取 Uint8Array
			var in = make([]byte, args[1].Length())
			js.CopyBytesToGo(in, args[1])
			// 读取私钥
			priv, err := x509.ReadPrivateKeyFromHex(khex)
			if err != nil {
				panic(err)
			}
			// 加密
			data, err := sm2.Decrypt(priv, in, mode)
			if err != nil {
				panic(err)
			}
			// 拷贝到JS Uint8Array
			out := js.Global().Get("Uint8Array").New(len(data))
			js.CopyBytesToJS(out, data)
			return out
		}),
		"sign": js.FuncOf(func(this js.Value, args []js.Value) any {
			// 公钥 16进制字符串
			var khex = args[0].String()
			// 从JS读取 Uint8Array
			var in = make([]byte, args[1].Length())
			js.CopyBytesToGo(in, args[1])
			var uid []byte
			if !this.Get("sm2SignUid").IsNull() {
				uid = make([]byte, this.Get("sm2SignUid").Length())
				js.CopyBytesToGo(in, this.Get("sm2SignUid"))
			}
			// 读取公钥
			priv, err := x509.ReadPrivateKeyFromHex(khex)
			if err != nil {
				panic(err)
			}
			// 签名
			r, s, err := sm2.Sm2Sign(priv, in, uid, rand.Reader)
			if err != nil {
				panic(err)
			}
			rb := r.Bytes()
			sb := s.Bytes()
			// 拷贝到JS Uint8Array
			out := js.Global().Get("Uint8Array").New(64)
			js.CopyBytesToJS(out, append(rb, sb...))
			return out
		}),
		"verify": js.FuncOf(func(this js.Value, args []js.Value) any {
			// 私钥 16进制字符串
			var khex = args[0].String()
			// 从JS读取 Uint8Array
			var in = make([]byte, args[1].Length())
			js.CopyBytesToGo(in, args[1])
			var rs = make([]byte, args[2].Length())
			js.CopyBytesToGo(in, args[2])
			var uid []byte
			if !this.Get("sm2SignUid").IsNull() {
				uid = make([]byte, this.Get("sm2SignUid").Length())
				js.CopyBytesToGo(in, this.Get("sm2SignUid"))
			}
			// 读取私钥
			pub, err := x509.ReadPublicKeyFromHex(khex)
			if err != nil {
				panic(err)
			}
			// 验签
			flag := sm2.Sm2Verify(pub, in, uid, (&big.Int{}).SetBytes(rs[:32]),
				(&big.Int{}).SetBytes(rs[32:]))
			return flag
		}),
	})
	js.Global().Set("Sm2", JsSm2)

	// SM3 哈希
	JsSm3 := js.FuncOf(func(this js.Value, args []js.Value) any {
		return this
	})
	JsSm3.Set("prototype", map[string]interface{}{
		"write": js.FuncOf(func(this js.Value, args []js.Value) any {
			this.Set("data", args[0])
			return this
		}),
		"sum": js.FuncOf(func(this js.Value, args []js.Value) any {
			hw := sm3.New()
			// 从JS读取 Uint8Array
			var in = make([]byte, this.Get("data").Length())
			js.CopyBytesToGo(in, this.Get("data"))
			// 计算
			hw.Write(in)
			hash := hw.Sum(nil)
			// 拷贝到JS Uint8Array
			out := js.Global().Get("Uint8Array").New(len(hash))
			js.CopyBytesToJS(out, hash)
			return out
		}),
	})
	js.Global().Set("Sm3", JsSm3)

	// SM4 对称加密
	JsSm4 := js.FuncOf(func(this js.Value, args []js.Value) any {
		this.Set("blockMode", js.ValueOf("ECB"))
		return this
	})
	JsSm4.Set("prototype", map[string]interface{}{
		"setBlockMode": js.FuncOf(func(this js.Value, args []js.Value) any {
			this.Set("blockMode", args[0])
			return this
		}),
		"encrypt": js.FuncOf(func(this js.Value, args []js.Value) any {
			// 分组模式
			var mode = strings.ToUpper(this.Get("blockMode").String())
			// 从JS读取 Uint8Array
			var key = make([]byte, args[0].Length())
			js.CopyBytesToGo(key, args[0])
			var in = make([]byte, args[1].Length())
			js.CopyBytesToGo(in, args[1])
			// 加密
			var data []byte
			var err error
			switch mode {
			default:
				fallthrough
			case "ECB":
				data, err = sm4.Sm4Ecb(key, in, true)
			case "CBC":
				data, err = sm4.Sm4Cbc(key, in, true)
			case "CFB":
				data, err = sm4.Sm4CFB(key, in, true)
			case "OFB":
				data, err = sm4.Sm4OFB(key, in, true)
			}
			if err != nil {
				panic(err)
			}
			// 拷贝到JS Uint8Array
			out := js.Global().Get("Uint8Array").New(len(data))
			js.CopyBytesToJS(out, data)
			return out
		}),
		"decrypt": js.FuncOf(func(this js.Value, args []js.Value) any {
			// 分组模式
			var mode = strings.ToUpper(this.Get("blockMode").String())
			// 从JS读取 Uint8Array
			var key = make([]byte, args[0].Length())
			js.CopyBytesToGo(key, args[0])
			var in = make([]byte, args[1].Length())
			js.CopyBytesToGo(in, args[1])
			// 解密
			var data []byte
			var err error
			switch mode {
			default:
				fallthrough
			case "ECB":
				data, err = sm4.Sm4Ecb(key, in, false)
			case "CBC":
				data, err = sm4.Sm4Cbc(key, in, false)
			case "CFB":
				data, err = sm4.Sm4CFB(key, in, false)
			case "OFB":
				data, err = sm4.Sm4OFB(key, in, false)
			}
			if err != nil {
				panic(err)
			}
			// 拷贝到JS Uint8Array
			out := js.Global().Get("Uint8Array").New(len(data))
			js.CopyBytesToJS(out, data)
			return out
		}),
	})
	js.Global().Set("Sm4", JsSm4)

	<-make(chan bool)
}
