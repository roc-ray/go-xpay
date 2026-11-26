package xpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// PaySig 计算支付请求签名 pay_sig。
//
//		pay_sig = hex(hmac_sha256(app_key, uri + '&' + post_body))
//
//	 uri 为不带 query 的接口路径，如 /xpay/query_user_balance；
//	 post_body 必须与实际发送的 HTTP 请求体完全一致；
//	 app_key 为支付环境(env)对应的 AppKey，获取自小程序后台「虚拟支付-支付基础配置」。
//
// 该函数同时也适用于校验虚拟支付消息推送（如 xpay_goods_deliver_notify）的签名。
func PaySig(appKey, uri, postBody string) string {
	return hex.EncodeToString(hmacSHA256([]byte(appKey), uri+"&"+postBody))
}

// UserSignature 计算用户态签名 signature。
//
//	signature = hex(hmac_sha256(session_key, post_body))
//
// session_key 为用户当前有效的登录态密钥，通过 auth.code2Session 获取。
func UserSignature(sessionKey, postBody string) string {
	return hex.EncodeToString(hmacSHA256([]byte(sessionKey), postBody))
}

func hmacSHA256(key []byte, msg string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(msg))
	return h.Sum(nil)
}
