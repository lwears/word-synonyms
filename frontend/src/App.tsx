import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { Toaster } from 'sonner';
import { toast } from 'sonner';
import { zodResolver } from '@hookform/resolvers/zod';
import ky from 'ky';

import WordWithSynonyms from './components/WordWithSynonyms';
import SynonymWithWords from './components/SynonymWithWords';
import Loading from './components/Loading';
import Button from './components/Button';
import Input from './components/Input';
import { schema } from './schema';
import { handleHttpError } from './helpers';

import type {
  GetSynonymsResponse,
  GetWordsForSynonymResponse,
  AddWordResponse,
  AddSynonymResponse,
  FormData,
} from './types';

const BASE_URL = 'http://localhost:8090';

function App() {
  const [synonyms, setSynonyms] = useState<GetSynonymsResponse | null>(null);
  const [synonymWithWords, setSynonymWithWords] =
    useState<GetWordsForSynonymResponse | null>(null);
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  useEffect(() => {
    if (loading) {
      setSynonyms(null);
      setSynonymWithWords(null);
    }
  }, [loading]);

  const watchedSynonym = watch('synonym');
  const watchedword = watch('word');

  const handleAddWord = (data: FormData) => {
    setLoading(true);
    return ky
      .post(`${BASE_URL}/word`, { json: data })
      .json<AddWordResponse>()
      .then((d) => {
        reset();
        toast.success('Word Added', { description: d.word });
      })
      .catch((err) => handleHttpError(err, `Error Adding Word: ${data.word}`))
      .finally(() => setLoading(false));
  };

  const handleAddSynonymToWord = ({ synonym, word }: FormData) => {
    setLoading(true);
    return ky
      .post(`${BASE_URL}/synonym/${word}`, { json: { synonym } })
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
      .get(`${BASE_URL}/synonyms/${data.word}`)
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
      .get(`${BASE_URL}/words/${data.synonym}`)
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
              <Input
                title="Enter Word"
                type="text"
                placeholder="Type word here"
                name="word"
                register={register}
                error={errors.word}
                required
                disabled={loading}
              />
              <Input
                title="Enter Synonym"
                type="text"
                placeholder="Type synonym here"
                name="synonym"
                register={register}
                error={errors.synonym}
                required
                disabled={loading}
              />
              <Button
                onClick={handleSubmit(handleAddWord)}
                content="Add Word"
                loading={loading}
                disabled={!watchedword}
              />
              <Button
                onClick={handleSubmit(handleAddSynonymToWord)}
                content="Add Synonym to Word"
                loading={loading}
                disabled={!watchedSynonym && !watchedword}
              />
              <Button
                content="Get Synonyms"
                onClick={handleSubmit(handleGetSynonyms)}
                loading={loading}
                disabled={!watchedword}
              />
              <Button
                content="Get Words for Synonym"
                onClick={handleSubmit(handleGetWordsForSynonym)}
                loading={loading}
                disabled={!watchedSynonym}
              />
            </div>
          </form>
        </div>
        <div className="flex border border-primary rounded flex-1 p-4 justify-center max-w-[500px]">
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
