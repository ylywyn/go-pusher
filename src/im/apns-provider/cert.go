/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : cert.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/pkcs12"
)

//
var (
	ErrFailedToDecryptKey           = errors.New("failed to decrypt private key")
	ErrFailedToParsePKCS1PrivateKey = errors.New("failed to parse PKCS1 private key")
	ErrFailedToParseCertificate     = errors.New("failed to parse certificate PEM data")
	ErrNoPrivateKey                 = errors.New("no private key")
	ErrNoCertificate                = errors.New("no certificate")
)

func LoadP12File(filename string, passwd string) (tls.Certificate, error) {
	p12, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return decodeP12Bytes(p12, passwd)
}

func decodeP12Bytes(bytes []byte, passwd string) (tls.Certificate, error) {
	key, cert, err := pkcs12.Decode(bytes, passwd)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}, nil
}

func LoadPemFile(filename string, password string) (tls.Certificate, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return decodePemBytes(bytes, password)
}

func decodePemBytes(bytes []byte, password string) (tls.Certificate, error) {
	var cert tls.Certificate
	var block *pem.Block
	for {
		block, bytes = pem.Decode(bytes)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		}
		if block.Type == "PRIVATE KEY" || strings.HasSuffix(block.Type, "PRIVATE KEY") {
			key, err := unencryptPrivateKey(block, password)
			if err != nil {
				return tls.Certificate{}, err
			}
			cert.PrivateKey = key
		}
	}
	if len(cert.Certificate) == 0 {
		return tls.Certificate{}, ErrNoCertificate
	}
	if cert.PrivateKey == nil {
		return tls.Certificate{}, ErrNoPrivateKey
	}
	if c, e := x509.ParseCertificate(cert.Certificate[0]); e == nil {
		cert.Leaf = c
	}
	return cert, nil
}

func unencryptPrivateKey(block *pem.Block, password string) (crypto.PrivateKey, error) {
	if x509.IsEncryptedPEMBlock(block) {
		bytes, err := x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return nil, ErrFailedToDecryptKey
		}
		return parsePrivateKey(bytes)
	}
	return parsePrivateKey(block.Bytes)
}

func parsePrivateKey(bytes []byte) (crypto.PrivateKey, error) {
	key, err := x509.ParsePKCS1PrivateKey(bytes)
	if err != nil {
		return nil, ErrFailedToParsePKCS1PrivateKey
	}
	return key, nil
}
