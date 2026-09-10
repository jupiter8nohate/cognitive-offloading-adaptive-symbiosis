package coas

import "fmt"

const (
	GLITCHOLOGYSourceRepository = "jupiter8nohate/computational-metacognitive-bilingualism"
	GLITCHOLOGYSourcePath       = "books/ERR_404_GLITCHOLOGY.md"
	GLITCHOLOGYSourceSHA        = "fcbb823faf824bb751c96855fd1100e4e60efea4"
	GLITCHOLOGYPublicName       = "Err⃝or⃟⃤GLITCHOLOGY"
	GLITCHOLOGYDisplayTitle     = "E⃟ r⃟ r⃟⃝ o⃟ r⃟⃤ G⃟ L⃟ I⃟ T⃟ C⃟ H⃟ O⃟ L⃟ O⃟ G⃟ Y⃟"
	GLITCHOLOGYGrammar          = "<GLYPH> [RUNTIME] CLAIM :: STATE :: AUTHORITY"
)

var GLITCHOLOGYCoreLaws = []string{
	"PATTERN != PROOF",
	"PROFILE != PERSON",
	"MODEL != MIND",
	"PREDICTION != DESTINY",
	"CAPABILITY != AUTHORITY",
	"ACCESS != CONSENT",
	"ERROR != GUILT",
	"UNKNOWN != FALSE",
	"DIFFERENCE != DEFECT",
	"SIGNAL != VERDICT",
	"DATA != CONTEXT",
	"INTERPRETATION != EVENT",
	"REPRESENTATION != PRESENCE",
	"CONFIDENCE != CERTAINTY",
	"RECOVERY > PROPAGATION",
	"HUMAN_AGENCY > MACHINE_AUTHORITY",
	"MACHINE_CAN_READ != MACHINE_CAN_DEFINE",
}

var GLITCHOLOGYGlyphs = []string{
	"(⓿_⓿)",
	"𐦂",
	"﹖",
	"︖",
	"⁇",
	"¿",
	"⸮",
	"‽",
	"？",
	"𖨆",
	"𐀪",
	"𖠋",
	"Err ⃝or⃟⃤",
	"𒈔𒅒𒇫𒄆",
	"☻⃟❦",
	"«ADMIN»",
	"꩜",
	"𓁇𓁋",
	"🫐⃟",
	"𖤍",
}

func GLITCHOLOGYStatement(glyph, runtime, claim, state, authority string) string {
	return fmt.Sprintf("%s [%s] %s :: %s :: %s", glyph, runtime, claim, state, authority)
}

func GLITCHOLOGYBanner(label string) string {
	return fmt.Sprintf("𒄆𓁹✞𒀱✞𓁹𒄆  %s  //  %s  𒄆𓁹✞𒀱✞𓁹𒄆", GLITCHOLOGYDisplayTitle, label)
}
