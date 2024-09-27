package types

import "strings"

//最终版本的term_query结构

type TermQuery struct {
	Must    []*TermQuery
	Should  []*TermQuery
	Keyword *Keyword
}

func NewTermQuery(field, word string) *TermQuery {
	return &TermQuery{Keyword: &Keyword{Field: field, Word: word}}
}

func (q *TermQuery) Empty() bool {
	return len(q.Must) == 0 && len(q.Should) == 0 && q.Keyword == nil
}

func (q *TermQuery) And(querys ...*TermQuery) *TermQuery {
	if len(querys) == 0 {
		return q
	}
	arr := make([]*TermQuery, 0, len(querys)+1)
	if !q.Empty() {
		arr = append(arr, q)
	}
	for _, ele := range querys {
		if !ele.Empty() {
			arr = append(arr, ele)
		}
	}
	return &TermQuery{Must: arr}
}

func (q *TermQuery) Or(querys ...*TermQuery) *TermQuery {
	if len(querys) == 0 {
		return q
	}
	arr := make([]*TermQuery, 0, len(querys)+1)
	if !q.Empty() {
		arr = append(arr, q)
	}
	for _, ele := range querys {
		if !ele.Empty() {
			arr = append(arr, ele)
		}
	}
	return &TermQuery{Should: arr}
}

func (q *TermQuery) ToString() string {
	if q.Keyword != nil {
		return q.Keyword.Tostring()
	} else if len(q.Must) != 0 {
		if len(q.Must) == 1 {
			return q.Must[0].ToString()
		} else {
			sb := strings.Builder{}
			sb.WriteByte('(')
			for _, ele := range q.Must {
				s := ele.ToString()
				if len(s) > 0 {
					sb.WriteString(s)
					sb.WriteByte('&')
				}
			}
			s := sb.String()
			return s[0:len(s)-1] + ")"
		}
	} else if len(q.Should) != 0 {
		if len(q.Should) == 1 {
			return q.Should[0].ToString()
		} else {
			sb := strings.Builder{}
			sb.WriteByte('(')
			for _, ele := range q.Should {
				s := ele.ToString()
				sb.WriteString(s)
				sb.WriteByte('|')
			}
			s := sb.String()
			return s[0:len(s)-1] + ")"
		}
	}
	return ""
}
