package test

import "time"

// OrmChinaDomain1 原始名的话，则直接前缀转小写跟数据库字段对应
type OrmChinaDomain1 struct {
	Ts      time.Time
	Name    string
	Age     int
	Address string
}

type OrmChinaDomainTag struct {
	Station string `json:"station"`
}

// OrmChinaDomain2 使用json标签
type OrmChinaDomain2 struct {
	TsJson      time.Time `json:"ts"`
	NameJson    string    `json:"name,omitempty"`
	AgeJson     int       `json:"age,omitempty"`
	AddressJson string    `json:"address,omitempty"`
}

// OrmChinaDomain3 使用column标签
type OrmChinaDomain3 struct {
	Timestamp time.Time `column:"ts"`
	Na        string    `column:"name"`
	Ag        int       `column:"age"`
	Add       string    `column:"address"`
}
