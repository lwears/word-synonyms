import { z } from 'zod'

const envVars = z.object({
  VITE_BASE_URI: z.string().url(),
})

const env = envVars.parse(import.meta.env)

export default env
