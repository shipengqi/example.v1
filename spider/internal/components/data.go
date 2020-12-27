package component

import "net/http"

type Request struct {
	req *http.Request
	// request depth, must be a positive int
	depth uint32
}

func NewRequest(req *http.Request, depth uint32) *Request {
	return &Request{req: req, depth: depth}
}

func (r *Request) GetReq() *http.Request {
	return r.req
}

func (r *Request) GetDepth() uint32 {
	return r.depth
}

func (r *Request) Valid() bool {
	return r.req != nil && r.req.URL != nil
}

type Response struct {
	res   *http.Response
	depth uint32
}

func NewResponse(res *http.Response, depth uint32) *Response {
	return &Response{res: res, depth: depth}
}

func (r *Response) GetRes() *http.Response {
	return r.res
}

func (r *Response) GetDepth() uint32 {
	return r.depth
}

func (r *Response) Valid() bool {
	return r.res != nil && r.res.Body != nil
}

type Item map[string]interface{}

func (i *Item) Valid() bool {
	return i != nil
}

type Data interface {
	Valid() bool
}
