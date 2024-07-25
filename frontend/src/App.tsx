import { Toaster, toast } from 'sonner'
import { Form } from './Form'
import {
  QueryClient,
  QueryClientProvider,
  QueryCache,
} from '@tanstack/react-query'
import React from 'react'

function App() {
  const [queryClient] = React.useState(
    () =>
      new QueryClient({
        queryCache: new QueryCache({
          onError: (error, query) => {
            if (query.meta && query.meta.errorMessage) {
              toast.error(query.meta.errorMessage as string)
            }
          },
        }),
      })
  )

  return (
    <QueryClientProvider client={queryClient}>
      <Form />
      <Toaster position="top-right" expand richColors />
    </QueryClientProvider>
  )
}

export default App
