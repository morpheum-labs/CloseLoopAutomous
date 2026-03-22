package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/closeloopautomous/arms/internal/domain"
	"github.com/closeloopautomous/arms/internal/nlp/tfidftags"
)

type postTfidfSuggestTagsReq struct {
	Corpus       []string `json:"corpus"`
	Text         string   `json:"text"`
	TopK         int      `json:"top_k"`
	MinTokenLen  int      `json:"min_token_len"`
}

type postProductTfidfSuggestTagsReq struct {
	Text         string   `json:"text"`
	IdeaID       string   `json:"idea_id"`
	ExtraCorpus  []string `json:"extra_corpus"`
	TopK         int      `json:"top_k"`
	MinTokenLen  int      `json:"min_token_len"`
}

func (h *Handlers) postTfidfSuggestTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "POST required")
		return
	}
	var req postTfidfSuggestTagsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "text is required")
		return
	}
	tags := tfidftags.Suggest(req.Corpus, req.Text, req.TopK, req.MinTokenLen)
	writeJSON(w, http.StatusOK, map[string]any{
		"tags":               tags,
		"method":             tfidfMethodFromCorpusCount(countNonEmptyCorpus(req.Corpus)),
		"corpus_documents":   countNonEmptyCorpus(req.Corpus),
	})
}

func (h *Handlers) postProductTfidfSuggestTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "POST required")
		return
	}
	pid := domain.ProductID(r.PathValue("id"))
	if _, err := h.Autopilot.Products.ByID(r.Context(), pid); err != nil {
		if mapDomainErr(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	var req postProductTfidfSuggestTagsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	ideas, err := h.Autopilot.Ideas.ListByProduct(r.Context(), pid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	var target string
	var exclude domain.IdeaID
	if id := strings.TrimSpace(req.IdeaID); id != "" {
		exclude = domain.IdeaID(id)
		idea, ierr := h.Autopilot.Ideas.ByID(r.Context(), exclude)
		if ierr != nil {
			if mapDomainErr(w, ierr) {
				return
			}
			writeError(w, http.StatusInternalServerError, "internal", ierr.Error())
			return
		}
		if idea.ProductID != pid {
			writeError(w, http.StatusBadRequest, "bad_request", "idea does not belong to this product")
			return
		}
		target = ideaText(idea)
	} else {
		target = strings.TrimSpace(req.Text)
	}
	if target == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "text or idea_id is required")
		return
	}

	corpus := make([]string, 0, len(ideas)+len(req.ExtraCorpus))
	for i := range ideas {
		if exclude != "" && ideas[i].ID == exclude {
			continue
		}
		corpus = append(corpus, ideaText(&ideas[i]))
	}
	for _, s := range req.ExtraCorpus {
		if t := strings.TrimSpace(s); t != "" {
			corpus = append(corpus, t)
		}
	}
	tags := tfidftags.Suggest(corpus, target, req.TopK, req.MinTokenLen)
	writeJSON(w, http.StatusOK, map[string]any{
		"tags":             tags,
		"method":           tfidfMethodFromCorpusCount(len(corpus)),
		"corpus_documents": len(corpus),
		"product_id":       string(pid),
		"idea_id":          strings.TrimSpace(req.IdeaID),
	})
}

func ideaText(idea *domain.Idea) string {
	var b strings.Builder
	b.WriteString(idea.Title)
	b.WriteByte(' ')
	b.WriteString(idea.Description)
	b.WriteByte(' ')
	b.WriteString(idea.Reasoning)
	b.WriteByte(' ')
	b.WriteString(strings.Join(idea.Tags, " "))
	return strings.TrimSpace(b.String())
}

func countNonEmptyCorpus(c []string) int {
	n := 0
	for _, s := range c {
		if strings.TrimSpace(s) != "" {
			n++
		}
	}
	return n
}

func tfidfMethodFromCorpusCount(usableCorpusDocs int) string {
	if usableCorpusDocs == 0 {
		return "frequency"
	}
	return "tfidf"
}
