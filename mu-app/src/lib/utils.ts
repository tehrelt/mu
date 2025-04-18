import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import dayjs from "dayjs";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function fio(
  {
    firstName,
    lastName,
    middleName,
  }: {
    firstName?: string;
    lastName?: string;
    middleName?: string;
  },
  opts?: { full?: boolean },
) {
  if (opts?.full) {
    return [lastName ?? "", firstName ?? "", middleName ?? ""].join(" ");
  }

  return [
    lastName ? lastName : "",
    firstName ? firstName.charAt(0) + "." : "",
    middleName ? middleName.charAt(0) + "." : "",
  ].join(" ");
}

export function datef(date: Date, format?: string): string {
  return dayjs(date).format(format ?? "DD/MM/YYYY HH:mm:ss");
}
