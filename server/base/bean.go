package base

type OBean struct {
	Text    string      `json:"text" column:"text"`
	Value   interface{} `json:"value" column:"value"`
	Comment string      `json:"comment" column:"comment"`
	Color   string      `json:"color" column:"color"`
}

func NewOBean(text string, value interface{}) (res OBean) {
	res = OBean{
		Text:  text,
		Value: value,
	}
	return
}

type UserTotalBean struct {
	User       *UserEntity         `json:"user"`
	Password   *UserPasswordEntity `json:"password"`
	Persona    *UserPersonaBean    `json:"persona"`
	Enterprise *UserEnterpriseBean `json:"enterprise"`
}

type UserPersonaBean struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Sex    int8    `json:"sex"`
	Photo  string  `json:"photo"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
}

type UserEnterpriseBean struct {
	Name   string                  `json:"name"`
	Salary float32                 `json:"salary"`
	Orgs   []UserEnterpriseOrgBean `json:"orgs"`
}

type UserEnterpriseOrgBean struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Position string `json:"position"`
}

type RequestBean struct {
	Session *SessionBean   `json:"session"`
	User    *LoginUserBean `json:"user"`
}

type SessionBean struct {
	User *LoginUserBean `json:"user"`
}

type LoginUserBean struct {
	ServerId int64  `json:"serverId"`
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Account  string `json:"account"`
	Email    string `json:"email"`
}

type InstallInfo struct {
	Module string              `json:"module"`
	Stages []*InstallStageInfo `json:"stages"`
}

type InstallStageInfo struct {
	Stage    string   `json:"stage"`
	SqlParam SqlParam `json:"sqlParam"`
}

type SqlParam struct {
	Sql    string        `json:"sql,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func NewSqlParam(sql_ string, params []interface{}) (sqlParam SqlParam) {
	if params == nil {
		params = []interface{}{}
	}
	sqlParam = SqlParam{
		Sql:    sql_,
		Params: params,
	}
	return
}
