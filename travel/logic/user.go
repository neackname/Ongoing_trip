package logic

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"travel/TravelModel"
)

// DecryptUserInfo @title DecryptUserInfo
// @description	解析密文，获取详细的用户信息
// @auth	Snactop		2023-11-13	19:27
// @param	userInfo *TravelModel.TraUser, sessionKey, encryptedData, iv string
// @return	error	传出错误信息
func DecryptUserInfo(userInfo *TravelModel.TraUser, sessionKey, encryptedData, iv string) error {
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return err
	}

	encryptedDataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return err
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(sessionKeyBytes)
	if err != nil {
		return err
	}

	if len(encryptedDataBytes)%aes.BlockSize != 0 {
		return errors.New("encryptedData is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	decrypted := make([]byte, len(encryptedDataBytes))
	mode.CryptBlocks(decrypted, encryptedDataBytes)

	decrypted = pkcs7Unpad(decrypted)
	if decrypted == nil {
		return errors.New("decryption failed")
	}

	if err := json.Unmarshal(decrypted, &userInfo); err != nil {
		return err
	}

	return nil
}

// pkcs7Unpad @title pkcs7Unpad
// @description  PKCS7解码函数
// @auth	Snactop		2023-11-13	19:27
// @param	data []byte
// @return	[]byte
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return nil
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil
	}
	return data[:(length - unpadding)]
}
