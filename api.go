package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"regexp"

	"github.com/dwikng/gofali/storage"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

//go:embed admin/*
var adminFiles embed.FS

type api struct {
	server *fasthttp.Server

	store storage.Storage
}

func (a *api) run() error {
	r := router.New()

	r.GET("/", a.handleRoot)
	r.GET("/{slug}", a.handleSlug)

	adminGroup := r.Group("/@" + conf.adminPath)
	adminGroup.POST("/create", a.handleAdminCreate)
	adminGroup.GET("/all", a.handleAdminGetAll)
	adminGroup.PUT("/edit", a.handleAdminEdit)
	adminGroup.DELETE("/delete", a.handleAdminDelete)

	files, _ := fs.Sub(adminFiles, "admin")
	adminGroup.GET("/{file:*}", fasthttpadaptor.NewFastHTTPHandler(
		http.StripPrefix("/@"+conf.adminPath, http.FileServer(http.FS(files))),
	))

	a.server = &fasthttp.Server{
		Handler:           r.Handler,
		ReduceMemoryUsage: true,
	}

	addr := fmt.Sprintf("%s:%d", conf.host, conf.port)
	return a.server.ListenAndServe(addr)
}

func (a *api) handleRoot(ctx *fasthttp.RequestCtx) {
	if conf.rootRedirect != "" {
		ctx.Redirect(conf.rootRedirect, fasthttp.StatusFound)
		return
	}

	ctx.NotFound()
}

func (a *api) handleSlug(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	link := a.store.TryUse(slug)
	if link != nil {
		ctx.Redirect(link.Url, fasthttp.StatusFound)
		return
	}

	ctx.NotFound()
}

func (a *api) handleAdminRoot(ctx *fasthttp.RequestCtx) {
	path := "/" + ctx.UserValue("file").(string)
	if path == "//" {
		path = "/index.html"
	}

	bytes, err := adminFiles.ReadFile(path)
	if err != nil {
		ctx.NotFound()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/html")
	_, _ = ctx.Write(bytes)
}

func (a *api) handleAdminCreate(ctx *fasthttp.RequestCtx) {
	url := ctx.PostArgs().Peek("url")

	if len(url) == 0 || !isValidURL(string(url)) {
		respondJSON(ctx, fasthttp.StatusBadRequest, errorResponse{Message: "Invalid url"})
		return
	}

	link, err := a.store.Create(string(url))
	if err != nil {
		respondJSON(ctx, fasthttp.StatusInternalServerError, errorResponse{Message: err.Error()})
		return
	}

	respondJSON(ctx, fasthttp.StatusCreated, link)
}

func (a *api) handleAdminGetAll(ctx *fasthttp.RequestCtx) {
	links := a.store.All()

	respondJSON(ctx, fasthttp.StatusOK, links)
}

func (a *api) handleAdminEdit(ctx *fasthttp.RequestCtx) {
	slug := ctx.PostArgs().Peek("slug")
	url := ctx.PostArgs().Peek("url")

	if len(url) == 0 || !isValidURL(string(url)) {
		respondJSON(ctx, fasthttp.StatusBadRequest, errorResponse{Message: "Invalid url"})
		return
	}

	err := a.store.Update(string(slug), string(url))
	if err != nil {
		respondJSON(ctx, fasthttp.StatusInternalServerError, errorResponse{Message: err.Error()})
	}

	respondJSON(ctx, fasthttp.StatusOK, map[string]string{
		"message": fmt.Sprintf("Slug '%s' updated successfully", string(slug)),
	})
}

func (a *api) handleAdminDelete(ctx *fasthttp.RequestCtx) {
	slug := ctx.PostArgs().Peek("slug")

	err := a.store.Delete(string(slug))
	if err != nil {
		respondJSON(ctx, fasthttp.StatusInternalServerError, errorResponse{Message: err.Error()})
		return
	}

	respondJSON(ctx, fasthttp.StatusOK, map[string]string{
		"message": fmt.Sprintf("Slug '%s' deleted successfully", string(slug)),
	})
}

func (a *api) shutdown() error {
	return a.server.ShutdownWithContext(context.Background())
}

func isValidURL(url string) bool {
	const urlRegex = `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`
	re := regexp.MustCompile(urlRegex)
	return re.MatchString(url)
}

type errorResponse struct {
	Message string `json:"message"`
}

func respondJSON(ctx *fasthttp.RequestCtx, code int, data interface{}) {
	marshaled, _ := json.Marshal(data)
	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json")
	ctx.SetBody(marshaled)
}
