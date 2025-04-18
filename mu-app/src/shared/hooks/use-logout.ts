import { authService } from "@/shared/services/auth.service";
import { sessionService } from "@/shared/services/session.service";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const useLogout = () => {
  const client = useQueryClient();

  return useMutation({
    mutationKey: ["logout"],
    mutationFn: async () => await authService.logout(),
    onSuccess: async () => {
      await client.invalidateQueries({
        queryKey: ["profile"],
      });

      sessionService.clear();
    },
  });
};
