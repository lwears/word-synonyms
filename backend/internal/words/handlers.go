package words

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type WordRequest struct {
	Word string `json:"word"`
}
type SynonymRequest struct {
	Synonym string `json:"synonym"`
}

type AddWordResponse struct {
	ID   int64  `json:"id"`
	Word string `json:"word"`
}

type GetSynonymsResponse struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
}

type GetWordsForSynonymResponse struct {
	Synonym string   `json:"synonym"`
	Words   []string `json:"words"`
}

type AddSynonymResponse struct {
	ID int64 `json:"id"`
}

type WordError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type WordsHTTPHandler struct {
	wordService WordsService
	mux         *http.ServeMux
}

func isValidWord(input string) bool {
	pattern := `^[a-zA-Z]{1,50}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

func normalizeWord(word string) string {
	return strings.ToLower(word)
}

func NewWordsHTTPHandler(wordService WordsService) *WordsHTTPHandler {
	handler := &WordsHTTPHandler{
		wordService: wordService,
		mux:         http.NewServeMux(),
	}
	handler.routes()
	return handler
}

func (h *WordsHTTPHandler) routes() {
	h.mux.HandleFunc("POST /word", h.AddWordHandler)
	h.mux.HandleFunc("GET /words", h.GetAllWordsHandler)
	h.mux.HandleFunc("GET /words/{synonym}", h.GetWordsForSynonymHandler)

	h.mux.HandleFunc("POST /synonym/{word}", h.AddSynonymHandler)
	h.mux.HandleFunc("GET /synonyms/{word}", h.GetSynonymsHandler)
}

func (h *WordsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *WordsHTTPHandler) AddWordHandler(w http.ResponseWriter, r *http.Request) {
	var req WordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !isValidWord(req.Word) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid word value")
		return
	}

	req.Word = normalizeWord(req.Word)

	wordDbRow, err := h.wordService.GetWord(req.Word)
	if err == nil && wordDbRow != nil {
		h.errorResponse(w, http.StatusConflict, "Word already exists")
		return
	}

	newWord, err := h.wordService.AddWord(req.Word)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response := AddWordResponse{
		ID:   newWord.ID,
		Word: newWord.Word,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *WordsHTTPHandler) AddSynonymHandler(w http.ResponseWriter, r *http.Request) {
	var req SynonymRequest
	word := r.PathValue("word")

	// Validate word
	if !isValidWord(word) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid word value")
		return
	}

	word = normalizeWord(word)

	wordDbRow, err := h.wordService.GetOrAddWord(word)
	if err != nil {
		h.errorResponse(w, http.StatusConflict, err.Error())
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !isValidWord(req.Synonym) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid synonym field")
		return
	}

	synonym := normalizeWord(req.Synonym)

	synonymWordDbRow, err := h.wordService.GetOrAddWord(synonym)
	if err != nil {
		h.errorResponse(w, http.StatusConflict, err.Error())
		return
	}

	newWordId, err := h.wordService.AddSynonym(wordDbRow.ID, synonymWordDbRow.ID)
	if err != nil {
		if err, ok := err.(*sqlite.Error); ok && err.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
			h.errorResponse(w, http.StatusConflict, "Synonym link already exists")
			return
		}
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response := AddSynonymResponse{
		ID: newWordId,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *WordsHTTPHandler) GetAllWordsHandler(w http.ResponseWriter, r *http.Request) {
	words, err := h.wordService.GetAll()
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if err := json.NewEncoder(w).Encode(words); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *WordsHTTPHandler) GetSynonymsHandler(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")

	if !isValidWord(word) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid word value")
		return
	}

	word = normalizeWord(word)

	wordDbRow, err := h.wordService.GetWord(word)
	if err != nil || wordDbRow == nil {
		h.errorResponse(w, http.StatusNotFound, "Word does not exist")
		return
	}

	synonyms, err := h.wordService.GetSynonyms(wordDbRow)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response := &GetSynonymsResponse{
		Word:     synonyms.Word,
		Synonyms: synonyms.Synonyms,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *WordsHTTPHandler) GetWordsForSynonymHandler(w http.ResponseWriter, r *http.Request) {
	synonym := r.PathValue("synonym")

	if !isValidWord(synonym) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid synonym value")
		return
	}

	synonym = normalizeWord(synonym)

	wordDbRow, err := h.wordService.GetWord(synonym)
	if err != nil || wordDbRow == nil {
		h.errorResponse(w, http.StatusNotFound, "Word not found")
		return
	}

	synonyms, err := h.wordService.GetWordsForSynonym(wordDbRow)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response := &GetWordsForSynonymResponse{
		Synonym: synonyms.Synonym,
		Words:   synonyms.Words,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *WordsHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(WordError{
		StatusCode: statusCode,
		Error:      errorString,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
