package ssml

import (
	"testing"
	"time"
)

const (
	test1 = `<speak>asdf <break time="1000ms" /> <phoneme alphabet="ipa" ph="d͡ʒð">test</phoneme></speak>`
	test2 = `<speak>asdf2 <say-as interpret-as="date" format="dmy">20170228</say-as></speak>`
)

func TestSSMLBuilder(t *testing.T) {
	r := NewBuilder()
	result := r.Text("asdf").Space().Break(time.Duration(1)*time.Second).Space().Phoneme("test", ALPHABET_IPA, "d͡ʒð").String()
	if result != test1 {
		t.Fatalf("Expected %s, got %s", test1, result)
	}

	r = NewBuilder()
	result = r.Text("asdf2").Space().Date(time.Now(), DATE_FORMAT_DMY).String()

	if result != test2 {
		t.Fatalf("Expected %s, got %s", test2, result)
	}
}
