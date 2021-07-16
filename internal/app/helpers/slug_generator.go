package helpers

import "strings"

var toLat = strings.NewReplacer(
	" ", "_",
	"а", "a",
	"б", "b",
	"в", "v",
	"г", "g",
	"д", "d",
	"е", "e",
	"ё", "yo",
	"ж", "zh",
	"з", "z",
	"и", "i",
	"й", "y",
	"к", "k",
	"л", "l",
	"м", "m",
	"н", "n",
	"о", "o",
	"п", "p",
	"р", "r",
	"с", "s",
	"т", "t",
	"у", "u",
	"ф", "f",
	"х", "kh",
	"ц", "ts",
	"ч", "ch",
	"ш", "sh",
	"щ", "sch",
	"ъе", "ye",
	"ъ", "",
	"ый", "iy",
	"ий", "iy",
	"ы", "y",
	"ь", "",
	"э", "e",
	"ю", "yu",
	"я", "ya",
)

func GenerateSlug(s string) string {
	return toLat.Replace(strings.ToLower(s))
}
