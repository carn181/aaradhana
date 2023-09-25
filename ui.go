package main

import (
	"fmt"
	"strings"
	"github.com/gotk3/gotk3/gtk"
)

func closeWindow() {
	gtk.MainQuit()
	fmt.Println("Ending GTK Program")
}

func addLyrics(listbox *gtk.ListBox, songs []Song) {
	for _, song := range songs {
		row, _ := gtk.ListBoxRowNew()

		label, _ := gtk.LabelNew(song.name)
		row.Add(label)

		listbox.Add(row)
	}
}

func displayLyrics(listbox *gtk.ListBox, row *gtk.ListBoxRow) {
	label, _ := listbox.GetSelectedRow().GetChild()
	s, _ := label.(*gtk.Label).GetText()
	fmt.Println(s)
}


func main() {
	gtk.Init(nil)
	builder, _ := gtk.BuilderNewFromFile("ui.glade")
	obj, _ := builder.GetObject("window1")
	window := obj.(*gtk.Window)

	var signals = map[string]interface{}{
		"window-destroy": closeWindow,
		//		"show-song": displayLyrics,
	}

	songs := db_query()

	
	obj, _ = builder.GetObject("searchresults")
	lyricsBox := obj.(*gtk.ListBox)
	
	lyricsBox.Connect("row-activated", func(){
			
		obj, _ = builder.GetObject("lyrics")
		lyrics := obj.(*gtk.FlowBox)
		
		children := lyrics.GetChildren()
		children.Foreach(func(item interface{}){
			item.(*gtk.Widget).Destroy()
		})
		
		label, _ := lyricsBox.GetSelectedRow().GetChild()
		songTitle, _ := label.(*gtk.Label).GetText()
		song := lookupSong(songs, songTitle)
		stanzas1 := strings.Split(strings.Replace(song.lyrics, "<BR>", "\n", -1), "<slide>")
		stanzas2 := strings.Split(
			strings.Replace(
				strings.Replace(song.lyrics2,
					"<BR>", "\n", -1),
				"<br>", "\n", -1),
			"<slide>")

		multipleLyricsExist := len(stanzas2)==len(stanzas1)

		if multipleLyricsExist {
			for i := 0 ; i < len(stanzas1) ; i++ {
				stanza1 := stanzas1[i]
				stanza2 := stanzas2[i]
				stanza1lb, _ := gtk.LabelNew(stanza1)
				stanza1lb.SetLineWrap(true)
				stanza2lb, _ := gtk.LabelNew(stanza2)
				stanza2lb.SetLineWrap(true)

				box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
				box.PackStart(stanza1lb, true, true, 3)
				box.PackStart(stanza2lb, true, true, 3)				
				
				lyrics.Insert(box, -1)
			}
		} else {
			for i := 0 ; i < len(stanzas1) ; i++ {
				stanza1 := stanzas1[i]
				stanzalb, _ := gtk.LabelNew(stanza1)
				stanzalb.SetLineWrap(true)
				
				lyrics.Insert(stanzalb, -1)
			}
		}
		

		lyrics.ShowAll()
		//		lyric.SetText(strings.Split(strings.Replace(song.lyrics, "<BR>", "\n", -1), "<slide>")[0])
	})

	addLyrics(lyricsBox, songs)

	builder.ConnectSignals(signals)

	window.ShowAll()
	gtk.Main()
}
