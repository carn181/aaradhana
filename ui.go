package main

// TODO
// - Display Lyrics DONE
// - Title in lyrics
// - Time update
// - Change Wallpaper
// - Filter by Song Category
// - Editing lyrics ?

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"strings"
	"time"
)

func closeWindow() {
	gtk.MainQuit()
	fmt.Println("Ending GTK Program")
}

var current = 0
var lyricsDisplay = false
var secondlyric = true

func addResult(listbox *gtk.ListBox, song Song, queryNum int) {
	addResult := func() bool {
		if current == queryNum {
			row, _ := gtk.ListBoxRowNew()

			label, _ := gtk.LabelNew(song.name)
			row.Add(label)

			listbox.Add(row)
			listbox.ShowAll()
		}

		return false
	}

	glib.IdleAdd(addResult)
	time.Sleep(time.Second / 100) // To avoid conflicts and reasonable waiting time
}

func addResults(listbox *gtk.ListBox, songs []Song, query int) {
	for _, song := range songs {
		addResult(listbox, song, query)
	}
}

func clearResults(listbox *gtk.ListBox) {
	results := listbox.GetChildren()
	results.Foreach(func(item interface{}) {
		item.(*gtk.Widget).Destroy()
	})
}

//func displaySong(string text, possibly string text2){
//	win
//	styling
//	label
//	if(lyrics2) label2s
//	showAll
//}

func sanitizeLyrics(text string) []string {
	s := strings.Split(
		strings.Replace(
			strings.Replace(text,
				"<BR>", "\n", -1),
			"<br>", "\n", -1),
		"<slide>")
	return s
}

type LyricsWindow struct {
	win        *gtk.Window
	overlay    *gtk.Overlay
	lyric1     *gtk.Label
	lyric2     *gtk.Label
	box        *gtk.Box
	background *gtk.Image
}

func AddWidgetClass(widget gtk.IWidget, class_name string) {
	context, _ := widget.ToWidget().GetStyleContext()
	context.AddClass(class_name)
}

func (lw *LyricsWindow) init(lyrics1 string, lyrics2 string) {
	css := ".lyrics1 {color: rgb(255,255,255); font-size:40px; font-family: Iosevka, serif; background-color: rgba(0,0,0,0.2);}\n.lyrics2 {font-size:40px; background-color: rgba(0,0,0,0.2);}\nbox {background-color: rgba(0,0,0,0.2);}\n"

	provider, _ := gtk.CssProviderNew()
	provider.LoadFromData(css)

	lw.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	display, _ := gdk.DisplayGetDefault()
	monitor, _ := display.GetMonitor(0)
	rect := monitor.GetWorkarea()
	lw.win.SetDefaultSize(rect.GetWidth(), rect.GetHeight())
	lw.win.Fullscreen()

	lw.overlay, _ = gtk.OverlayNew()

	lw.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	lw.lyric1, _ = gtk.LabelNew(lyrics1)
	lw.lyric1.SetName("lyrics1")
	AddWidgetClass(lw.lyric1, "lyrics1")

	lw.lyric2, _ = gtk.LabelNew(lyrics2)
	lw.lyric2.SetName("lyrics2")
	AddWidgetClass(lw.lyric2, "lyrics2")

	context, _ := lw.lyric1.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	context, _ = lw.lyric2.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	w, h := lw.win.GetSize()

	pixbuf, _ := gdk.PixbufNewFromFile("walls/3cross.jpg")
	pixbufResized, _ := pixbuf.ScaleSimple(w, h, gdk.INTERP_BILINEAR)
	lw.background, _ = gtk.ImageNewFromPixbuf(pixbufResized)

	lw.box.PackStart(lw.lyric1, true, true, 0)

	if lyrics2 == "" {
		lw.box.Remove(lw.lyric2)
		secondlyric = false
	} else {
        lw.box.PackStart(lw.lyric2, true, true, 0)
    }

	lw.overlay.AddOverlay(lw.background)
	lw.overlay.AddOverlay(lw.box)

	lw.win.Add(lw.overlay)

	lw.win.Connect("destroy",
		func() {
			lw.win.Destroy()
			lyricsDisplay = false
		})

	lw.win.ShowAll()
	lyricsDisplay = true
}

