import z from 'zod'

export const rateCreateSchema = z.object({
    initialRate: z.string().transform(v => parseFloat(v)),
    measureUnit: z.string(),
    name: z.string(),
})

export type RateCreate = z.infer<typeof rateCreateSchema>

export const rateSchema = z.object({
    id: z.string(),
    rate: z.number().multipleOf(0.01),
    measureUnit: z.string(),
    name: z.string(),
})

export type Rate = z.infer<typeof rateSchema>

export const rateListSchema = z.object({
    rates: z.array(rateSchema),
})

export type RateList = z.infer<typeof rateListSchema>