package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strings"

	"github.com/fogleman/gg"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/witchert/hitomezashi/hitomezashi"
)

const MAX_SIZE = 257

// Convert a given string to binary representation for each rune
func stringToBinary(str string) string {
	var buffer bytes.Buffer
	for _, runeValue := range str {
		fmt.Fprintf(&buffer, "%b", runeValue)
	}
	return fmt.Sprintf("%s", buffer.Bytes())
}

// Get commit hash from remote github repository
func hashFromLSRemote(owner string, repo string, branch string) (string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{Name: "origin", URLs: []string{"https://github.com/" + owner + "/" + repo}})
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	hash := ""
	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name().String() == "refs/heads/"+branch {
			hash = ref.Hash().String()
		}
	}

	return hash, nil
}

func generateImage(owner string, repo string, branch string) (*gg.Context, error) {
	size := 280
	if size >= MAX_SIZE {
		size = MAX_SIZE
	}

	hash, err := hashFromLSRemote(owner, repo, branch)
	if err != nil {
		return nil, err
	}
	hash = stringToBinary(hash)

	acrossBinary := hash[:size]
	downBinary := hash[len(hash)-(size):]

	across := strings.Split(acrossBinary, "")
	down := strings.Split(downBinary, "")

	h := hitomezashi.New(across, down, 12, hash)

	return h.Canvas, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	// Dead end anything that isnt the root path
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	canvas, err := generateImage(req.URL.Query().Get("owner"), req.URL.Query().Get("repo"), req.URL.Query().Get("branch"))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	err = png.Encode(w, canvas.Image())
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":1111", nil)
}
