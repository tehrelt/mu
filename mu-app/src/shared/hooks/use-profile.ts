import { useQuery } from "@tanstack/react-query";
import { authService } from "../services/auth.service";

export const useProfile = () => {
  const data = useQuery({
    queryKey: ["profile"],
    queryFn: () => authService.profile(),
    retry: false,
  });

  return data;
};
