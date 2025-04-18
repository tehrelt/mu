import { useQuery } from "@tanstack/react-query";
import { authService } from "../services/auth.service";

export const useProfile = () => {
  const data = useQuery({
    queryKey: ["profile"],
    queryFn: async () => {
      const profile = await authService.profile();
      return profile;
    },
  });

  return data;
};
