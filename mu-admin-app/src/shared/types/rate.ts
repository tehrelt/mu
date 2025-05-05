import z from "zod";

export const serviceTypeSchema = z.enum([
  "WATER_SUPPLY",
  "HEATING",
  "POWER_SUPPLY",
  "GAS_SUPPLY",
]);
export type ServiceType = z.infer<typeof serviceTypeSchema>;

export const rateCreateSchema = z.object({
  initialRate: z.string().transform((v) => parseFloat(v)),
  measureUnit: z.string(),
  name: z.string(),
  serviceType: serviceTypeSchema,
});

export type RateCreate = z.infer<typeof rateCreateSchema>;

export const rateSchema = z.object({
  id: z.string(),
  rate: z.number().multipleOf(0.01),
  measureUnit: z.string(),
  name: z.string(),
  serviceType: serviceTypeSchema,
});

export type Rate = z.infer<typeof rateSchema>;

export const rateListSchema = z.object({
  rates: z.array(rateSchema),
});

export type RateList = z.infer<typeof rateListSchema>;

export const localizeServiceType = (serviceType: ServiceType) => {
  switch (serviceType) {
    case "WATER_SUPPLY":
      return "Водоснабжение";
    case "HEATING":
      return "Отопление";
    case "POWER_SUPPLY":
      return "Электроснабжение";
    case "GAS_SUPPLY":
      return "Газоснабжение";
  }
};
