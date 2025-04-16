import { Link } from "@remix-run/react";
import { Profile } from "~/schemes/profile";
import { Button } from "./ui/button";
import { UserCard } from "./usercard";

type HeadersProps = {
  profile?: Profile;
};

export const Header = ({ profile }: HeaderProps) => {
  return (
    <header className="flex py-4 px-2">
      <div className="flex-grow">
        <h1>Мои услуги</h1>
      </div>
      <div>
        <div className="flex items-center">
          {profile ? (
            <div>
              <UserCard profile={profile} />
              <Button onClick={() => {}}>Выйти</Button>
            </div>
          ) : (
            <Link to="/login">
              <Button>Войти</Button>
            </Link>
          )}
        </div>
      </div>
    </header>
  );
};
