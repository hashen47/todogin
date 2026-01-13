package todo

type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
	UserId  int    `json:"user_id"`
}

type TodoCreateReq struct {
	Title   string `json:"title" binding:"required,min=5,max=100"`
	Content string `json:"content" binding:"required,min=5,max=255"`
}

type TodoGetReq struct {
	Offset int `json:"offset" binding:"gte=0"`
	Limit  int `json:"limit" binding:"required,gte=1,lte=100"`
}

type TodoUpdateReq struct {
	Id      int    `json:"id" binding:"required,gte=1"`
	Title   string `json:"title" binding:"required,min=5,max=100"`
	Content string `json:"content" binding:"required,min=5,max=255"`
	Done    bool   `json:"done" binding:"boolean"`
}

type TodoDeleteReq struct {
	Id int `json:"id" binding:"required,gte=1"`
}
