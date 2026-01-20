package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"time"

	web "github.com/ReilEgor/SiteSentinel/frontEnd/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.New()
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))
	tmpl := template.Must(template.New("").ParseFS(web.TemplateFS, "templates/*.html", "templates/partials/*.html"))
	server.SetHTMLTemplate(tmpl)
	staticFS, _ := fs.Sub(web.StaticFS, "static")
	server.StaticFS("/static", http.FS(staticFS))

	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	server.Run(":8080")
}
