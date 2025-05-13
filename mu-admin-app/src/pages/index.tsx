import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { routes } from "@/shared/routes";
import { Link } from "react-router-dom";

export const Index = () => {
  return (
    <div className="grid gap-4">
      <h1 className="text-xl font-bold">Главная</h1>
      <div className="flex gap-2">
        <Link to={routes.rate.list}>
          <Card className="group">
            <CardContent>
              <CardTitle className="group-hover:underline">Тарифы</CardTitle>
            </CardContent>
          </Card>
        </Link>
        <Link to={routes.users.list}>
          <Card className="group">
            <CardContent>
              <CardTitle className="group-hover:underline">
                Пользователи
              </CardTitle>
            </CardContent>
          </Card>
        </Link>
        <Link to={routes.tickets.list}>
          <Card className="group">
            <CardContent>
              <CardTitle className="group-hover:underline">Заявки</CardTitle>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  );
};
