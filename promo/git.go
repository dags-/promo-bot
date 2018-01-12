package promo

import (
	"fmt"
	"github.com/dags-/promo-bot/util"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"os"
	"strings"
)

func GetPromotions(owner, repo string) (map[string]map[string]Promotion, error) {
	var promos map[string]map[string]Promotion

	url := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{URL: url})
	if err != nil {
		return promos, err
	}

	w, err := r.Worktree()
	if err != nil {
		return promos, err
	}

	files, err := w.Filesystem.ReadDir("/")
	if err != nil {
		return promos, err
	}

	promos = readPromotions(files, w.Filesystem)

	return promos, nil
}

func readPromotions(files []os.FileInfo, fs billy.Filesystem) (map[string]map[string]Promotion) {
	promos := map[string]map[string]Promotion{
		"server":  make(map[string]Promotion),
		"twitch":  make(map[string]Promotion),
		"youtube": make(map[string]Promotion),
	}

	for _, fi := range files {
		if !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}

		f, err := fs.Open(fi.Name())
		if err != nil {
			fmt.Println("git.readpromos.open: ", err)
			continue
		}

		var pr Promotion
		err = utils.DecodeJson(&pr, f)
		if err != nil {
			fmt.Println("git.readpromos.decode: ", err)
			continue
		}

		err = f.Close()
		if err != nil {
			fmt.Println("git.readpromos.close: ", err)
		}

		if section, ok := promos[pr.Type]; ok {
			section[pr.ID] = pr
		} else {
			fmt.Println("git.readpromos.add: invalid promo type '", pr.Type, "'")
		}
	}

	return promos
}
