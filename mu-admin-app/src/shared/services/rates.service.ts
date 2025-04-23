import { api } from "@/app/api"
import { RateCreate } from "../types/rate"

class RateService {
    async create(data: RateCreate) {
        const response = await api.post('/rates', data)
        return response.data
    }
}

export const rateService = new RateService()