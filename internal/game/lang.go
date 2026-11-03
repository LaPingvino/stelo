package game

// Translations (Joop, 2026-11-02: "Translating is ok!"). The original's text
// is Esperanto and stays the default; Lang "en" draws the English below
// instead. Keys are the original strings exactly as they are in levels/ and
// the code (x-system, '#' line breaks); a missing key draws the original.
// The title picture's own lettering is an image and is not translated.

// Lang is the language the game's text is drawn in: "eo" (the original) or "en".
var Lang = "eo"

// Tr gives s in Lang.
func Tr(s string) string {
	if Lang == "en" {
		if t, ok := english[s]; ok {
			return t
		}
	}
	return s
}

var english = map[string]string{
	// HUD and highscores
	"Poentoj: ": "Score: ",
	"Vi tiom bone ludis ke vi aperos en#la listo de la plej bonaj ludantoj!##Via nomo: ":   "You played so well that you will appear#in the list of the best players!##Your name: ",
	"Ho ve... Vi ne suficxe bone ludis por#aperi en la listo de la plej bonaj#ludantoj...": "Alas... You did not play well enough to#appear in the list of the best#players...",
	// klarigo (the rules screen)
	"Klarigo de la reguloj de la ludo":                                                                                   "Explanation of the rules of the game",
	"La ludanto#Movu per la sagoklavoj#Salto = sago supren":                                                              "The player#Move with the arrow keys#Jump = arrow up",
	"Muroj#Oni ne povas nek rajtas#trairi ilin...":                                                                       "Walls#You can't, and may not,#go through them...",
	"Malamiko#Saltu sur la supron por mortigi#tiun kanajlon kaj gajni poentojn":                                          "Enemy#Jump on its top to kill#that rascal and earn points",
	"Malamiko#Nemortigebla malamiko":                                                                                     "Enemy#An enemy that can't be killed",
	"Verda stelo#Prenu por gajni poentojn":                                                                               "Green star#Take it to earn points",
	"Fino de la nivelo":                                                                                                  "End of the level",
	"Sxtupetaro#Por iri supren kaj malsupren":                                                                            "Ladder#To go up and down",
	"Fajro#Atentu ke vi ne falu en la fajron":                                                                            "Fire#Take care not to fall into the fire",
	"Koro, donas pli da vivo#La steletoj super la kapo de vi#indikas kiom da vivo vi havas":                              "Heart, gives more life#The little stars above your head#show how much life you have",
	"Transirejo#Se vi preteriris gxin kaj vi mortas,#vi cxi tie rekomencos":                                              "Checkpoint#If you have passed it and you die,#you will start again here",
	"Bonusnivelo#Gajnu pli da poentoj per#mortigado de malamikoj#Estas po unu bonusnivelo#en cxiuj niveloj":              "Bonus level#Earn more points by#killing enemies#There is one bonus level#in every level",
	"La ludo estas dedicxita al#Tsukuru, mia fora frato el Japanio#           kaj#Aleks, mia plej bona amiko el Francio": "The game is dedicated to#Tsukuru, my faraway brother from Japan#           and#Aleks, my best friend from France",
	// nivelo5
	"NUN IRU AL LA FINSIGNO POR FINLUDI LA TUTAN LUDON": "NOW GO TO THE FINISH SIGN TO FINISH THE WHOLE GAME",
	// danklisto (the credits)
	"Cxi tiu ludo ne efektivigxus cxi tiel sen la sekvaj:":                                   "This game would not have come about like this without the following:",
	"DIO - Sen kiu nek programistoj nek ludantoj ekzistus":                                   "GOD - Without whom neither programmers nor players would exist",
	"Mark Overmars - Kiu kreis la programon Gamemaker kaj tiel#ebligis facilan ludfaradon":   "Mark Overmars - Who created the program Game Maker and so#made making games easy",
	"LaPingvino - La programisto kiu faris la ludon":                                         "LaPingvino - The programmer who made the game",
	"Kolonjano - Kiu helpis min per testado kaj utilaj konsiloj":                             "Kolonjano - Who helped me with testing and useful advice",
	"L. L. Zamenhof - Pro lia miriga unuiga lingvo":                                          "L. L. Zamenhof - For his amazing unifying language",
	"Miaj gefratoj - Pro la plezuro per kiu ili testadis la ludon":                           "My brothers and sisters - For the pleasure with which they tested the game",
	"Miaj gepatroj - Pro la tekokomputilo per kiu mi la ludon programis":                     "My parents - For the laptop on which I programmed the game",
	"Tsukuru - Kiu estante en Italio kuragxigis min plufari la ludon":                        "Tsukuru - Who, while in Italy, encouraged me to go on with the game",
	"La Interreto, precipe vgmusic.com - Pro la muziko":                                      "The Internet, especially vgmusic.com - For the music",
	"La Esperantistaro - Mia nekompreneble bona celgrupo":                                    "The Esperanto community - My incomprehensibly good audience",
	"Retposxtu al mi ke vi venis cxi tie, por ke mi faru pliajn nivelojn#[ikojba@gmail.com]": "E-mail me that you got this far, so that I make more levels#[ikojba@gmail.com]",
}

// MenuWords are the title menu's words (cmd/stelo/title.go), translated in
// english like the rest.
var MenuWords = []string{"Ludi", "Daŭrigi", "Niveloj", "Lingvo", "Arto", "nova", "de 2006", "Redaktilo", "Fini",
	"Esperanto", "Platforma ludo · 2006, refarita 2026", "nivelo",
	"Paŭzo", "Plej bonaj", "Via nomo:", "Nivelo finita!", "Nova nivelo", "Malfermi", "El dosiero…", "Konservi", "Daŭrigi redakti", "Nivelo", "Bonusnivelo", "Poentoj", "Rekomenci la nivelon", "Sono", "ŝaltita", "malŝaltita", "Ĉefmenuo"}

func init() {
	for k, v := range map[string]string{"Ludi": "Play", "Daŭrigi": "Continue", "Niveloj": "Levels", "Lingvo": "Language",
		"Arto": "Art", "nova": "new", "de 2006": "2006", "Redaktilo": "Level editor", "Fini": "Quit",
		"Esperanto": "Esperanto", "Platforma ludo · 2006, refarita 2026": "A platform game · 2006, remade 2026", "nivelo": "level",
		"Paŭzo": "Paused", "Plej bonaj": "Best players", "Via nomo:": "Your name:", "Nivelo finita!": "Level complete!", "Nova nivelo": "New level", "Malfermi": "Open", "El dosiero…": "From a file…",
		"Konservi": "Save", "Daŭrigi redakti": "Back to editing", "Nivelo": "Level", "Bonusnivelo": "Bonus level", "Poentoj": "Score", "Rekomenci la nivelon": "Restart the level", "Sono": "Sound", "ŝaltita": "on", "malŝaltita": "off",
		"Ĉefmenuo": "Main menu"} {
		english[k] = v
	}
}
