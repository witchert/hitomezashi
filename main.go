package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/muesli/gamut"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/witchert/hitomezashi/hitomezashi"
)

const MAX_SIZE = 257

var colors = []color.Color{
	color.NRGBA64{57818, 35930, 12168, 65535},
	color.NRGBA64{19216, 48673, 16633, 65535},
	color.NRGBA64{57090, 34948, 47660, 65535},
	color.NRGBA64{36003, 45435, 13660, 65535},
	color.NRGBA64{20404, 39496, 20721, 65535},
	color.NRGBA64{52060, 15473, 58198, 65535},
	color.NRGBA64{46641, 36500, 57269, 65535},
	color.NRGBA64{57480, 15399, 30954, 65535},
	color.NRGBA64{42716, 42029, 24161, 65535},
	color.NRGBA64{39231, 22794, 37294, 65535},
	color.NRGBA64{55575, 38184, 27239, 65535},
	color.NRGBA64{16859, 42870, 53708, 65535},
	color.NRGBA64{41379, 27522, 9102, 65535},
	color.NRGBA64{49700, 24602, 29228, 65535},
	color.NRGBA64{55766, 29478, 56914, 65535},
	color.NRGBA64{42260, 24550, 16624, 65535},
	color.NRGBA64{28496, 23990, 58742, 65535},
	color.NRGBA64{47353, 17967, 10505, 65535},
	color.NRGBA64{58946, 25652, 10518, 65535},
	color.NRGBA64{20285, 45089, 37379, 65535},
	color.NRGBA64{59253, 28945, 25262, 65535},
	color.NRGBA64{55792, 14417, 47691, 65535},
	color.NRGBA64{57658, 13672, 10787, 65535},
	color.NRGBA64{27779, 29712, 10939, 65535},
	color.NRGBA64{27078, 27904, 54256, 65535},
	color.NRGBA64{33328, 15023, 59635, 65535},
	color.NRGBA64{40308, 19184, 46614, 65535},
	color.NRGBA64{27566, 39157, 59264, 65535},
	color.NRGBA64{55481, 14389, 20760, 65535},
	color.NRGBA64{50442, 41923, 13432, 65535},
	color.NRGBA64{23731, 27990, 45070, 65535},
	color.NRGBA64{52365, 18426, 38505, 65535},
}

// Convert a given string to binary representation for each rune
func stringToBinary(str string) []bool {
	var buffer bytes.Buffer
	for _, runeValue := range str {
		fmt.Fprintf(&buffer, "%b", runeValue)
	}
	binaryString := fmt.Sprintf("%s", buffer.Bytes())
	var out []bool
	for _, x := range binaryString {
		b, _ := strconv.ParseBool(string(x))
		out = append(out, b)
	}
	return out
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

func generateImage(owner string, repo string, branch string, size int) (*gg.Context, string, error) {
	if size >= MAX_SIZE {
		size = MAX_SIZE
	}

	hash, err := hashFromLSRemote(owner, repo, branch)
	if err != nil {
		return nil, "", err
	}
	hashBinary := stringToBinary(hash)
	acrossBinary := hashBinary[:size]
	downBinary := hashBinary[len(hashBinary)-(size):]

	h := hitomezashi.New(acrossBinary, downBinary, 12)

	seed, _ := strconv.ParseInt(hash[:7], 16, 64)
	r := rand.New(rand.NewSource(seed))
	colorA := colors[r.Intn(len(colors)-1)+1]
	colorB := gamut.Complementary(colorA)
	h.SetColors(colorA, colorB)

	h.Draw()

	return h.Canvas, hash, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	// Dead end anything that isnt the root path
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	var err error
	size := 24

	if req.URL.Query().Get("size") != "" {
		size, err = strconv.Atoi(req.URL.Query().Get("size"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid size parameter!"))
			return
		}
	}

	if req.URL.Query().Get("owner") == "" || req.URL.Query().Get("repo") == "" || req.URL.Query().Get("branch") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Missing required parameters!"))
		return
	}

	canvas, hash, err := generateImage(req.URL.Query().Get("owner"), req.URL.Query().Get("repo"), req.URL.Query().Get("branch"), size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", hash)

	err = png.Encode(w, canvas.Image())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - " + err.Error()))
	}
}

func main() {
	newrelicAppName := os.Getenv("NEWRELIC_APP_NAME")
	newrelicLicense := os.Getenv("NEWRELIC_LICENSE")
	if newrelicAppName != "" && newrelicLicense != "" {
		app, _ := newrelic.NewApplication(
			newrelic.ConfigAppName(newrelicAppName),
			newrelic.ConfigLicense(newrelicLicense),
			newrelic.ConfigDistributedTracerEnabled(true),
		)

		http.HandleFunc(newrelic.WrapHandleFunc(app, "/", handler))
	} else {
		http.HandleFunc("/", handler)
	}
	http.ListenAndServe(":1111", nil)
}
