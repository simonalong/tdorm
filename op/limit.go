package op

type Limit struct {
	//    [LIMIT limit_val [OFFSET offset_val]]
}

type Slimit struct {
	//    [SLIMIT limit_val [SOFFSET offset_val]]
}

func (receiver *Limit) GenerateSql() string {
	// todo
	return ""
}

func (receiver *Slimit) GenerateSql() string {
	// todo
	return ""
}
