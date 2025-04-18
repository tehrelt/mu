import { Link } from "react-router-dom";
import { useProfile } from "../shared/hooks/use-profile";
import { Button } from "./ui/button";
import { UserCard } from "./usercard";

export const Header = () => {
  const profileData = useProfile();

  return (
    <header className="flex justify-around py-4 ">
      <div className="flex items-center">
        <h1 className="text-xl font-bold">Мои услуги</h1>
      </div>
      <div className="flex items-center">
        {profileData.data ? (
          <div className="flex gap-x-1 items-center">
            <UserCard user={profileData.data} />
          </div>
        ) : (
          <Link to={"/sign-in"}>
            <Button>Авторизация</Button>
          </Link>
        )}
      </div>
    </header>
  );
};
