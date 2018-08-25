package cmap

import "fmt"

type IllegalParameterError struct {
	msg string
}

func newIllegalParameterError(errMsg string)IllegalParameterError{
	return IllegalParameterError{
		msg:fmt.Sprintf("并发安全字典：<非法的参数错误：%s>",errMsg),
	}
}

func (pe IllegalParameterError)Error()string{
	return pe.msg
}

type IllegalPairTypeError struct {
	msg string
}

func newIllegalPairTypeError(pair Pair)IllegalPairTypeError{
	return IllegalPairTypeError{
		msg:fmt.Sprintf("并发安全字典：<Pair类型错误：%T>",pair),
	}
}

func (pte IllegalPairTypeError)Error()string{
	return pte.msg
}

type IllegalRedistributorError struct {
	msg string
}

func newIllegalRedistributorError(errMsg string)IllegalRedistributorError{
	return IllegalRedistributorError{
		msg:fmt.Sprintf("并发安全字典：<键值对再分布器错误:%s>",errMsg),
	}
}

func (re IllegalRedistributorError)Error()string{
	return re.msg
}