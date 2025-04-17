import { logout } from "@/actions/auth";
import { LogoutButton } from "@/components/logout-button";
import { Button } from "@/components/ui/button";
import { authService } from "@/service/auth.service";
import Link from "next/link";

export default async function Home() {
  const profile = await authService.profile();

  return (
    <div className="flex justify-around">
      <h1>Мои Услуги</h1>
      {profile ? (
        <div>
          {profile.lastName} {profile.firstName}
          <LogoutButton />
        </div>
      ) : (
        <Link href="/login" className="cursor-pointer">
          <Button variant={"link"}>Login</Button>
        </Link>
      )}
    </div>
  );
}
