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

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate word
	if !isValidWord(req.Word) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid word value")
		return
	}

	// Normalize word to lowercase
	req.Word = strings.ToLower(req.Word)

	// Check word exists
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

	// Create the response to ensure keys are lowercase
	response := AddWordResponse{
		ID:   newWord.ID,
		Word: newWord.Word,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *WordsHTTPHandler) AddSynonymHandler(w http.ResponseWriter, r *http.Request) {
	var req SynonymRequest
	word := r.PathValue("word")

	// Validate word
	if !isValidWord(word) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid synonym value")
		return
	}

	// Normalise
	word = strings.ToLower(word)

	// Check word exists and create if not
	wordDbRow, err := h.wordService.GetOrAddWord(word)
	if err != nil {
		h.errorResponse(w, http.StatusConflict, err.Error())
		return
	}

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate synonym word
	if !isValidWord(req.Synonym) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid synonym field")
		return
	}

	// Normalise
	synonym := strings.ToLower(req.Synonym)

	// Check synonym word exists and create if not
	synonymWordDbRow, err := h.wordService.GetOrAddWord(synonym)
	if err != nil {
		h.errorResponse(w, http.StatusConflict, err.Error())
		return
	}

	newWordId, err := h.wordService.AddSynonym(wordDbRow.ID, synonymWordDbRow.ID)
	if err != nil {
		// This could be done better. Running out of time and need to handle conflict
		// Also would need to handle all other error codes potentially
		if err, ok := err.(*sqlite.Error); ok {
			if err.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
				h.errorResponse(w, http.StatusConflict, "Synonym link already exists")
				return
			}
		}
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// Create the response to ensure keys are lowercase
	response := AddSynonymResponse{
		ID: newWordId,
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

// This was just added for testing
func (h *WordsHTTPHandler) GetAllWordsHandler(w http.ResponseWriter, r *http.Request) {
	words, err := h.wordService.GetAll()
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = json.NewEncoder(w).Encode(words)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *WordsHTTPHandler) GetSynonymsHandler(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")

	// validate word
	if !isValidWord(word) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid word value")
		return
	}

	word = strings.ToLower(word)

	// Check word exists
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

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (h *WordsHTTPHandler) GetWordsForSynonymHandler(w http.ResponseWriter, r *http.Request) {
	synonym := r.PathValue("synonym")

	// validate word
	if !isValidWord(synonym) {
		h.errorResponse(w, http.StatusBadRequest, "Invalid synonym value")
		return
	}

	synonym = strings.ToLower(synonym)

	// Check word exists
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

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *WordsHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(WordError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
