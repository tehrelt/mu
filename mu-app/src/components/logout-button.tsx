"use client";

import { logout } from "@/actions/auth";
import { Button } from "./ui/button";

export const LogoutButton = () => {
  return <Button onClick={() => logout()}>Logout</Button>;
};
