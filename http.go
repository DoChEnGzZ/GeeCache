package GeeCahce

import (
	"fmt"
	"github.com/DoChEnGzZ/GeeDo"
	"log"
	"net/http"
)

const defaultBasePath = "/_geeCache"

func NewHTTPPool(self string) *HttpPool {
	p := &HttpPool{
		engine:   Gee.Default(),
		self:     self,
		basePath: defaultBasePath,
	}
	BaseGroup := p.engine.Group(defaultBasePath)
	QueryGroup := BaseGroup.Group("/Query")
	QueryGroup.Get("/", func(c *Gee.Context) {
		//if !strings.HasPrefix(c.Req.URL.Path, p.basePath){
		//	c.String(http.StatusBadRequest,"bad request")
		//	p.Log("incorrect request url:%s",c.Req.URL.Path)
		//	return
		//}
		p.Log("%s %s", c.Req.Method, c.Req.URL.Path)
		//r.URL.Path:/<basepath>/<groupname>/<key>
		//parts:/<groupname>/<key>
		//parts := strings.SplitN(c.Req.URL.Path[len(p.basePath):], "/", 3)
		//if len(parts) != 3 {
		//	c.String(http.StatusBadRequest,"Bad request")
		//	p.Log("incorrect request url:%s",c.Req.URL.Path)
		//	return
		//}
		groupName := c.PostForm("GroupName")
		key := c.PostForm("Key")
		group := GetGroup(groupName)
		if group == nil {
			c.String(http.StatusNotFound, "no such group:"+groupName)
			p.Log("no such group:%s", groupName)
			return
		}
		v, err := group.Get(key)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			p.Log("geeCache get error:%s", err.Error())
			return
		}
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		_, err = c.Writer.Write(v.ByteSlice())
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			p.Log("write error", err.Error())
		}
	})
	//p.engine.Post()
	return p
}

//
// HttpPool
// @Description: implements PeerPicker for a pool of HTTP peers. handler in fact
//
type HttpPool struct {
	engine   *Gee.Engine
	self     string
	basePath string
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s],%s", p.self, fmt.Sprintf(format, v))
}

func (p *HttpPool) Run() (err error) {
	err = p.engine.Run(p.self)
	return err
}
