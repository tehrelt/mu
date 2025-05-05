import { Link } from "react-router-dom";
import { useProfile } from "../shared/hooks/use-profile";
import { Button } from "./ui/button";
import { UserCard } from "./usercard";
import { GalleryVerticalEnd } from "lucide-react";
import { routes } from "@/shared/routes";

export const Header = () => {
  const profileData = useProfile();

  return (
    <header className="flex justify-around py-4 ">
      <div className="flex items-center gap-x-2">
        <div className="flex items-center justify-center rounded-md bg-primary text-primary-foreground p-1">
          <GalleryVerticalEnd className="size-6" />
        </div>
        <h1 className="text-xl font-bold">Мои услуги</h1>
      </div>
      <div className="flex items-center">
        {profileData.data ? (
          <div className="flex gap-x-1 items-center">
            <UserCard user={profileData.data} />
          </div>
        ) : (
          <Link to={routes.signIn}>
            <Button>Авторизация</Button>
          </Link>
        )}
      </div>
    </header>
  );
};
