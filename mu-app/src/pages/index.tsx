import { Card, CardContent } from "@/components/ui/card";
import { BanknoteArrowUp, Bell, Rocket, Smartphone } from "lucide-react";

export const Index = () => {
  return (
    <div className="flex justify-center w-screen mt-8">
      <div className="flex justify-around gap-x-12 items-center">
        <div className="grid gap-y-2">
          <Feature
            icon={<Smartphone />}
            text="Одно приложение для разных услуг"
          />
          <Feature
            icon={<BanknoteArrowUp />}
            text="Оплачивайте прям в приложении"
          />
          <Feature
            icon={<Bell />}
            text="Получайте уведомления в Telegram или на почту"
          />
        </div>
        <div>
          <h1 className="font-bold text-xl">МОИ УСЛУГИ</h1>
          <h1 className="">Единый личный кабинет ЖКХ-услуг</h1>
        </div>
      </div>
    </div>
  );
};

const Feature = ({ text, icon }: { text: string; icon: React.ReactNode }) => {
  return (
    <div className="max-w-lg">
      <Card>
        <CardContent className="grid grid-cols-9 grid-rows-1">
          <p className="col-span-7">{text}</p>
          <div className="p-2 rounded-full bg-muted col-start-9 flex items-center justify-center ">
            {icon}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
