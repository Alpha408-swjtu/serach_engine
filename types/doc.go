package types

func (kw Keyword) Tostring() string { //拼接field和keyword
	if len(kw.Word) > 0 {
		return kw.Field + "\001" + kw.Word
	} else {
		return ""
	}
}
