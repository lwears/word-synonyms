import { z } from 'zod'

const bothSchema = z.object({
  word: z
    .string({
      required_error: 'word is required',
      invalid_type_error: 'word must be a string',
    })
    .toLowerCase()
    .trim()
    .min(1, { message: 'word must be minimum 1 character' })
    .max(45)
    .refine((value) => isNaN(Number(value)), {
      message: 'value should not be a number',
    }),
  synonym: z
    .string({
      required_error: 'synonym is required',
      invalid_type_error: 'synonym must be a string',
    })
    .toLowerCase()
    .trim()
    .min(1, { message: 'synonym must be minimum 1 character' })
    .max(45)
    .refine((value) => isNaN(Number(value)), {
      message: 'value should not be a number',
    }),
})

const wordOnlySchema = z.object({
  word: z
    .string({
      required_error: 'word is required',
      invalid_type_error: 'word must be a string',
    })
    .toLowerCase()
    .trim()
    .min(1, { message: 'word must be minimum 1 character' })
    .max(45)
    .refine((value) => isNaN(Number(value)), {
      message: 'value should not be a number',
    }),
  synonym: z.literal(''),
})

const synonymOnlySchema = z.object({
  word: z.literal(''),
  synonym: z
    .string({
      required_error: 'synonym is required',
      invalid_type_error: 'synonym must be a string',
    })
    .toLowerCase()
    .trim()
    .min(1, { message: 'synonym must be minimum 1 character' })
    .max(50, { message: 'synonym must be maximum 50 characters' })
    .refine((value) => isNaN(Number(value)), {
      message: 'value should not be a number',
    }),
})

export const schema = z.union([bothSchema, wordOnlySchema, synonymOnlySchema])
