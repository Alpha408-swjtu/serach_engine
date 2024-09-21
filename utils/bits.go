package utils

func IsBit1(n uint64, i int) bool {
	if i > 64 {
		panic(i)
	}
	c := uint64(1 << (i - 1))
	if n&c == c {
		return true
	} else {
		return false
	}
}

func SetBit1(n uint64, i int) uint64 {
	if i > 64 {
		panic(i)
	}
	c := uint64(1 << (i - 1))
	return n | c
}

func CountBit1(n uint64) int {
	c := uint64(1)
	sum := 0
	for i := 0; i < 64; i++ {
		if c&n == c {
			sum++
		}
		c = c << 1
	}
	return sum
}

// 将document属性编辑到位图
type Candidate struct {
	Id     int
	Gender string
	Vip    bool
	Active int
	Bits   uint64 //存储上方信息到bit中
}

const (
	MALE        = 1
	VIP         = 1 << 1
	WEEK_ACTIVE = 1 << 2
)

func (c *Candidate) SetMale() {
	c.Gender = "男"
	c.Bits = c.Bits | MALE
}

func (c *Candidate) SetVip() {
	c.Vip = true
	c.Bits = c.Bits | VIP
}

func (c *Candidate) SetActive(day int) {
	c.Active = day
	if day <= 7 {
		c.Bits = c.Bits | WEEK_ACTIVE
	}
}

// 判断多个条件是否满足：将条件先编码进入on这个bits
func (c Candidate) Filter2(on uint64) bool {
	return c.Bits&on == on
}
