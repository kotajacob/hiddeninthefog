package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (app *application) routes() http.Handler {
	static := http.FileServer(http.FS(Static))

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.path)
	mux.Handle("/static/", static)

	return app.logRequest(mux)
}

// path is an http.HandlerFunc which passes the request to either artist,
// video, or file depending on if the request is for a file, video file, or
// directory.
func (app *application) path(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	info, err := os.Stat(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	if !info.IsDir() {
		app.video(w, r)
		return
	}

	app.list(w, r)
}

// ListPage is the datastructure used on list pages.
type ListPage struct {
	Title   string
	Entries []DirEntry
}

type DirEntry struct {
	Name string
	URL  string
	Ext  string
	Size int64
	Time string
}
type Entries []DirEntry

func (e Entries) Len() int      { return len(e) }
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Entries) Less(i, j int) bool {
	return e[i].Name > e[j].Name
}

// list is an http.HandlerFunc which displays a directories listings.
func (app *application) list(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	title := strings.TrimPrefix(path, filepath.Clean(app.dir))
	title = strings.TrimPrefix(title, "/")
	if title == "" {
		title = "hidden in the fog"
	}
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		app.errLog.Println("failed reading directory:", path)
		http.NotFound(w, r)
		return
	}

	var entries Entries
	for _, e := range dirEntries {
		if strings.HasSuffix(e.Name(), ".vtt") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		if len(name) > 60 {
			name = name[:57] + "..."
		}
		if e.IsDir() {
			name += "/"
		}

		info, err := e.Info()
		if err != nil {
			app.errLog.Println("failed reading file info:", path)
			http.NotFound(w, r)
			return
		}

		entries = append(entries, DirEntry{
			Name: name,
			URL:  filepath.Join(strings.TrimPrefix(path, app.dir), e.Name()),
			Size: info.Size(),
			Time: info.ModTime().Format("Jan 02 15:04 2006"),
		})
	}
	sort.Sort(entries)

	tsName := "list.tmpl"
	ts, ok := app.templateCache[tsName]
	if !ok {
		app.errLog.Println(fmt.Errorf(
			"the template %s is missing",
			tsName,
		))
		http.NotFound(w, r)
		return
	}
	err = ts.ExecuteTemplate(w, tsName, ListPage{
		Title:   title,
		Entries: entries,
	})
	if err != nil {
		app.errLog.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}

// VideoPage is the datastructure used in the video handler for the video
// template.
type VideoPage struct {
	Title     string
	Directory string
	URL       string
	VTT       string
	Ext       string
}

// video is an http.HandlerFunc for videos.
// The video is displayed in browser with a helpful player.
func (app *application) video(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	_, direct := q["direct"]
	if direct {
		path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
		info, err := os.Stat(path)
		if err != nil {
			app.errLog.Println(err)
			http.NotFound(w, r)
			return
		}
		f, err := os.Open(path)
		if err != nil {
			app.errLog.Println(err)
			http.NotFound(w, r)
			return
		}
		http.ServeContent(w, r, path, info.ModTime(), f)
		return
	}

	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	title := strings.Trim(filepath.Base(path), filepath.Ext(path))
	if title == "" {
		title = "hidden in the fog"
	}
	tsName := "video.tmpl"
	ts, ok := app.templateCache[tsName]
	if !ok {
		app.errLog.Println(fmt.Errorf(
			"the template %s is missing",
			tsName,
		))
		http.NotFound(w, r)
		return
	}
	err := ts.ExecuteTemplate(w, tsName, VideoPage{
		Title: title,
		Directory: strings.TrimPrefix(
			filepath.Dir(filepath.Clean(r.URL.Path)),
			"/",
		),
		URL: filepath.Clean(r.URL.Path),
		VTT: strings.TrimSuffix(
			filepath.Clean(r.URL.Path),
			filepath.Ext(r.URL.Path),
		) + ".vtt",
		Ext: strings.TrimPrefix(filepath.Ext(r.URL.Path), "."),
	})
	if err != nil {
		app.errLog.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
