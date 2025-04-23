import z from 'zod'

export const rateCreateSchema = z.object({
    initialRate: z.string().transform(v => parseFloat(v)),
    measureUnit: z.string(),
    name: z.string(),
})

export type RateCreate = z.infer<typeof rateCreateSchema>

