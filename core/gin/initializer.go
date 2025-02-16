package gin

import (
	"io/fs"
	"net/http"
	"strings"

	"chatgpt-adapter/core/logger"
	"chatgpt-adapter/www"

	"github.com/gin-gonic/gin"
	"github.com/iocgo/sdk"
	"github.com/iocgo/sdk/env"
	"github.com/iocgo/sdk/router"
)

var (
	debug bool
)

// @Inject(lazy="false", name="ginInitializer")
func Initialized(env *env.Environment) sdk.Initializer {
	debug = env.GetBool("server.debug")

	public, _ := fs.Sub(www.Dist, "public")
	_next, _ := fs.Sub(public, "_next")

	return sdk.InitializedWrapper(0, func(container *sdk.Container) (err error) {
		sdk.ProvideTransient(container, sdk.NameOf[*gin.Engine](), func() (engine *gin.Engine, err error) {
			if !debug {
				gin.SetMode(gin.ReleaseMode)
			}

			engine = gin.New()
			{
				engine.Use(gin.Recovery())
				engine.Use(cros)
				engine.Use(token)
			}
			engine.Static("/file/", "tmp")
			engine.StaticFS("/_next", http.FS(_next))
			engine.NoRoute(func(c *gin.Context) {
				filename := strings.TrimLeft(c.Request.URL.Path, "/")
				if filename == "/" || filename == "" {
					filename = "index.html"
				}
				if _, err = fs.Stat(public, filename); err != nil {
					logger.Errorf("文件不存在: %s %v", filename, err)
					c.AbortWithStatus(http.StatusNotFound)
					return
				}
				http.FileServer(http.FS(public)).ServeHTTP(c.Writer, c.Request)
			})
			beans := sdk.ListInvokeAs[router.Router](container)
			for _, route := range beans {
				route.Routers(engine)
			}

			return
		})
		return
	})
}

func token(gtx *gin.Context) {
	str := gtx.Request.Header.Get("X-Api-Key")
	if str == "" {
		str = strings.TrimPrefix(gtx.Request.Header.Get("Authorization"), "Bearer ")
	}
	gtx.Set("token", str)
}

func cros(gtx *gin.Context) {
	method := gtx.Request.Method
	gtx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	gtx.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
	gtx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
	gtx.Header("Access-Control-Allow-Headers", "*")
	gtx.Header("Access-Control-Expose-Headers", "*")
	gtx.Header("Access-Control-Max-Age", "172800")
	gtx.Header("Access-Control-Allow-Credentials", "false")
	//gtx.Set("content-type", "application/json")

	if method == "OPTIONS" {
		gtx.Status(http.StatusOK)
		return
	}

	if gtx.Request.RequestURI == "/" ||
		gtx.Request.RequestURI == "/favicon.ico" ||
		strings.Contains(gtx.Request.URL.Path, "/v1/models") ||
		strings.HasPrefix(gtx.Request.URL.Path, "/file/") {
		// 处理请求
		gtx.Next()
		return
	}
	// 处理请求
	gtx.Next()
}
