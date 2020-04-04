package server

import (
	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"html/template"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
	"yj-app/app/service/system/config"
	"yj-app/app/service/system/dict"
	"yj-app/app/service/system/permission"
	"yj-app/app/yjgframe/router"
)

//已注册HTTP服务列表
var servers = make(map[string]*Server, 0)

//HTTP服务结构体
type Server struct {
	ServerName     string        //服务名称
	Addr           string        //监听地址端口
	ServerRoot     string        //静态资源文件夹
	Handler        *gin.Engine   //HTTP Handler
	ReadTimeout    time.Duration //读取超时时间
	WriteTimeout   time.Duration //写入超时时间
	MaxHeaderBytes int           //http头大小设置
}

//根据name获取Server
func GetServer(name string) *Server {
	return servers[name]
}

// 创建服务
func New(name, addr string, middleware ...gin.HandlerFunc) *Server {
	var s Server
	s.WriteTimeout = 60 * time.Second
	s.ReadTimeout = 60 * time.Second
	s.Addr = addr
	s.ServerName = name
	s.MaxHeaderBytes = 1 << 20
	s.Handler = gin.New()
	s.Handler.Use(middleware...)
	//注册路由
	if len(router.GroupList) > 0 {
		for _, group := range router.GroupList {
			if group.ServerName == name {
				grp := s.Handler.Group(group.RelativePath, group.Handlers...)
				for _, r := range group.Router {
					if r.Method == "ANY" {
						grp.Any(r.RelativePath, r.HandlerFunc...)
					} else {
						grp.Handle(r.Method, r.RelativePath, r.HandlerFunc...)
					}

				}
			}
		}
	}

	servers[name] = &s
	return &s
}

//加载模板
func (s *Server) Template(path string) *Server {
	//s.Handler.SetFuncMap(template.FuncMap{
	//	"hasPermi":          permission.HasPermi,
	//	"getPermiButton":    permission.GetPermiButton,
	//	"getDictLabel":      dict.GetDictLabel,
	//	"getDictTypeSelect": dict.GetDictTypeSelect,
	//	"getDictTypeRadio":  dict.GetDictTypeRadio,
	//	"getDictTypeData":   dict.GetDictTypeData,
	//	"OssUrl":            config.GetOssUrl,
	//})

	s.Handler.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "template",
		Extension: ".html",
		Master:    "",
		Partials:  []string{"header", "footer", "system/menu/icon"},
		Funcs: template.FuncMap{
			"hasPermi":          permission.HasPermi,
			"getPermiButton":    permission.GetPermiButton,
			"getDictLabel":      dict.GetDictLabel,
			"getDictTypeSelect": dict.GetDictTypeSelect,
			"getDictTypeRadio":  dict.GetDictTypeRadio,
			"getDictTypeData":   dict.GetDictTypeData,
			"OssUrl":            config.GetOssUrl,
		},
		DisableCache: true,
	})

	return s
}

//加载静态资源文件
func (s *Server) Static(path string) *Server {
	tmp, _ := os.Getwd()
	tmp += "/" + path
	s.Handler.StaticFS("/static", http.Dir(tmp))
	s.Handler.StaticFile("/favicon.ico", "./favicon.ico")
	return s
}

//启动服务
func (s *Server) Start(g errgroup.Group) {

	server := &http.Server{
		Addr:           s.Addr,
		Handler:        s.Handler,
		ReadTimeout:    s.ReadTimeout,
		WriteTimeout:   s.WriteTimeout,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}

	log.Printf("[%v]Server listen: %v Actual pid is %d", s.ServerName, s.Addr, syscall.Getpid())

	g.Go(func() error {
		return server.ListenAndServe()
	})
}
