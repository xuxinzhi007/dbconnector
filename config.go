package dbconnector

type MysqlConfig struct {
	// 方式一：完整的DSN字符串
	DataSource string `json:"dataSource" yaml:"DataSource"`

	// 方式二：分拆的配置项
	Host     string `json:"host" yaml:"Host"`
	Port     int    `json:"port" yaml:"Port"`
	User     string `json:"user" yaml:"User"`
	Password string `json:"password" yaml:"Password"`
	DBName   string `json:"dbName" yaml:"DBName"`
	Charset  string `json:"charset" yaml:"Charset"`
	ParseTime bool   `json:"parseTime" yaml:"ParseTime"`
	Loc      string `json:"loc" yaml:"Loc"`
}
