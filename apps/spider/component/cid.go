package component

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// cidTemplate 代表组件 ID 的模板
var cidTemplate = "%s%d|%s"

// CID 代表组件 ID
type CID string

// GenCID 会根据给定参数生成组件ID
func GenCID(mtype Type, sn uint64, maddr net.Addr) (CID, error) {
	if !LegalType(mtype) {
		errMsg := fmt.Sprintf("illegal module type: %s", mtype)
		return "", errors.NewIllegalParameterError(errMsg)
	}
	letter := legalTypeLetterMap[mtype]
	var midStr string
	if maddr == nil {
		midStr = fmt.Sprintf(cidTemplate, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(cidTemplate, letter, sn, maddr.String())
	}
	return CID(midStr), nil
}


// SplitCID 用于分解组件 ID。
// 第一个结果值表示分解是否成功。
// 若分解成功，则第二个结果值长度为3，
// 并依次包含组件类型字母、序列号和组件网络地址（如果有的话）。
func SplitCID(cid CID) ([]string, error) {
	var ok bool
	var letter string
	var snStr string
	var addr string
	cidStr := string(cid)
	if len(cidStr) <= 1 {
		return nil, errors.NewIllegalParameterError("insufficient MID")
	}
	letter = cidStr[:1]
	if _, ok = legalLetterTypeMap[letter]; !ok {
		return nil, errors.NewIllegalParameterError(
			fmt.Sprintf("illegal module type letter: %s", letter))
	}
	snAndAddr := cidStr[1:]
	index := strings.LastIndex(snAndAddr, "|")
	if index < 0 {
		snStr = snAndAddr
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module SN: %s", snStr))
		}
	} else {
		snStr = snAndAddr[:index]
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module SN: %s", snStr))
		}
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module address: %s", addr))
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module IP: %s", ipStr))
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64); err != nil {
			return nil, errors.NewIllegalParameterError(
				fmt.Sprintf("illegal module port: %s", portStr))
		}
	}
	return []string{letter, snStr, addr}, nil
}
