package ssml

import (
	"bytes"
	"fmt"
	"net/url"
	"time"
)

type Alphabet int

const (
	ALPHABET_IPA Alphabet = iota
	ALPHABET_X_SAMPA
)

func (s Alphabet) String() string {
	switch s {
	case ALPHABET_X_SAMPA:
		return "x-sampa"
	default:
		return "ipa"
	}
}

type Role int

const (
	ROLE_VERB Role = iota
	ROLE_PAST_PARTICIPLE
	ROLE_NOUN
	ROLE_SENSE
)

func (s Role) String() string {
	switch s {
	case ROLE_PAST_PARTICIPLE:
		return "x-ivona:VBD"
	case ROLE_NOUN:
		return "x-ivona:NN"
	case ROLE_SENSE:
		return "x-ivona:SENSE_1"
	default:
		return "ivona:VB"
	}
}

type InterpretAs int

const (
	INTERPRET_AS_CHARACTERS InterpretAs = iota
	INTERPRET_AS_CARDINAL
	INTERPRET_AS_NUMBER
	INTERPRET_AS_ORDINAL
	INTERPRET_AS_DIGITS
	INTERPRET_AS_FRACTION
	INTERPRET_AS_UNIT
	INTERPRET_AS_TIME
	INTERPRET_AS_TELEPHONE
	INTERPRET_AS_ADDRESS
	INTERPRET_AS_INTERJECTION
)

func (s InterpretAs) String() string {
	switch s {
	case INTERPRET_AS_CARDINAL, INTERPRET_AS_NUMBER:
		return "number"
	case INTERPRET_AS_ORDINAL:
		return "ordinal"
	case INTERPRET_AS_DIGITS:
		return "digits"
	case INTERPRET_AS_FRACTION:
		return "fraction"
	case INTERPRET_AS_UNIT:
		return "unit"
	case INTERPRET_AS_TIME:
		return "time"
	case INTERPRET_AS_TELEPHONE:
		return "telephone"
	case INTERPRET_AS_ADDRESS:
		return "address"
	case INTERPRET_AS_INTERJECTION:
		return "interjection"
	default:
		return "characters"
	}
}

type DateFormat int

const (
	DATE_FORMAT_MDY DateFormat = iota
	DATE_FORMAT_DMY
	DATE_FORMAT_YMD
	DATE_FORMAT_MD
	DATE_FORMAT_DM
	DATE_FORMAT_YM
	DATE_FORMAT_MY
	DATE_FORMAT_D
	DATE_FORMAT_M
	DATE_FORMAT_Y
)

func (s DateFormat) String() string {
	switch s {
	case DATE_FORMAT_DMY:
		return "dmy"
	case DATE_FORMAT_YMD:
		return "ymd"
	case DATE_FORMAT_MD:
		return "md"
	case DATE_FORMAT_DM:
		return "dm"
	case DATE_FORMAT_YM:
		return "ym"
	case DATE_FORMAT_MY:
		return "my"
	case DATE_FORMAT_D:
		return "d"
	case DATE_FORMAT_M:
		return "m"
	case DATE_FORMAT_Y:
		return "y"
	default:
		return "mdy"
	}
}

type Builder interface {
	Text(string) Builder
	Space() Builder
	Newline() Builder
	Paragraph(string) Builder
	Sentence(string) Builder
	Break(time.Duration) Builder
	StrongBreak() Builder
	Audio(*url.URL) Builder
	Word(string, Role) Builder
	SayAs(string, InterpretAs) Builder
	Date(time.Time, DateFormat) Builder
	Phoneme(string, Alphabet, string) Builder
	String() string
}

type builder struct {
	buf bytes.Buffer
}

func (r *builder) Text(text string) Builder {
	r.buf.WriteString(text)
	return r
}

func (r *builder) Space() Builder {
	r.buf.WriteByte(' ')
	return r
}

func (r *builder) Newline() Builder {
	r.buf.WriteByte('\n')
	return r
}

func (r *builder) Paragraph(text string) Builder {
	return r.Text("<p>").Text(text).Text("</p>")
}

func (r *builder) Sentence(text string) Builder {
	return r.Text("<s>").Text(text).Text("</s>")
}

func (r *builder) Break(duration time.Duration) Builder {
	return r.Text(fmt.Sprintf(`<break time="%dms" />`, int64(duration.Seconds()*1000)))
}

func (r *builder) StrongBreak() Builder {
	return r.Text(`<break time="strong" />`)
}

func (r *builder) Audio(URL *url.URL) Builder {
	return r.Text(fmt.Sprintf(`<audio src="%s" />`, URL.String()))
}

func (r *builder) Word(text string, role Role) Builder {
	return r.Text(fmt.Sprintf(`<w role="%s" />`, role)).Text(text).Text("</w>")
}

func (r *builder) SayAs(text string, interpretAs InterpretAs) Builder {
	return r.Text(fmt.Sprintf(`<say-as interpret-as="%s">`, interpretAs)).Text(text).Text("</say-as>")
}

func (r *builder) Date(date time.Time, format DateFormat) Builder {
	return r.Text(fmt.Sprintf(`<say-as interpret-as="date" format="%s">`, format)).Text(date.Format("20060102")).Text("</say-as>")
}

func (r *builder) Phoneme(text string, alphabet Alphabet, ph string) Builder {
	return r.Text(fmt.Sprintf(`<phoneme alphabet="%s" ph="%s">`, alphabet, ph)).Text(text).Text("</phoneme>")
}

func (r builder) String() string {
	return fmt.Sprintf("<speak>%s</speak>", r.buf.String())
}

func NewBuilder() Builder {
	return &builder{}
}
