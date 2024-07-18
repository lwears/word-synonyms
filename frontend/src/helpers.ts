import { toast } from 'sonner'

import type { HTTPError } from 'ky'

export const handleHttpError = (err: HTTPError, text: string) =>
  err.response
    .json()
    .then(({ error }) => toast.error(text, { description: error }))
