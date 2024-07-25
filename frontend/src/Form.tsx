import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { Toaster, toast } from 'sonner'

import Button from './components/Button'
import Input from './components/Input'
import Loading from './components/Loading'
import SynonymWithWords from './components/SynonymWithWords'
import WordWithSynonyms from './components/WordWithSynonyms'
import { schema } from './schema'
import { api } from './api'

import type { FormData } from './types'
import { useEffect } from 'react'

export const Form = () => {
  const queryClient = useQueryClient()

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const watchedSynonym = watch('synonym')
  const watchedword = watch('word')

  const mutateAddWord = useMutation({
    mutationFn: api.addWord,
    onSuccess: (d) => {
      reset()
      toast.success('Word Added', { description: d.word })
      queryClient.resetQueries()
    },
    onError: (error, data) => {
      toast.error(`Error Adding Word: ${data.word}`, {
        description: error.message,
      })
    },
  })

  const mutateAddSynonym = useMutation({
    mutationFn: api.addSynonymToWord,
    onSuccess: (res, data) => {
      reset()
      toast.success(`Synonym added to ${data.word}`, { description: res.id })
      queryClient.resetQueries()
    },
    onError: (error, data) => {
      toast.error(`Error Adding Synonym: ${data.synonym}`, {
        description: error.message,
      })
    },
  })

  const getSynonyms = () => {
    queryClient.resetQueries({ queryKey: ['words'] })
    return api.getSynonyms(watchedword)
  }

  const getWordsForSynonym = () => {
    queryClient.resetQueries({ queryKey: ['synonyms'] })
    return api.getWordsForSynonym(watchedSynonym)
  }

  const queryGetSynonyms = useQuery({
    queryKey: ['synonyms'],
    queryFn: getSynonyms,
    enabled: false,
    meta: {
      errorMessage: 'Error fetching Synonyms',
    },
  })

  const queryGetWordsForSynonym = useQuery({
    queryKey: ['words'],
    queryFn: getWordsForSynonym,
    enabled: false,
    meta: {
      errorMessage: 'Error fetching words for synonym',
    },
  })

  useEffect(() => {
    if (queryGetSynonyms.isSuccess || queryGetWordsForSynonym.isSuccess) {
      reset()
    }
  }, [queryGetSynonyms.isSuccess, queryGetWordsForSynonym.isSuccess, reset])

  const loading =
    mutateAddWord.isPending ||
    mutateAddSynonym.isPending ||
    queryGetSynonyms.isLoading ||
    queryGetWordsForSynonym.isLoading

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
                onClick={handleSubmit((data) => mutateAddWord.mutate(data))}
                content="Add Word"
                loading={loading}
                disabled={!watchedword}
              />
              <Button
                onClick={handleSubmit((data) => mutateAddSynonym.mutate(data))}
                content="Add Synonym to Word"
                loading={loading}
                disabled={!watchedSynonym || !watchedword}
              />
              <Button
                content="Get Synonyms"
                onClick={handleSubmit(() => queryGetSynonyms.refetch())}
                loading={loading}
                disabled={!watchedword}
              />
              <Button
                content="Get Words for Synonym"
                onClick={handleSubmit(() => queryGetWordsForSynonym.refetch())}
                loading={loading}
                disabled={!watchedSynonym}
              />
            </div>
          </form>
        </div>
        <div className="flex border border-primary rounded-lg flex-1 p-4 justify-center max-w-[500px]">
          {loading && <Loading />}
          {queryGetWordsForSynonym.data && (
            <SynonymWithWords {...queryGetWordsForSynonym.data} />
          )}
          {queryGetSynonyms.data && (
            <WordWithSynonyms {...queryGetSynonyms.data} />
          )}
        </div>
      </div>
      <Toaster position="top-right" expand richColors />
    </main>
  )
}