func (lw *LyricsWindow) update(lyrics1 string, lyrics2 string) {
	lw.lyric1.SetText(lyrics1)
	if lyrics2 != "" {
		if !secondlyric {
			lw.box.PackStart(lw.lyric2, true, true, 30)
			secondlyric = true
		}
		lw.lyric2.SetText(lyrics2)
        lw.lyric2.QueueDraw()
	} else {
		if secondlyric {
			lw.box.Remove(lw.lyric2)
            lw.lyric1.QueueDraw()
            secondlyric = false
		}
	}
	lw.win.ShowAll()
}

func main() {
	gtk.Init(nil)
	builder, _ := gtk.BuilderNewFromFile("ui.glade")
	var ld LyricsWindow
	obj, _ := builder.GetObject("window1")
	window := obj.(*gtk.Window)

	var signals = map[string]interface{}{
		"window-destroy": closeWindow,
		//		"show-song": displayLyrics,
	}

	songs := db_query()

	obj, _ = builder.GetObject("searchresults")
	lyricsBox := obj.(*gtk.ListBox)

	obj, _ = builder.GetObject("search")
	searchEntry := obj.(*gtk.SearchEntry)

	lyricsBox.Connect("row-activated", func() {

		obj, _ = builder.GetObject("lyrics")
		lyrics := obj.(*gtk.FlowBox)

		children := lyrics.GetChildren()
		children.Foreach(func(item interface{}) {
			item.(*gtk.Widget).Destroy()
		})

		label, _ := lyricsBox.GetSelectedRow().GetChild()
		songTitle, _ := label.(*gtk.Label).GetText()
		song := lookupSong(songs, songTitle)
		stanzas1 := sanitizeLyrics(song.lyrics)
		stanzas2 := sanitizeLyrics(song.lyrics2)

		multipleLyricsExist := len(stanzas2) == len(stanzas1)

		if multipleLyricsExist {
			for i := 0; i < len(stanzas1); i++ {
				stanza1 := stanzas1[i]
				stanza2 := stanzas2[i]
				stanza1lb, _ := gtk.LabelNew(stanza1)
				stanza1lb.SetLineWrap(true)
				stanza2lb, _ := gtk.LabelNew(stanza2)
				stanza2lb.SetLineWrap(true)

				button, _ := gtk.ButtonNew()
				box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
				box.PackStart(stanza1lb, true, true, 3)
				box.PackStart(stanza2lb, true, true, 3)

				button.Add(box)
				button.Connect("clicked", func() {
					if lyricsDisplay {
						ld.update(stanza1, stanza2)
					} else {
						ld.init(stanza1, stanza2)
					}
				})
				lyrics.Insert(button, -1)
			}
		} else {
			for i := 0; i < len(stanzas1); i++ {
				stanza1 := stanzas1[i]
				stanzalb, _ := gtk.LabelNew(stanza1)
				stanzalb.SetLineWrap(true)

				button, _ := gtk.ButtonNew()
				button.Add(stanzalb)

				button.Connect("clicked", func() {
					if lyricsDisplay {
						ld.update(stanza1, "")
					} else {
						ld.init(stanza1, "")
					}
				})
				lyrics.Insert(button, -1)
			}
		}

		lyrics.ShowAll()
		//		lyric.SetText(strings.Split(strings.Replace(song.lyrics, "<BR>", "\n", -1), "<slide>")[0])
	})

	searchEntry.Connect("search-changed", func() {
		clearResults(lyricsBox)
		text, _ := searchEntry.GetText()
		songs = db_query(text)
		current++
		go addResults(lyricsBox, songs, current)
	})

	builder.ConnectSignals(signals)

	window.ShowAll()

	gtk.Main()
}
