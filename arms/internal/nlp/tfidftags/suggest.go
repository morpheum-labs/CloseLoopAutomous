// Package tfidftags scores tokens in a target string with TF-IDF against a text corpus (no LLM).
package tfidftags

import (
	"math"
	"sort"
	"strings"
	"unicode"

	tfidf "github.com/go-nlp/tfidf"
)

// TagScore is one suggested tag with a relative TF-IDF salience score.
type TagScore struct {
	Token string  `json:"token"`
	Score float64 `json:"score"`
}

type intDoc []int

func (d intDoc) IDs() []int { return []int(d) }

const (
	defaultTopK        = 12
	maxTopK            = 64
	maxCorpusStrings  = 2000
	maxCharsPerString = 65536
)

var englishStopwords = map[string]struct{}{
	"a": {}, "an": {}, "the": {}, "and": {}, "or": {}, "but": {}, "in": {}, "on": {}, "at": {}, "to": {}, "for": {},
	"of": {}, "as": {}, "by": {}, "with": {}, "from": {}, "is": {}, "are": {}, "was": {}, "were": {}, "be": {},
	"been": {}, "being": {}, "have": {}, "has": {}, "had": {}, "do": {}, "does": {}, "did": {}, "will": {},
	"would": {}, "could": {}, "should": {}, "may": {}, "might": {}, "must": {}, "shall": {}, "can": {}, "this": {},
	"that": {}, "these": {}, "those": {}, "it": {}, "its": {}, "i": {}, "you": {}, "he": {}, "she": {}, "we": {},
	"they": {}, "them": {}, "their": {}, "what": {}, "which": {}, "who": {}, "whom": {}, "when": {}, "where": {},
	"why": {}, "how": {}, "all": {}, "each": {}, "every": {}, "both": {}, "few": {}, "more": {}, "most": {},
	"other": {}, "some": {}, "such": {}, "no": {}, "nor": {}, "not": {}, "only": {}, "own": {}, "same": {},
	"so": {}, "than": {}, "too": {}, "very": {}, "just": {}, "also": {}, "into": {}, "about": {}, "over": {},
	"after": {}, "before": {}, "between": {}, "through": {}, "during": {}, "without": {}, "within": {}, "again": {},
	"further": {}, "then": {}, "once": {}, "here": {}, "there": {}, "any": {}, "if": {}, "because": {},
}

// Tokenize lowercases and splits on non-letter/non-digit runes.
func Tokenize(s string) []string {
	s = strings.ToLower(s)
	var cur strings.Builder
	var out []string
	flush := func() {
		if cur.Len() == 0 {
			return
		}
		tok := cur.String()
		cur.Reset()
		if _, stop := englishStopwords[tok]; stop {
			return
		}
		out = append(out, tok)
	}
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			cur.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()
	return out
}

func clampStr(s string) string {
	if len(s) <= maxCharsPerString {
		return s
	}
	return s[:maxCharsPerString]
}

func buildVocab(tokens [][]string) ([]string, map[string]int) {
	seen := make(map[string]struct{}, 256)
	order := make([]string, 0, 256)
	for _, toks := range tokens {
		for _, t := range toks {
			if t == "" {
				continue
			}
			if _, ok := seen[t]; ok {
				continue
			}
			seen[t] = struct{}{}
			order = append(order, t)
		}
	}
	v2i := make(map[string]int, len(order))
	for i, w := range order {
		v2i[w] = i
	}
	return order, v2i
}

func toIntDoc(toks []string, v2i map[string]int) intDoc {
	if len(toks) == 0 {
		return nil
	}
	out := make([]int, len(toks))
	for i, t := range toks {
		out[i] = v2i[t]
	}
	return intDoc(out)
}

func freqSuggest(tokens []string, minLen, topK int) []TagScore {
	count := make(map[string]int)
	for _, t := range tokens {
		if len(t) < minLen {
			continue
		}
		count[t]++
	}
	out := make([]TagScore, 0, len(count))
	for tok, n := range count {
		out = append(out, TagScore{Token: tok, Score: float64(n)})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Token < out[j].Token
	})
	if topK > 0 && len(out) > topK {
		out = out[:topK]
	}
	return out
}

// Suggest returns up to topK tokens from targetText ranked by TF-IDF against corpus strings.
// Corpus documents are not included in the IDF model as the scored document; pass other docs only.
// If corpus is empty (or all empty after tokenize), falls back to term frequency in the target.
// minTokenLen filters suggested tokens (default 2 if <= 0).
func Suggest(corpus []string, targetText string, topK, minTokenLen int) []TagScore {
	if minTokenLen <= 0 {
		minTokenLen = 2
	}
	if topK <= 0 {
		topK = defaultTopK
	}
	if topK > maxTopK {
		topK = maxTopK
	}

	targetText = clampStr(strings.TrimSpace(targetText))
	targetToks := Tokenize(targetText)
	filteredTarget := make([]string, 0, len(targetToks))
	for _, t := range targetToks {
		if len(t) >= minTokenLen {
			filteredTarget = append(filteredTarget, t)
		}
	}
	if len(filteredTarget) == 0 {
		return nil
	}

	var corpToks [][]string
	nCorpus := 0
	for _, c := range corpus {
		if nCorpus >= maxCorpusStrings {
			break
		}
		c = clampStr(strings.TrimSpace(c))
		if c == "" {
			continue
		}
		toks := Tokenize(c)
		if len(toks) == 0 {
			continue
		}
		corpToks = append(corpToks, toks)
		nCorpus++
	}

	if len(corpToks) == 0 {
		return freqSuggest(filteredTarget, minTokenLen, topK)
	}

	allToks := append([][]string{}, corpToks...)
	allToks = append(allToks, filteredTarget)
	_, v2i := buildVocab(allToks)

	tf := tfidf.New()
	for _, toks := range corpToks {
		tf.Add(toIntDoc(toks, v2i))
	}
	if tf.Docs == 0 {
		return freqSuggest(filteredTarget, minTokenLen, topK)
	}
	tf.CalculateIDF()

	targetDoc := toIntDoc(filteredTarget, v2i)
	scores := tf.Score(targetDoc)
	// Aggregate max score per token id (targetDoc order aligns with scores).
	best := make(map[int]float64)
	for i, id := range []int(targetDoc) {
		if i >= len(scores) {
			break
		}
		s := scores[i]
		if math.IsNaN(s) || math.IsInf(s, 0) {
			continue
		}
		if prev, ok := best[id]; !ok || s > prev {
			best[id] = s
		}
	}

	inv := make([]string, len(v2i))
	for w, id := range v2i {
		inv[id] = w
	}
	out := make([]TagScore, 0, len(best))
	for id, s := range best {
		tok := inv[id]
		if len(tok) < minTokenLen {
			continue
		}
		out = append(out, TagScore{Token: tok, Score: s})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Token < out[j].Token
	})
	if len(out) > topK {
		out = out[:topK]
	}
	return out
}
