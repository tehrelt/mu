import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useProfile } from "@/shared/hooks/use-profile";
import { integrationService } from "@/shared/services/integrations.service";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Send } from "lucide-react";

export const IntegrationsSettingsPage = () => {
  const profileQuery = useProfile();

  const settingsQuery = useQuery({
    queryKey: ["profile", "integrations", "settings"],
    queryFn: async () => await integrationService.settings(),
  });

  return (
    <div>
      <Card>
        <CardHeader>Настройка интеграций</CardHeader>
        <CardContent className="grid gap-6">
          <div className="grid gap-2">
            <Label>Электронная почта</Label>
            <Input value={profileQuery.data ? profileQuery.data.email : ""} />
          </div>

          <div className="grid gap-2">
            <Label>Telegram</Label>
            <div>
              {settingsQuery.data?.hasTelegram ? (
                <Button disabled>Телеграм подключен</Button>
              ) : (
                <ConnectTelegramButton />
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

const ConnectTelegramButton = () => {
  const mutation = useMutation({
    mutationKey: ["integrations", "link-telegram"],
    mutationFn: async () => {
      const { otp, userId } = await integrationService.getOtpCode();
      window.open(`https://t.me/mu_myservices_bot?start=${otp}_${userId}`);
    },
  });

  const handleConnectClick = () => {
    mutation.mutate();
  };

  return (
    <Button onClick={() => handleConnectClick()}>
      <Send /> Подключить
    </Button>
  );
};
