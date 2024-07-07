import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { Toaster } from 'sonner';
import { toast } from 'sonner';
import { zodResolver } from '@hookform/resolvers/zod';
import ky from 'ky';
import clsx from 'clsx';

import WordWithSynonyms from './components/WordWithSynonyms';
import SynonymWithWords from './components/SynonymWithWords';
import Loading from './components/Loading';
import Button from './components/Button';
import { schema } from './schema';

import type { HTTPError } from 'ky';
import type {
  GetSynonymsResponse,
  GetWordsForSynonymResponse,
  AddWordResponse,
  AddSynonymResponse,
  FormData,
} from './types';

const handleHttpError = (err: HTTPError, text: string) =>
  err.response
    .json()
    .then(({ error }) => toast.error(text, { description: error }));

function App() {
  const [synonyms, setSynonyms] = useState<GetSynonymsResponse | null>(null);
  const [synonymWithWords, setSynonymWithWords] =
    useState<GetWordsForSynonymResponse | null>(null);
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setError,
  } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  useEffect(() => {
    if (loading) {
      setSynonyms(null);
      setSynonymWithWords(null);
    }
  }, [loading]);

  const handleAddWord = (data: FormData) => {
    setLoading(true);
    setTimeout(() => {
      return ky
        .post('http://localhost:8090/word', { json: data })
        .json<AddWordResponse>()
        .then((d) => {
          reset();
          toast.success('Word Added', { description: d.word });
        })
        .catch((err) => handleHttpError(err, `Error Adding Word: ${data.word}`))
        .finally(() => setLoading(false));
    }, 5000);
  };

  const handleAddSynonymToWord = ({ synonym, word }: FormData) => {
    setLoading(true);
    // TODO: fix
    if (!synonym && !word) {
      setError('root', { message: 'Please fill in both fields' });
      return;
    }
    return ky
      .post(`http://localhost:8090/synonym/${word}`, { json: { synonym } })
      .json<AddSynonymResponse>()
      .then(() => {
        reset();
        toast.success(`Synonym added to ${word}`, {
          description: `synonym: ${synonym}`,
        });
      })
      .catch((err) =>
        handleHttpError(err, `Error adding synonym to word: ${word}`)
      )
      .finally(() => setLoading(false));
  };

  const handleGetSynonyms = (data: { word: string }) => {
    setLoading(true);
    return ky
      .get(`http://localhost:8090/synonyms/${data.word}`)
      .json<GetSynonymsResponse>()
      .then((d) => {
        reset();
        setSynonyms(d);
      })
      .catch((err) => handleHttpError(err, `Error fetching synonyms`))
      .finally(() => setLoading(false));
  };

  const handleGetWordsForSynonym = (data: { synonym: string }) => {
    setLoading(true);
    return ky
      .get(`http://localhost:8090/words/${data.synonym}`)
      .json<GetWordsForSynonymResponse>()
      .then((d) => {
        reset();
        setSynonymWithWords(d);
      })
      .catch((err) => handleHttpError(err, `Error fetching words for synonyms`))
      .finally(() => setLoading(false));
  };

  return (
    <main className="size-full flex flex-col gap-8 justify-center items-center ">
      <h1 className="text-white text-4xl">Word-Synonyms</h1>
      <div className="grid grid-cols-2 gap-4">
        <div className="flex flex-1 justify-center items-center">
          <form>
            <div className="flex flex-col gap-2">
              <label className="form-control w-full max-w-xs">
                <div className="label">
                  <span className="label-text">Enter Word</span>
                </div>
                <input
                  id="word"
                  type="text"
                  placeholder="Type word here"
                  className={clsx(
                    'input input-bordered w-full max-w-xs',
                    errors.word && 'input-error'
                  )}
                  required
                  disabled={loading}
                  {...register('word')}
                />
                {errors.word && (
                  <div className="label">
                    <span className="label-text-alt text-error">
                      {errors.word.message}
                    </span>
                  </div>
                )}
              </label>
              <label className="form-control w-full max-w-xs">
                <div className="label">
                  <span className="label-text">Enter Synonym</span>
                </div>
                <input
                  id="synonym"
                  type="text"
                  placeholder="Type synonym here"
                  className={clsx(
                    'input input-bordered w-full max-w-xs',
                    errors.synonym && 'input-error'
                  )}
                  required
                  disabled={loading}
                  {...register('synonym')}
                />
                <div className="label">
                  {errors.synonym && (
                    <span className="label-text-alt text-error">
                      {errors.synonym.message}
                    </span>
                  )}
                </div>
              </label>
              <Button
                onClick={handleSubmit(handleAddWord)}
                content="Add Word"
                loading={loading}
              />
              <Button
                content="Get Synonyms"
                onClick={handleSubmit(handleGetSynonyms)}
                loading={loading}
              />
              <Button
                content="Get Words for Synonym"
                onClick={handleSubmit(handleGetWordsForSynonym)}
                loading={loading}
              />
              <Button
                onClick={handleSubmit(handleAddSynonymToWord)}
                content="Add Synonym to Word"
                loading={loading}
              />
            </div>
          </form>
        </div>
        <div className="flex border border-primary rounded flex-1 p-4 justify-center w-[500px]">
          {loading && <Loading />}
          {synonymWithWords && <SynonymWithWords {...synonymWithWords} />}
          {synonyms && <WordWithSynonyms {...synonyms} />}
        </div>
      </div>

      <Toaster position="top-right" expand richColors />
    </main>
  );
}

export default App;
