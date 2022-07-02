package encrypter

import (
	"bint.com/internal/const/options"
	. "bint.com/pkg/serviceTools"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

func Encrypt(rootSource string, rootDest string, keyDest string) {
	var encryptDest *os.File
	var kDest *os.File
	var key string
	var filesListToExecute []string
	var fileToExecute string
	var splitedChunk []string
	var encryptedChunk []string
	var err error
	var res string
	var trash []string
	var generateTrash bool
	var splitedChunkNumber int
	var minVal int

	trash = []string{"AND", "OR", "NOT", "XOR", ".", "+", "-", "*", "/", "^", "print", "len",
		"index", "is_letter", "is_digit", "pop", "push", "input", "next_command", "get_root_source",
		"get_root_dest", "send_command", "goto", "SET_SOURCE", "SET_DEST", "SEND_DEST", "DEL_DEST",
		"UNDEFINE", "UNSET_SOURCE", "REROUTE", "UNSET_DEST", "RESET_SOURCE", "str", "int", "float",
		"bool", "(", ")", "[", "]", ":", "True", "False", "=", "<=", ">=", "==", "<", ">", ",", "string",
		"stack"}
	rand.Seed(time.Now().Unix())

	encryptDest, err = os.Create(rootDest)
	if nil != err {
		panic(err)
	}
	kDest, err = os.Create(keyDest)
	if nil != err {
		panic(err)
	}

	filesListToExecute = []string{rootSource}

	for _, fileToExecute = range filesListToExecute {
		f, err := os.Open(fileToExecute)

		if nil != err {
			panic(err)
		}

		newChunk := EachChunk(f)
		var trashNumbersList []int
		var trashNumber int

		for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
			splitedChunk = strings.Split(chunk, options.BendSep)
			trashLen := rand.Intn(options.MaxTrashLen)
			encryptedChunkLen := trashLen + len(splitedChunk)

		Outer:
			for i := 0; i < trashLen; i++ {
				if "#" == string(strings.TrimSpace(splitedChunk[0])[0]) {
					minVal = 1
				} else {
					minVal = 0
				}
				if encryptedChunkLen-1-minVal > 0 {
					trashNumber = rand.Intn(encryptedChunkLen-minVal) + minVal
				} else {
					encryptedChunk = splitedChunk
					encryptedChunkLen = len(encryptedChunk)
					break
				}
				for j := 0; j < len(trashNumbersList); j++ {
					if trashNumbersList[j] == trashNumber {
						i--
						continue Outer
					}
				}
				trashNumbersList = append(trashNumbersList, trashNumber)
			}

			sort.Ints(trashNumbersList)

			encryptedChunk = nil
			for i := 0; i < encryptedChunkLen; i++ {
				encryptedChunk = append(encryptedChunk, "")
			}

			for i := 0; i < encryptedChunkLen; i++ {
				generateTrash = false
				for _, j := range trashNumbersList {
					if i == j {
						generateTrash = true
						break
					}
				}
				if generateTrash {
					t := trash[rand.Intn(len(trash))]
					encryptedChunk[i] = t
				} else {
					encryptedChunk[i] = splitedChunk[splitedChunkNumber]
					splitedChunkNumber++
				}
			}
			splitedChunkNumber = 0
			for i := 0; i < len(encryptedChunk); i++ {
				res += encryptedChunk[i] + options.BendSep
			}
			_, err = encryptDest.WriteString(res + ";")
			if nil != err {
				panic(err)
			}
			var isTrash bool

			for i := 0; i < len(encryptedChunk); i++ {
				for _, tNum := range trashNumbersList {
					if i == tNum {
						isTrash = true
						break
					}
				}
				if !isTrash {
					key += fmt.Sprintf("%v", i) + ", "
				}
				isTrash = false
			}

			key = key[:len(key)-2]
			key += ";\n"

			_, err = kDest.WriteString(key)
			if nil != err {
				panic(err)
			}

			res = ""
			key = ""
			trashNumbersList = nil
		}
	}
}
