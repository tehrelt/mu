import { createSession } from "@/actions/session.server";
import { authService } from "@/service/auth.service";

export const PUT = async (req: Request, res: Response) => {
  const token = req.headers.get("authorization")?.split(" ")[1];

  if (!token) {
    return res;
  }

  const newPair = await authService.refresh(token);

  await createSession(newPair);

  return res;
};
