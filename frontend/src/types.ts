export interface FormData {
  word: string;
  synonym: string;
}

export interface AddWordResponse {
  word: string;
  synonym: string;
}

export interface AddSynonymResponse {
  id: number;
}

export interface GetSynonymsResponse {
  word: string;
  synonyms: string[];
}

export interface GetWordsForSynonymResponse {
  synonym: string;
  words: string[];
}
