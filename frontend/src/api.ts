import ky from 'ky'
import env from './env'
import type {
  AddWordResponse,
  AddSynonymResponse,
  GetSynonymsResponse,
  GetWordsForSynonymResponse,
  FormData,
} from './types'

const addWord = (data: FormData) =>
  ky.post(`${env.VITE_BASE_URI}/word`, { json: data }).json<AddWordResponse>()

const addSynonymToWord = ({ synonym, word }: FormData) =>
  ky
    .post(`${env.VITE_BASE_URI}/synonym/${word}`, { json: { synonym } })
    .json<AddSynonymResponse>()

const getSynonyms = (word: string) =>
  ky.get(`${env.VITE_BASE_URI}/synonyms/${word}`).json<GetSynonymsResponse>()

const getWordsForSynonym = (synonym: string) =>
  ky
    .get(`${env.VITE_BASE_URI}/words/${synonym}`)
    .json<GetWordsForSynonymResponse>()

export const api = {
  addWord,
  addSynonymToWord,
  getSynonyms,
  getWordsForSynonym,
}
