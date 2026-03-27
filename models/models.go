package models

type Article struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Cid        int    `gorm:"column:cid" json:"cid"`
	Did        int    `gorm:"column:did" json:"did"`
	Title      string `gorm:"column:title;size:255" json:"title"`
	Desc       string `gorm:"column:desc;size:255" json:"desc"`
	Keyword    string `gorm:"column:keyword;size:255" json:"keyword"`
	Author     string `gorm:"column:author;size:255" json:"author"`
	Thumb      string `gorm:"column:thumb;size:255" json:"thumb"`
	Summary    string `gorm:"column:summary;type:text" json:"summary"`
	Content    string `gorm:"column:content;type:longtext" json:"content"`
	IsHot      int    `gorm:"column:is_hot" json:"is_hot"`
	IsNew      int    `gorm:"column:is_new" json:"is_new"`
	IsRecom    int    `gorm:"column:is_recom" json:"is_recom"`
	Weight     int    `gorm:"column:weight" json:"weight"`
	PublicTime int    `gorm:"column:public_time" json:"public_time"`
	Status     int    `gorm:"column:status" json:"status"`
	Month      string `gorm:"column:month;size:255" json:"month"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
}

func (Article) TableName() string {
	return "sl_article"
}

type Category struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Type       string `gorm:"column:type;size:16" json:"type"`
	Parent     string `gorm:"column:parent;size:255" json:"parent"`
	Mark       string `gorm:"column:mark;size:60" json:"mark"`
	Author     string `gorm:"column:author;size:32" json:"author"`
	Name       string `gorm:"column:name;size:100" json:"name"`
	Weight     int    `gorm:"column:weight" json:"weight"`
	Status     int    `gorm:"column:status" json:"status"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
}

func (Category) TableName() string {
	return "sl_category"
}

type Comment struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Aid        int    `gorm:"column:aid" json:"aid"`
	Referer    string `gorm:"column:referer;size:255" json:"referer"`
	Ip         string `gorm:"column:ip;size:50" json:"ip"`
	Name       string `gorm:"column:name;size:255" json:"name"`
	Email      string `gorm:"column:email;size:255" json:"email"`
	Url        string `gorm:"column:url;size:255" json:"url"`
	IsReply    int    `gorm:"column:is_reply" json:"is_reply"`
	Content    string `gorm:"column:content;type:text" json:"content"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
}

func (Comment) TableName() string {
	return "sl_comment"
}

type Menu struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Type       int    `gorm:"column:type" json:"type"`
	Mark       string `gorm:"column:mark;size:60" json:"mark"`
	Author     string `gorm:"column:author;size:32" json:"author"`
	Name       string `gorm:"column:name;size:100" json:"name"`
	Url        string `gorm:"column:url;size:256" json:"url"`
	Parent     int    `gorm:"column:parent" json:"parent"`
	Icon       string `gorm:"column:icon;size:128" json:"icon"`
	Weight     int    `gorm:"column:weight" json:"weight"`
	Status     int    `gorm:"column:status" json:"status"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
}

func (Menu) TableName() string {
	return "sl_menu"
}

type Tags struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Type       string `gorm:"column:type;size:16" json:"type"`
	Mark       string `gorm:"column:mark;size:60" json:"mark"`
	Author     string `gorm:"column:author;size:32" json:"author"`
	Name       string `gorm:"column:name;size:100" json:"name"`
	Weight     int    `gorm:"column:weight" json:"weight"`
	Status     int    `gorm:"column:status" json:"status"`
	IsHot      int    `gorm:"column:is_hot" json:"is_hot"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
}

func (Tags) TableName() string {
	return "sl_tags"
}

type Directory struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Cid        int    `gorm:"column:cid" json:"cid"`
	Type       string `gorm:"column:type;size:16" json:"type"`
	Parent     string `gorm:"column:parent;size:255" json:"parent"`
	Mark       string `gorm:"column:mark;size:60" json:"mark"`
	Author     string `gorm:"column:author;size:32" json:"author"`
	Name       string `gorm:"column:name;size:100" json:"name"`
	Weight     int    `gorm:"column:weight" json:"weight"`
	Status     int    `gorm:"column:status" json:"status"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
}

func (Directory) TableName() string {
	return "sl_directory"
}

type ArticleTags struct {
	ID         int `gorm:"primaryKey;column:id" json:"id"`
	Aid        int `gorm:"column:aid" json:"aid"`
	Tid        int `gorm:"column:tid" json:"tid"`
	CreateTime int `gorm:"column:create_time" json:"create_time"`
	UpdateTime int `gorm:"column:update_time" json:"update_time"`
}

func (ArticleTags) TableName() string {
	return "sl_article_tags"
}

type Website struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Name       string `gorm:"column:name;size:255" json:"name"`
	Key        string `gorm:"column:key;size:32" json:"key"`
	Value      string `gorm:"column:value;type:text" json:"value"`
	Staus      int    `gorm:"column:staus" json:"staus"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
}

func (Website) TableName() string {
	return "sl_website"
}

type Admin struct {
	ID             int    `gorm:"primaryKey;column:id" json:"id"`
	Username       string `gorm:"column:username;size:64" json:"username"`
	Password       string `gorm:"column:password;size:64" json:"password"`
	Salt           string `gorm:"column:salt;size:64" json:"salt"`
	GroupID        int    `gorm:"column:group_id" json:"group_id"`
	Name           string `gorm:"column:name;size:64" json:"name"`
	UpdateTime     int    `gorm:"column:update_time" json:"update_time"`
	CreateTime     int    `gorm:"column:create_time" json:"create_time"`
	ExpirationTime int    `gorm:"column:expiration_time" json:"expiration_time"`
	LoginNum       int    `gorm:"column:login_num" json:"login_num"`
	LastLoginTime  int    `gorm:"column:last_login_time" json:"last_login_time"`
	LastLoginIP    string `gorm:"column:last_login_ip;size:64" json:"last_login_ip"`
	Lang           string `gorm:"column:lang;size:255" json:"lang"`
	Status         int    `gorm:"column:status" json:"status"`
}

func (Admin) TableName() string {
	return "sl_admin"
}

type AdminGroup struct {
	ID          int    `gorm:"primaryKey;column:id" json:"id"`
	Mark        string `gorm:"column:mark;size:50" json:"mark"`
	Name        string `gorm:"column:name;size:64" json:"name"`
	Description string `gorm:"column:description;size:256" json:"description"`
	MenuPermit  string `gorm:"column:menu_permit;type:text" json:"menu_permit"`
	MenuModules string `gorm:"column:menu_modules;type:text" json:"menu_modules"`
	AllowIP     string `gorm:"column:allow_ip;type:text" json:"allow_ip"`
	Status      int    `gorm:"column:status" json:"status"`
	CreateTime  int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime  int    `gorm:"column:update_time" json:"update_time"`
}

func (AdminGroup) TableName() string {
	return "sl_admin_group"
}
