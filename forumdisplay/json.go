package forumdisplay

type thread struct {
	Charset   string          `json:"Charset"`
	Variables threadVariables `json:"Variables"`
	Version   string          `json:"Version"`
}

type threadVariables struct {
	ForumThreadlist []threadVariablesForumThreadlist `json:"forum_threadlist"`
}

type threadVariablesForumThreadlist struct {
	Attachment   string `json:"attachment"`
	Author       string `json:"author"`
	Authorid     string `json:"authorid"`
	Dateline     string `json:"dateline"`
	Dbdateline   string `json:"dbdateline"`
	Dblastpost   string `json:"dblastpost"`
	Digest       string `json:"digest"`
	Displayorder string `json:"displayorder"`
	Lastpost     string `json:"lastpost"`
	Lastposter   string `json:"lastposter"`
	Price        string `json:"price"`
	Readperm     string `json:"readperm"`
	Recommend    string `json:"recommend"`
	RecommendAdd string `json:"recommend_add"`
	Replies      string `json:"replies"`
	Replycredit  string `json:"replycredit"`
	Rushreply    string `json:"rushreply"`
	Special      string `json:"special"`
	Subject      string `json:"subject"`
	Tid          string `json:"tid"`
	Typeid       string `json:"typeid"`
	Views        string `json:"views"`
}
