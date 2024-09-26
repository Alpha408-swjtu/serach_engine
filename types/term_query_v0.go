package types

import "strings"

type TermQueryV0 struct {
	Should  []TermQueryV0
	Must    []TermQueryV0
	Keyword string
}

func (exp TermQueryV0) Empty() bool {
	return len(exp.Keyword) == 0 && len(exp.Must) == 0 && len(exp.Should) == 0
}

func KeywordExpression(keyword string) TermQueryV0 {
	return TermQueryV0{Keyword: keyword}
}

func ShouldExpression(exps ...TermQueryV0) TermQueryV0 {
	if len(exps) == 0 {
		return TermQueryV0{}
	}
	arr := make([]TermQueryV0, 0, len(exps))
	for _, ele := range exps {
		if !ele.Empty() {
			arr = append(arr, ele)
		}
	}
	return TermQueryV0{Should: arr}
}

func MustExpression(exps ...TermQueryV0) TermQueryV0 {
	if len(exps) == 0 {
		return TermQueryV0{}
	}
	arr := make([]TermQueryV0, 0, len(exps))
	for _, ele := range exps {
		if !ele.Empty() {
			arr = append(arr, ele)
		}
	}
	return TermQueryV0{Must: arr}
}

func (exp TermQueryV0) String() string {
	if len(exp.Keyword) != 0 {
		return exp.Keyword
	} else if len(exp.Must) > 0 {
		if len(exp.Must) == 1 {
			return exp.Must[0].String()
		} else {
			sb := strings.Builder{}
			sb.WriteString("(")
			for _, ele := range exp.Must {
				s := ele.String()
				sb.WriteString(s)
				sb.WriteString("&")
			}
			s := sb.String()
			return s[0:len(s)-1] + ")"
		}
	} else if len(exp.Should) > 0 {
		if len(exp.Should) == 1 {
			return exp.Should[0].String()
		} else {
			sb := strings.Builder{}
			sb.WriteString("(")
			for _, ele := range exp.Should {
				s := ele.String()
				sb.WriteString(s)
				sb.WriteString("|")
			}
			s := sb.String()
			return s[0:len(s)-1] + ")"
		}
	}
	return ""
}
